# GoSQL

GoSQL is an SQL ORM for the Go programming language.
It is intended to add type safety and ease of use to SQL.

## Installation

```shell
go get github.com/tkdeng/gosql

# Any SQL Driver

# SQLite
go get github.com/mattn/go-sqlite3

# MySQL
go get github.com/go-sql-driver/mysql
```

## Usage

```go

import (
  "github.com/tkdeng/gosql"

  // SQLite
  _ "github.com/mattn/go-sqlite3"

  // MySQL
  _ "github.com/go-sql-driver/mysql"
)

func main(){
  // Local File
  db, err := gosql.Open("sqlite3", "/path/to/db.sqlite")
  
  // Memory/RAM
  db, err := gosql.Open("sqlite3", "")

  // Server
  db, err := gosql.Open("mysql", gosql.Server{
    username: "user",
    password: "p@ssw0rd!",
    host: "localhost",
    port: 1433,
    protocol: "tcp", // or udp
    database: "db",
  })

  // close database
  defer db.Close()
}

```

### Adding data to a table

```go
// Create or Select a Table
// note: this method will automatically create a table if it doesn't exist
table := db.Table("users",
  INT("id").Primary().Unique().Default(0), // note: `.AutoInc()` is not supported by sqlite
  TEXT("username"),
  TEXT("password"),
)

// INSERT new row
err := table.Set(map[string]any{
  "username": "user",
  "password": "p@ssw0rd!",
})

// INSERT or UPDATE row
err := table.Set(map[string]any{
  "username": "user",
  "password": "NewPassword!",
}, "username") // optional: specify unique keys

// UPDATE row
err := table.Where("id").Equal(0).Set(map[string]any{
  "username": "user",
  "password": "NewerPassword!",
})

// In the above example, if the database finds that
// "username" = "user" already exists, it will update the
// "password" of the existing user, instead of creating
// a new user. If not found, a new user will be created.
```

### Getting data from a table

```go
// SELECT data from database
err := table.Get([]string{"id", "username", "password"}, func(scan func(dest ...any) error) bool {
  var id int
  var username string
  var password string

  scan(&id, &username, &password) // runs rows.Scan from core sql module

  // return true to continue rows.Next()
  // return false to break the loop and close the query
  return true
})

// check if table has a row WHERE key = value
if table.Has(map[string]any{"username": "admin"}) {
  // admin user exists
}
```

### adding WHERE query

```go
err := table.Where("username") // WHERE username
  .Equal("user") // = 'user'
  .AND("password") // AND password
  .EQUAL("p@ssw0rd!") // = 'p@ssw0rd!'
  .Get(nil, func(scan func(dest ...any) error) bool { // SELECT * FROM table
    // do stuff
  })

query := table.Where("username", false).Equal("admin") // WHERE NOT username = 'admin'

query := table.Where("username").NotEqual("admin") // WHERE username <> 'admin'
// note: `<>` is equivalent to `!=` in sql

query := table.Where("id").GreaterThan(0) // WHERE id > 0

query := table.OrderBy("id") // ORDER BY id
query := table.OrderBy("id", true) // ORDER BY id DESC

// note: the Where method returns a new instance of the query
table.Where("id").Equal(0)

err := table.Delete() // will not register the above where query
err == gosql.Error_UnsafeQuery
// For safety, the Delete method will actually return an error if
// the WHERE query is empty by default.
// This prevents you from accidently deleting the entire table.
// (more info about safety checks below)
```

### Removing data from a table

```go
// delete row from database
err := table.Where("password").Equal("p@ssw0rd!").Delete()

// note: the above example will delete any user with the password "p@ssw0rd!"

// If you plan on removing insucure passwords, I would recommend locking the account
// by updating it with a random password, and let the user change it via email or the
// 'forgot password' button. Your users will likely be mad if you randomly delete their account.


// running Delete, without a Where query will return an error
err := table.Delete()
err == gosql.Error_UnsafeQuery

// to override this (set @force = true)
err := table.Delete(true)
err == nil
// note: the above method will only delete all rows, and keeps the empty table

// to Drop the table
table.Drop(true) // note: you must pass 'true' to confirm dropping the table
```

### Query safety checks

```go
// to run raw sql queries
db.Query() || db.Exec() || db.Prepare()
// these are a 1 to 1 of the core sql params, but with an additional safety check on the query
// note: this module uses the above methods when running sql queries (eccept for the Drop method)


// this method gets run by default, but you can also call it manually
if db.SafeQuery("SELECT * FROM users") {
  // this passed safety checks
}

if db.SafeQuery("SELECT * FROM users WHERE username = 'admin' OR 1=1") {
  // this will fail the safety check by default

  // common sql injection includes escaping an input and adding `OR 1=1`,
  // if any WHERE query checks if something is equal to itself,
  // it will fail the default safety check.
}

if db.SafeQuery("SELECT * FROM users WHERE username = 'admin'; DROP *") {
  // this will fail the safety check by default

  // common sql injection loves to abuse the `;` character.
  // this feature is not frequently needed, and can simply
  // be replaced by running another query function.
}

if db.SafeQuery("DROP *") {
  // this will fail the safety check by default

  // this is not only to protect against sql injection,
  // but to also prevent developers from accidently deleting
  // an entire database.

  // note: the `table.Drop(true)` method will override safety checks.
}


// adding your own safety checks

// this method will add another callback to a list
// that will be called whenever a query safety check is ran.
gosql.AddSafetyCheck(func(query string) bool {
  if safe {
    return true // return true to continue the safety check list
  } else {
    return false // return false to report an unsafe query string
  }
})

// for a simple regex match, this method will check if the
// query matched a regular expression (RE2).
// if the query matches this regex, it will repoort the query as unsafe
gosql.AddSafetyCheckRE(`(?i)(WHERE|AND|OR)\s+NOT`) // this example prevents any `WHERE NOT` queries

// pass @where: true, to only check after the WHERE query
gosql.AddSafetyCheckRE(`(?i)NOT`, true)


// overriding safety checks (Not Recommended)
db = db.Unsafe("I Know What Im Doing!") // returns a new database instance that allows unsafe queries

// note: the raw database object from the core sql module will also bypass query safety checks
var rawDB *sql.DB = db.SQL
```
