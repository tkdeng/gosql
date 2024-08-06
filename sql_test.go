package gosql

import (
	"testing"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/tkdeng/gosql/common"
)

func Test(t *testing.T) {
	db, err := Open("sqlite3", "")
	if err != nil {
		t.Error(err)
	}

	table := db.Table("users", TEXT("username"), TEXT("password"))

	err = table.Set(map[string]any{
		"username": "admin",
		"password": "12345",
	}, "username")
	if err != nil {
		t.Error(err)
	}

	err = table.Set(map[string]any{
		"username": "user",
		"password": "p@ssw0rd!",
	}, "username")
	if err != nil {
		t.Error(err)
	}

	expect := [][]string{
		{"admin", "12345"},
		{"user", "p@ssw0rd!"},
	}

	err = table.Get([]string{"username", "password"}, func(scan func(dest ...any) error) bool {
		var username string
		var password string

		scan(&username, &password)

		if i, err := common.IndexOf(expect, []string{username, password}); err == nil {
			expect = append(expect[:i], expect[i+1:]...)
		} else {
			t.Error("database does not contain:", "[" + username + " " + password + "]")
		}

		return true
	})
	if err != nil {
		t.Error(err)
	}

	if len(expect) != 0 {
		t.Error("database failed to get:", expect)
	}

	err = table.Where("password").Equal("p@ssw0rd!").Delete()
	if err != nil {
		t.Error(err)
	}

	err = table.Drop(true)
	if err != nil {
		t.Error(err)
	}

	db.Close()
}

func TestSaefty(t *testing.T) {
	db, err := Open("sqlite3", "")
	if err != nil {
		t.Error(err)
	}

	db.Table("users", TEXT("username"), TEXT("password"))

	testSaefty := func(query string) {
		_, err = db.Query(query)
		if err != Error_UnsafeQuery {
			if err == nil {
				t.Error("failed to detect unsafe query")
			} else {
				t.Error(err)
			}
		}
	}

	// deny empty query
	testSaefty("")

	// deny `DROP` keyword
	testSaefty("DROP *")

	// deny common `username = '*' OR password = '*'` from WHERE query
	// (but allow other `key = '*'` queries)
	testSaefty("SELECT * FROM users WHERE username = '*'")

	// deny `1=1` or `key=self` from WHERE query
	testSaefty("SELECT * FROM users WHERE username = 'admin' OR 1=1")

	// deny `;` in sql queries (we should be creating separate query requests instead)
	// the `;` is  often abused by hackers, and rarly needed by servers
	testSaefty("SELECT * FROM users WHERE username = 'admin'; SELECT * FROM users")

	db.Close()
}

func TestServer(t *testing.T) {
	//todo: test sql server
	/* db, err := Open("mysql", Server{

	})
	if err != nil {
		t.Error(err)
	}
	_ = db */
}
