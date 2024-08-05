package gosql

import (
	"strconv"
)

type Query struct {
	db    *DB
	table string

	where      string
	whereValue []any
	order      string
}

func (query Query) OrderBy(key string, desc ...bool) *Query {
	if query.order == "" {
		query.order = "ORDER BY "
	} else {
		query.order += ", "
	}

	query.order += toAlphaNumeric(key)

	if len(desc) != 0 && desc[0] {
		query.order += ` DESC`
	} else {
		query.order += ` ASC`
	}

	return &query
}

type whereQuery struct {
	query Query
	where string
}

func (query Query) Where(key string, not ...bool) *whereQuery {
	q := `WHERE `
	if query.where != `` {
		q = ` AND `
	}

	if len(not) != 0 && not[0] {
		q += `NOT `
	}
	q += toAlphaNumeric(key)

	return &whereQuery{
		query: query,
		where: q,
	}
}

func (query Query) And(key string, not ...bool) *whereQuery {
	q := ` AND `
	if query.where == `` {
		q = `WHERE `
	}

	if len(not) != 0 && not[0] {
		q += `NOT `
	}
	q += toAlphaNumeric(key)

	return &whereQuery{
		query: query,
		where: q,
	}
}

func (query Query) Or(key string, not ...bool) *whereQuery {
	q := ` OR `
	if query.where == `` {
		q = `WHERE `
	}

	if len(not) != 0 && not[0] {
		q += `NOT `
	}
	q += toAlphaNumeric(key)

	return &whereQuery{
		query: query,
		where: q,
	}
}

// Equal `=`
func (query whereQuery) Equal(value any) *Query {
	query.query.where += query.where + ` = ?`
	query.query.whereValue = append(query.query.whereValue, value)
	return &query.query
}

// Not Equal `<>` || `!=`
func (query whereQuery) NotEqual(value any) *Query {
	query.query.where += query.where + ` <> ?`
	query.query.whereValue = append(query.query.whereValue, value)
	return &query.query
}

// Like
func (query whereQuery) Like(value any) *Query {
	query.query.where += query.where + ` LIKE ?`
	query.query.whereValue = append(query.query.whereValue, value)
	return &query.query
}

// In
func (query whereQuery) In(values ...any) *Query {
	if len(values) == 0 {
		return &query.query
	}

	query.query.where += query.where + ` IN (`
	for i, val := range values {
		query.query.where += `?`
		query.query.whereValue = append(query.query.whereValue, val)
		if i != len(values)-1 {
			query.query.where += `,`
		}
	}
	query.query.where += `)`

	return &query.query
}

// Greater Than `>`
func (query whereQuery) GreaterThan(value int) *Query {
	query.query.where += query.where + ` > ` + strconv.Itoa(value)
	return &query.query
}

// Less Than `<`
func (query whereQuery) LessThan(value int) *Query {
	query.query.where += query.where + ` < ` + strconv.Itoa(value)
	return &query.query
}

// Greater Than or Equal `>=`
func (query whereQuery) GreaterEqual(value int) *Query {
	query.query.where += query.where + ` >= ` + strconv.Itoa(value)
	return &query.query
}

// Less Than or Equal `<=`
func (query whereQuery) LessEqual(value int) *Query {
	query.query.where += query.where + ` <= ` + strconv.Itoa(value)
	return &query.query
}

// Between
func (query whereQuery) Between(value1 int, value2 int) *Query {
	query.query.where += query.where + ` BETWEEN ` + strconv.Itoa(value1) + ` AND ` + strconv.Itoa(value2)
	return &query.query
}

// IsNull
func (query whereQuery) IsNull() *Query {
	query.query.where += query.where + ` IS NULL `
	return &query.query
}

// IsNotNull
func (query whereQuery) IsNotNull() *Query {
	query.query.where += query.where + ` IS NOT NULL `
	return &query.query
}
