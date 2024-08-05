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
    username: "SQLUser",
    password: "p@ssw0rd!",
    host: "localhost",
    port: 1433,
    protocol: "tcp", // or udp
    database: "",
  })
}

```
