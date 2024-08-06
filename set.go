package gosql

//todo: add methods for `CREATE INDEX` and `DROP INDEX`: https://www.w3schools.com/sql/sql_create_index.asp

// Set will INSERT or UPDATE values FROM table
//
// If a where query exists, this method will only use UPDATE.
//
// If the unique arg exists, this method will check if the database contains
// all matching key = value pairs. If it finds and, it will use UPDATE, and
// if nothing is found, it will use INSERT.
//
// If no unique args or where query exists, this method will default to INSERT.
func (query *Query) Set(values map[string]any, unique ...string) error {
	if len(values) == 0 {
		return nil
	}

	valList := []any{}

	// UPDATE if where query
	if query.where != "" {
		q := `UPDATE ` + query.table + ` SET `
		for key, val := range values {
			q += toAlphaNumeric(key) + ` = ?, `
			valList = append(valList, val)
		}
		q = q[:len(q)-2]

		q += ` ` + query.where
		valList = append(valList, query.whereValue...)

		st, err := query.db.Prepare(q)
		if err != nil {
			return err
		}

		_, err = st.Exec(valList...)
		return err
	}

	// UPDATE if unique keys found with matching values
	if len(unique) != 0 {
		where := `WHERE `
		whereValue := []any{}

		// build where query
		hasVal := false
		for _, key := range unique {
			key = toAlphaNumeric(key)
			if val, ok := values[key]; ok {
				if hasVal {
					where += ` AND `
				}
				hasVal = true
				where += key + ` = ?`
				whereValue = append(whereValue, val)
			}
		}

		// check if table contains existing rows
		if hasVal {
			if rows, err := query.db.Query(`SELECT * FROM `+query.table+` `+where, whereValue...); err == nil && rows.Next() {
				rows.Close()

				// UPDATE values in existing rows
				q := `UPDATE ` + query.table + ` SET `
				for key, val := range values {
					q += toAlphaNumeric(key) + ` = ?, `
					valList = append(valList, val)
				}
				q = q[:len(q)-2]

				q += ` ` + where
				valList = append(valList, whereValue...)

				st, err := query.db.Prepare(q)
				if err != nil {
					return err
				}

				_, err = st.Exec(valList...)
				return err
			}
		}
	}

	// INSERT values into table
	qKey := ``
	qVal := ``
	for key, val := range values {
		qKey += toAlphaNumeric(key) + `, `
		qVal += `?, `
		valList = append(valList, val)
	}
	qKey = qKey[:len(qKey)-2]
	qVal = qVal[:len(qVal)-2]

	st, err := query.db.Prepare(`INSERT INTO ` + query.table + ` (` + qKey + `) VALUES (` + qVal + `)`)
	if err != nil {
		return err
	}
	st.Exec(valList...)

	return nil
}

// Delete will remove a row from the database table
//
// ! Warning: setting @force to true, will allow the database to delete all rows from a table, if its missing a `where` query
func (query *Query) Delete(force ...bool) error {
	if query.where == "" {
		if len(force) != 0 && force[0] {
			st, err := query.db.Prepare(`DELETE FROM ` + query.table)
			if err != nil {
				return err
			}
			st.Exec()

			return nil
		}

		return Error_UnsafeQuery
	}

	st, err := query.db.Prepare(`DELETE FROM ` + query.table + ` ` + query.where)
	if err != nil {
		return err
	}
	st.Exec(query.whereValue...)

	return nil
}

// Drop will drop an entire table from the database, deleting everything
//
// ! Warning: This method will Delete the entire table, and will ignore any `where` queries
func (query *Query) Drop(force bool) error {
	if !force {
		return Error_UnsafeQuery
	}

	// Note: query.db.SQL will bypass the default safety checks,
	// since the `DROP` keyword will be denied by safety checks.
	st, err := query.db.SQL.Prepare(`DROP TABLE ` + query.table)
	if err != nil {
		return err
	}
	st.Exec()

	return nil
}
