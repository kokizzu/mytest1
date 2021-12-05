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
	"github.com/kokizzu/gotro/S"
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
	myDsn := `%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=true`
	myDsn = fmt.Sprintf(myDsn,
		dbUser,
		dbPass, // empty password
		dbHost,
		dbPort,
		dbName,
	) // mysql -u root -h 127.0.0.1 -P 4000 test

	// if without parseTime=true, must add these function to convert from &[]byte to each concrete value
	//convertFunc := map[string]func(interface{}) interface{}{
	//	`id`:         func(v interface{}) interface{} { return X.ToI(v) },
	//	`uniq`:       func(v interface{}) interface{} { return X.ToS(v) },
	//	`created_at`: func(v interface{}) interface{} { return X.ToTime(v) },
	//	`updated_at`: func(v interface{}) interface{} { return X.ToTime(v) },
	//}

	loggerAdapter := onelogadapter.New(onelog.New(os.Stdout, onelog.FATAL))
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
	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second) // everything must complete in 90s
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

	start := time.Now()
	defer func() {
		fmt.Println()
		fmt.Println(`Done in `, time.Since(start))
	}()
	msTiming := M.SX{}
	//defer L.Describe(msTiming)
	const N = 10000

	SetMinMaxTotal := func(key string, t time.Duration) {
		v := msTiming.GetInt(key + ` Min`)
		if v == 0 || v > int64(t) {
			msTiming.Set(key+` Min`, t)
		}
		v = msTiming.GetInt(key + ` Max`)
		if v == 0 || v < int64(t) {
			msTiming.Set(key+` Max`, t)
		}
		v = msTiming.GetInt(key + ` Total`)
		msTiming.Set(key+` Total`, v+int64(t))
	}

	// query/select example scan to map manual
	for z := 0; z < N; z++ {
		start := time.Now()
		cols := []string{`id`, `uniq`, `created_at`, `updated_at`}

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
			// convert to map[string]interface{}
			m := M.SX{}
			for k, col := range cols {
				if col == `uniq` {
					m[col] = X.ToS(row[k]) // convert &[]uint8 to string
				} else {
					m[col] = row[k] // convertFunc[col](row[k])
				}
			}
			res = append(res, m)
		}
		//L.Describe(res)
		SetMinMaxTotal(`[]MapManual`, time.Since(start))
	} //*/

	// query/select example scan to map
	for z := 0; z < N; z++ {
		start := time.Now()
		cols := []string{`id`, `uniq`, `created_at`, `updated_at`}
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
		res := []M.SX{}
		for rows.Next() {
			row := map[string]interface{}{}
			err := rows.MapScan(row)
			for k, v := range row {
				if k == `uniq` {
					row[k] = X.ToS(v) // convert &[]uint8 to string
				}
			}
			L.PanicIf(err, `row.ScanSlice:IN_LIKE_QUERY`)
			res = append(res, row)
		}
		//L.Describe(res)
		SetMinMaxTotal(`[]Map`, time.Since(start))
	} //*/

	// query/select example scan to slice manual
	for z := 0; z < N; z++ {
		start := time.Now()
		cols := []string{`id`, `uniq`, `created_at`, `updated_at`}
		isStringColumn := []bool{false, true, false, false}
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
		res := []interface{}{}
		for rows.Next() {
			err := rows.Scan(row...)
			copy := make([]interface{}, len(cols))
			for k, v := range row {
				if isStringColumn[k] {
					copy[k] = X.ToS(v)
				} else {
					copy[k] = v
				}
			}
			L.PanicIf(err, `row.Scan:IN_LIKE_QUERY`)
			res = append(res, copy)
		}
		//L.Describe(res)
		SetMinMaxTotal(`[]SliceManual`, time.Since(start))
	} //*/

	// query/select example scan to struct
	for z := 0; z < N; z++ {
		start := time.Now()
		cols := []string{`id`, `uniq`, `created_at`, `updated_at`}
		isStringCol := []bool{false, true, false, false}
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
		res := []interface{}{}
		for rows.Next() {
			row, err := rows.SliceScan()
			copy := make([]interface{}, len(cols))
			for k, v := range row {
				if isStringCol[k] {
					copy[k] = X.ToS(v)
				} else {
					copy[k] = v
				}
			}
			L.PanicIf(err, `row.ScanSlice:IN_LIKE_QUERY`)
			res = append(res, copy)
		}
		//L.Describe(res)
		SetMinMaxTotal(`[]Slice`, time.Since(start))
	} //*/

	// query/select example scan to pointer of struct
	for z := 0; z < N; z++ {
		start := time.Now()
		cols := []string{`id`, `uniq`, `created_at`, `updated_at`}
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
		res := []*Users{} // we don't know the size, if len is small/known and struct is small, it's better to use non-pointer
		for rows.Next() {
			row := Users{}
			err := rows.StructScan(&row)
			L.PanicIf(err, `row.ScanStruct:IN_LIKE_QUERY`)
			res = append(res, &row)
		}
		//L.Describe(res)
		SetMinMaxTotal(`[]*Struct`, time.Since(start))
	} //*/

	// query/select example scan to pointer of struct
	for z := 0; z < N; z++ {
		start := time.Now()
		cols := []string{`id`, `uniq`, `created_at`, `updated_at`}
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
		res := []Users{} // we don't know the size, if len is large and struct is large, it's better to use pointer
		for rows.Next() {
			row := Users{}
			err := rows.StructScan(&row)
			L.PanicIf(err, `row.ScanStruct:IN_LIKE_QUERY`)
			//L.Describe(row)
			res = append(res, row)
		}
		//L.Describe(res)
		SetMinMaxTotal(`[]Struct`, time.Since(start))
	} //*/

	// conclusion: scan to map or slice must convert from &[]uint8 manually when using any pointer interface{} (either to string; or to time.Time -- unless parseTime=true specified)

	keys := msTiming.SortedKeys()
	fmt.Printf("Time for %d operations:\n", N)
	last := ``
	for _, key := range keys {
		curr := S.LeftOf(key, ` `)
		if curr != last {
			last = curr
			fmt.Println()
		}
		ns := msTiming.GetInt(key)
		fmt.Printf("%28s %12d ns = %9d µs = %9.2f ms = %7.4f s\n", key, ns, ns/1000, float64(ns)/1000/1000, float64(ns)/1000/1000/1000)
		if S.EndsWith(key, `Total`) {
			ns := ns / N
			key = S.Replace(key, `Total`, `Avg`)
			fmt.Printf("%28s %12d ns = %9d µs = %9.2f ms = %7.4f s\n", key, ns, ns/1000, float64(ns)/1000/1000, float64(ns)/1000/1000/1000)
		}
	}
}

func callerComment(info *L.CallInfo) string {
	return fmt.Sprintf("-- %s %s.%s:%d", info.FuncName, info.PackageName, info.FuncName, info.Line)
}
