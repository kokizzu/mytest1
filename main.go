package main

import (
	"context"
	"fmt"
	"github.com/francoispqt/onelog"
	"github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/kokizzu/goproc"
	"github.com/kokizzu/gotro/A"
	"github.com/kokizzu/gotro/I"
	"github.com/kokizzu/gotro/L"
	"github.com/kokizzu/gotro/M"
	"github.com/kokizzu/gotro/X"
	"github.com/kokizzu/id64"
	sqldblogger "github.com/simukti/sqldb-logger"
	"github.com/simukti/sqldb-logger/logadapter/onelogadapter"
	"os"
	"time"
)

const driverDialect = "mysql"
const dbHost = `127.0.0.1`
const dbUser = `root`
const dbPass = ``
const dbPort = 4000
const dbName = `test`

var globalConn *sqlx.DB

func init() {
	myDsn := `%s:%s@tcp(%s:%d)/%s?charset=utf8`
	myDsn = fmt.Sprintf(myDsn,
		dbUser,
		dbPass, // empty password
		dbHost,
		dbPort,
		dbName,
	) // mysql -u root -h 127.0.0.1 -P 4000 test

	loggerAdapter := onelogadapter.New(onelog.New(os.Stdout, onelog.ALL))
	db := sqldblogger.OpenDriver(myDsn, &mysql.MySQLDriver{}, loggerAdapter /*, ...options */)
	globalConn = sqlx.NewDb(db, driverDialect)

	proc := goproc.New()
	proc.AddCommand(&goproc.Cmd{
		Program:    `migu`,
		Parameters: []string{`sync`, `-u`, dbUser, `-h`, dbHost, `-P`, I.ToS(dbPort), dbName, `schema.go`},
	})
	proc.StartAll()
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second) // everything must complete in 20s
	defer cancel()

	// NOTE: do not use PanicIf for long running service, handle error normally

	// exec/ddl example
	truncateQuery := callerComment(L.CallerInfo()) + `
TRUNCATE TABLE users`
	_, err := globalConn.ExecContext(ctx, truncateQuery)
	L.PanicIf(err, `ExecContext:TRUNCATE:users`)

	// exec/insert example
	insertQuery := callerComment(L.CallerInfo()) + `
INSERT INTO users(uniq) VALUES(?)`
	for z := 0; z < 100; z++ {
		//user := Users{Uniq: id64.SID()}
		//_, err := o.Insert(&user) // &mysql.MySQLError{Number:0x418, Message:"Column 'created_at' cannot be null"}
		_, err := globalConn.ExecContext(ctx, insertQuery, id64.SID())
		L.PanicIf(err, `ExecContext:INSERT:users`)
	}

	// query/select example
	cols := []string{`id`, `uniq`, `created_at`, `updated_at`}
	convertFunc := map[string]func(interface{}) interface{}{
		`id`:         func(v interface{}) interface{} { return X.ToI(v) },
		`uniq`:       func(v interface{}) interface{} { return X.ToS(v) },
		`created_at`: func(v interface{}) interface{} { return X.ToTime(v) },
		`updated_at`: func(v interface{}) interface{} { return X.ToTime(v) },
	}
	params := []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	sql := callerComment(L.CallerInfo()) + `
SELECT ` + A.StrJoin(cols, `,`) + ` 
FROM users 
WHERE id IN (?)
	OR uniq LIKE ?
LIMIT 20`
	selectQuery, args, err := sqlx.In(sql, params, `%a%`)
	L.PanicIf(err, `sqlx.In`)
	rows, err := globalConn.QueryxContext(ctx, selectQuery, args...)
	L.PanicIf(err, `QueryxContext:IN_LIKE_QUERY`)
	defer rows.Close()
	row := make([]interface{}, len(cols))
	for k := range row {
		row[k] = new(interface{})
	}
	res := []M.SX{}
	for rows.Next() {
		err := rows.Scan(row...)
		L.PanicIf(err, `row.Scan:IN_LIKE_QUERY`)
		L.Describe(row)

		// convert to map[string]interface{}
		m := M.SX{}
		for k, col := range cols {
			m[col] = convertFunc[col](row[k])
		}
		res = append(res, m)
	}
	L.Describe(res)
}

func callerComment(info *L.CallInfo) string {
	return fmt.Sprintf("-- %s %s.%s:%d", info.FuncName, info.PackageName, info.FuncName, info.Line)
}
