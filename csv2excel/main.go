// I always got lots of problems (different encoding between osx and windows)
// when sending csv file to PM or BM. So I have to convert to excel manually.
// And ssconvert from gnumeric has too much dependencies.
// Usage
// csv2excel -in xxx.csv -out yyy.xlsx
// OR csv2excel -in a/b/c/xxx.csv   (will got a/b/c/xxx.xlsx)
// OR cat xxx.csv | csv2excel -out yyy.xlsx

package main

import (
	"encoding/csv"
	"flag"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize"
)

type options struct {
	infile    string
	outfile   string
	delimiter string
}

var cliOpt options

func init() {
	flag.StringVar(&cliOpt.infile, "in", "", "path of csv file")
	flag.StringVar(&cliOpt.outfile, "out", "", "path of output excel file")
	flag.StringVar(&cliOpt.delimiter, "delimiter", "", "delimiter of csv")
}

func main() {
	flag.Parse()

	in := os.Stdin
	var err error
	if cliOpt.infile != "" {
		if in, err = os.Open(cliOpt.infile); err != nil {
			log.Fatalf("failed to open input file, %v\n", err)
		}
		defer in.Close()
	}

	if cliOpt.outfile == "" {
		// need specify output file if no infile (stream from Stdin)
		if in == os.Stdin {
			log.Fatalf("must specify input file or output file\n")
		}
		out := strings.TrimSuffix(cliOpt.infile, ".csv")
		cliOpt.outfile = out + ".xlsx"
	}

	reader := csv.NewReader(in)
	if cliOpt.delimiter != "" {
		reader.Comma = rune(cliOpt.delimiter[0])
	}
	f := excelize.NewFile()

	var excelHeader []int
	i := 0
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("error in read csv file %v\n", err)
		}
		if i == 0 {
			excelHeader = make([]int, len(record), len(record))
			for j, _ := range record {
				excelHeader[j] = 'A' + j
			}
		}
		i++
		for j, cell := range record {
			cellpos := string(excelHeader[j]) + strconv.Itoa(i)
			f.SetCellValue("Sheet1", cellpos, cell)
		}
	}

	err = f.SaveAs(cliOpt.outfile)
	if err != nil {
		log.Fatal(err)
	}
}
