package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"strings"

	_ "github.com/denisenkom/go-mssqldb" // ms sqlserver
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"github.com/zxdvd/go-libs/std-helper/M"
)

type options struct {
	db      string
	dialect string
	uri     string
	query   string
	limit   int
	tocsv   bool
	tojson  bool
}

var cliOpt options

func init() {
	flag.StringVar(&cliOpt.db, "db", "", "db settings in tools.json")
	flag.StringVar(&cliOpt.dialect, "dialect", "", "mysql/postgres/sqlserver")
	flag.StringVar(&cliOpt.uri, "uri", "", "database uri, postgres://xx@localhost/example")
	flag.StringVar(&cliOpt.query, "query", "", "query, select id,name from table1")
	flag.IntVar(&cliOpt.limit, "limit", 10000, "limit rows of result, 0 means 0 limit, default 10000")
	flag.BoolVar(&cliOpt.tocsv, "csv", false, "output to csv")
	flag.BoolVar(&cliOpt.tojson, "json", true, "output to json")
}

func parseOption() {
	flag.Parse()
	if cliOpt.db != "" {
		key := "dbquery." + cliOpt.db
		dbconf, err := Get("tools.json", key)
		if err != nil {
			fmt.Printf("err: %v\n", err)
			panic(err)
		}
		if cliOpt.dialect == "" {
			cliOpt.dialect = M.GetString(dbconf, "dialect", "")
		}
		if cliOpt.uri == "" {
			cliOpt.uri = M.GetString(dbconf, "uri", "")
		}
	}
	if cliOpt.tocsv {
		cliOpt.tojson = false
	}
}

func main() {
	parseOption()
	db, err := sql.Open(cliOpt.dialect, cliOpt.uri)
	if err != nil {
		log.Fatalf("failed to open database, %s", cliOpt.uri)
	}
	q, err := NewQuery(db, cliOpt.query)
	if err != nil {
		log.Fatal(err)
	}
	rowIdx := 0
	if cliOpt.tojson {
		fmt.Println("[")
		for row := range q.C {
			newRow := tomap(q.Columns, row)
			output, err := json.Marshal(newRow)
			if err != nil {
				log.Fatal(err)
			}
			fmt.Println(string(output))
			rowIdx++
			if cliOpt.limit > 0 && rowIdx >= cliOpt.limit {
				break
			}
		}
		fmt.Println("]")
		return
	}
	if cliOpt.tocsv {
		fmt.Println(strings.Join(q.Columns, ","))
		row := make([]string, len(q.Columns))
		for row_ := range q.C {
			for i, col := range row_ {
				if col == nil {
					row[i] = ""
					continue
				}
				v := valueFormat(col)
				row[i] = csvQuote(fmt.Sprintf("%v", v))
			}
			fmt.Println(strings.Join(row, ","))
			rowIdx++
			if cliOpt.limit > 0 && rowIdx >= cliOpt.limit {
				break
			}
		}
		return
	}
}

type Query struct {
	Columns []string
	C       chan []interface{}
}

func NewQuery(db *sql.DB, query string) (*Query, error) {
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	cols, err := rows.ColumnTypes()
	if err != nil {
		return nil, err
	}
	colnames := make([]string, len(cols))
	for i, col := range cols {
		colnames[i] = col.Name()
	}

	q := Query{
		Columns: colnames,
		C:       make(chan []interface{}, 100),
	}

	go func() {
		colVals := make([]interface{}, len(cols))
		colPointers := make([]interface{}, len(cols))
		for rows.Next() {
			for i := range colVals {
				colPointers[i] = &colVals[i]
			}
			if err := rows.Scan(colPointers...); err != nil {
				panic(err)
			}
			q.C <- append([]interface{}(nil), colVals...)
		}
		rows.Close()
		close(q.C)

	}()
	return &q, nil
}

func tomap(cols []string, row []interface{}) map[string]interface{} {
	output := make(map[string]interface{}, len(row))
	for j, col := range cols {
		output[col] = valueFormat(row[j])
	}
	return output
}

func valueFormat(v interface{}) interface{} {
	if bytes, ok := v.([]byte); ok {
		return string(bytes)
	}
	return v
}

func csvQuote(s string) string {
	return "\"" + strings.Replace(s, "\"", "\"\"", -1) + "\""
}
