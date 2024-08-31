package gosql

import (
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/tkdeng/goregex"
	"github.com/tkdeng/goutil"
)
type DB struct {
	SQL        *sql.DB
	initTables []string
	unsafe     bool
}

type Server struct {
	username string
	password string
	host     string
	port     uint16
	protocol string
	database string
}

var Error_UnsafeQuery = errors.New("unsafe query")

var querySafetyChecks []func(query string) bool = []func(query string) bool{}
var querySafetyChecksRE []*regex.Regexp = []*regex.Regexp{}
var querySafetyChecksWhereRE []*regex.Regexp = []*regex.Regexp{}

// Open opens a new database
func Open[T interface{ string | Server }](driverName string, dns T) (*DB, error) {
	//todo: add support for sql auth and cloudflare D1 or R2

	var dnsVal interface{} = dns

	var dbDNS string
	if path, ok := dnsVal.(string); ok {
		if path == "" {
			dbDNS = "file::memory:?cache=shared"
		} else {
			path = string(regex.Comp(`[^\w_\-:\\/@$#!+~\.\,\s ]`).RepStrLit([]byte(path), []byte{}))
			dbDNS = "file:" + path + "?cache=shared"
		}
	} else if server, ok := dnsVal.(Server); ok {
		if server.protocol == "" {
			server.protocol = "tcp"
		}
		dbDNS = fmt.Sprintf("%s:%s@%s(%s:%d)", server.username, server.password, server.protocol, server.host, server.port)
		if server.database != "" {
			dbDNS += "/" + server.database
		}
	} else {
		return nil, errors.New("invalid dns")
	}

	db, err := sql.Open(driverName, dbDNS)
	if err != nil {
		return nil, err
	}

	if strings.HasPrefix(dbDNS, "file:") || dbDNS == ":memory:" {
		db.SetMaxOpenConns(1)
	} else {
		db.SetConnMaxLifetime(time.Minute * 3)
		db.SetMaxOpenConns(10)
		db.SetMaxIdleConns(10)
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &DB{
		SQL:        db,
		initTables: []string{},
	}, nil
}

// Close closes the database
func (db *DB) Close() {
	db.SQL.Close()
}

// Table selects a database table
//
// if any rows are specified, this method will create a table if it does not exist
func (db *DB) Table(name string, rows ...*DataType) *Query {
	name = toAlphaNumeric(name)

	if len(rows) != 0 && !goutil.Contains(db.initTables, name) {
		query := `CREATE TABLE IF NOT EXISTS ` + name + ` (`
		for i, row := range rows {
			query += row.key

			if row.def != "" {
				query += ` DEFAULT ` + row.def
			}

			if i != len(rows)-1 {
				query += `, `
			}
		}
		query += `)`

		if st, err := db.SQL.Prepare(query); err == nil {
			if _, err = st.Exec(); err == nil {
				db.initTables = append(db.initTables, name)
			}
		}
	}

	return &Query{
		db:    db,
		table: name,
	}
}

// Unsafe will disable sql safety checks
//
// By default, the final output will be checked for potentially dangorous queries,
// like `DROP *` which could happen by accident.
//
// Note: this safety check does Not guarantee safety, and should Not be relied on.
//
// To use this method, you must pass the confirm argument as "I Know What Im Doing!",
// to confirm that you have read the documentation, and know what you are doing.
func (db DB) Unsafe(confirm string) *DB {
	if confirm == "I Know What Im Doing!" {
		db.unsafe = true
	}

	return &db
}

// Query executes a query that returns rows, typically a SELECT.
// The args are for any placeholder parameters in the query.
//
// Query uses [context.Background] internally; to specify the context, use
// [DB.QueryContext].
func (db *DB) Query(query string, args ...any) (*sql.Rows, error) {
	if !db.unsafe && !SafeQuery(query) {
		return nil, Error_UnsafeQuery
	}
	return db.SQL.Query(query, args...)
}

// Exec executes a query without returning any rows.
// The args are for any placeholder parameters in the query.
//
// Exec uses [context.Background] internally; to specify the context, use
// [DB.ExecContext].
func (db *DB) Exec(query string, args ...any) (sql.Result, error) {
	if !db.unsafe && !SafeQuery(query) {
		return nil, Error_UnsafeQuery
	}
	return db.SQL.Exec(query, args...)
}

// Prepare creates a prepared statement for later queries or executions.
// Multiple queries or executions may be run concurrently from the
// returned statement.
// The caller must call the statement's [*Stmt.Close] method
// when the statement is no longer needed.
//
// Prepare uses [context.Background] internally; to specify the context, use
// [DB.PrepareContext].
func (db *DB) Prepare(query string) (*sql.Stmt, error) {
	if !db.unsafe && !SafeQuery(query) {
		return nil, Error_UnsafeQuery
	}
	return db.SQL.Prepare(query)
}

// SafeQuery checks a query for common safety errors
//
// Note: this safety check does Not guarantee safety, and should Not be relied on.
//
// This method is called by default, unless `db.SQL` is used to access raw sql from
// the default `database/sql` module.
func SafeQuery(query string) bool {
	if strings.TrimSpace(query) == "" {
		return false
	}

	// check for bad keywords
	if regex.Comp(`(?is)(DROP|;)`).Match([]byte(query)) {
		return false
	}

	// check for bad where query
	safe := true
	regex.Comp(`WHERE(.*)$`).RepFunc([]byte(query), func(data func(int) []byte) []byte {
		q := data(1)

		// common sql injection: username = * AND password = *
		if regex.Comp(`(?is)["'\']?(user(name|id|)|pass(word|)|u*id)["'\']?\s*=\s*["'\']?\*["'\']?`).Match(q) {
			safe = false
			return nil
		}

		// common sql injection: 1=1
		regex.Comp(`(?is)["'\']?([\w_\-]*)["'\']?\s*=\s*["'\']?([\w_\-]*)["'\']?`).RepFunc(q, func(data func(int) []byte) []byte {
			if bytes.Equal(data(1), data(2)) {
				safe = false
			}
			return []byte{}
		})

		// check custom where regex list
		for _, reg := range querySafetyChecksWhereRE {
			if reg.Match([]byte(query)) {
				safe = false
				break
			}
		}

		return []byte{}
	})
	if !safe {
		return false
	}

	// check custom regex list
	for _, reg := range querySafetyChecksRE {
		if reg.Match([]byte(query)) {
			return false
		}
	}

	// check custom callback list
	for _, cb := range querySafetyChecks {
		if !cb(query) {
			return false
		}
	}

	return true
}

// AddSafetyCheck adds another safety check to the SafeQuery method
//
// return false, if you think the query looks unsafe.
// return true, to continue down the safety check list.
func AddSafetyCheck(cb func(query string) bool) {
	querySafetyChecks = append(querySafetyChecks, cb)
}

// AddSafetyCheckRE adds another safety check to the SafeQuery method
//
// The `RE` stands for RegExp, so you can simply pass a regex string,
// which will check for a match, instead of a full callback method.
//
// If the query matches this regex, the query will be seen as unsafe.
//
// @where: if true, will only check after the WHERE keyword
func AddSafetyCheckRE(re string, where ...bool) {
	if len(where) != 0 && where[0] {
		querySafetyChecksWhereRE = append(querySafetyChecksWhereRE, regex.Comp(re))
	} else {
		querySafetyChecksRE = append(querySafetyChecksRE, regex.Comp(re))
	}
}
