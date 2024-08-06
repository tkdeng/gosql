package gosql

import (
	"strings"
)

// Get will SELECT keys FROM table, and run a loop over the selected rows
//
// @cb: will be called for every row
//   - return true, to continue the loop
//   - return false, to close the query and break the loop
func (query *Query) Get(keys []string, cb func(scan func(dest ...any) error) bool) error {
	q := `SELECT `
	if len(keys) == 0 {
		q += `*`
	} else {
		for i := 0; i < len(keys); i++ {
			q += toAlphaNumeric(keys[i])
			if i != len(keys)-1 {
				q += `, `
			}
		}
	}

	q += ` FROM ` + query.table

	if query.where != "" {
		q += ` ` + query.where
	}

	if query.order != "" {
		q += ` ` + query.order
	}

	rows, err := query.db.Query(q, query.whereValue...)
	if err != nil {
		return err
	}

	for rows.Next() {
		if !cb(rows.Scan) {
			break
		}
	}
	rows.Close()

	return nil
}

// Has will check if key value pairs are found in the database (using SELECT WHERE)
//
// If a row is found, this method will return true.
// If nothing is found, or an error occurs, this method will return false.
func (query *Query) Has(values map[string]any) bool {
	if len(values) == 0 {
		return false
	}

	valList := []any{}

	q := `SELECT * FROM ` + query.table + ` WHERE `
	for key, val := range values {
		q += toAlphaNumeric(key) + ` = ? AND `
		valList = append(valList, val)
	}
	q = q[:len(q)-5]

	if query.where != "" {
		q += ` AND` + strings.TrimPrefix(query.where, "WHERE")
		valList = append(valList, query.whereValue...)
	}

	if rows, err := query.db.Query(q, valList...); err == nil && rows.Next() {
		rows.Close()
		return true
	}

	return false
}
