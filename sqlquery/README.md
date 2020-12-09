
## sqlquery

A simple cli help to query postgres/mysql/sqlserver and fetch result to csv/json.

Each row outputs as json map and it will output 10000 rows at most by default. Use `-limit 0`
 to output all.

It uses a small buffer to output streammingly so don't worry about memory usage
 even for large table.

### direct query without config

    $ ./sqlquery -dialect sqlserver -uri 'sqlserver://app:123456@127.0.0.1?database=db1&encrypt=disable' -query 'select 100;'
    $ ./sqlquery -dialect postgres -uri 'postgres://u1:123456@127.0.0.1:5432/db1?sslmode=disable' -query 'select * from pg_class limit 10'


### query using config
You can write `tools.json` file at `cwd` or `$HOME/tools.json` to setup multiple
 database connections like following:

```
{
  "__comment": "config for all my tools",
  "dbquery": {
    "__comment": "db config for dbquery",
    "db1": {
      "dialect": "postgres",
      "uri": "postgres://apollon@localhost:5432/db1?sslmode=disable"
    },
    "db2": {
      "dialect": "mysql",
       "uri": "user1:password1@tcp(192.168.1.1:3306)/user?parseTime=true"
    }
  }
}
```

Then you can query any database using alias

    $ ./sqlquery --db db1 -query 'select 10 as a'
    $ ./sqlquery --db db2 -query 'select 10 as b' -csv         # to csv




