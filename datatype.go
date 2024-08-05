package gosql

import (
	"strconv"

	"github.com/tkdeng/go-sqlorm/common"
)

type DataType struct {
	key     string
	valType string
	def     string
}

// Default sets a DEFAULT value
func (dataType DataType) Default(value any) *DataType {
	//todo: add optional sql function methods (with custom struct)

	if dataType.valType == "string" {
		dataType.def = `'` + sqlEscapeQuote(common.ToType[string](value)) + `'`
	} else {
		dataType.def = common.ToType[string](value)
	}
	return &dataType
}

// Append allows you to add custom type constraints to a DataType
//
// You can use this if an SQL DataType constraint is not supported by this module.
func (dataType DataType) Append(val string) *DataType {
	dataType.key += ` ` + val
	return &dataType
}

// Unique sets a type to UNIQUE
func (dataType DataType) Unique() *DataType {
	dataType.key += ` UNIQUE`
	return &dataType
}

// NotNull makes a type NOT NULL
func (dataType DataType) NotNull() *DataType {
	dataType.key += ` NOT NULL`
	return &dataType
}

// AutoInc makes a type AUTO_INCREMENT
//
// Note: AUTO_INCREMENT may Not be supported by sqlite
func (dataType DataType) AutoInc() *DataType {
	dataType.key += ` AUTO_INCREMENT`
	return &dataType
}

// Primary makes a type a PRIMARY KEY
func (dataType DataType) Primary() *DataType {
	dataType.key += ` PRIMARY KEY`
	return &dataType
}

// TYPE is a Custom DataType
//
// You can use this if an SQL DataType is not supported by this module.
// A list of SQL DataTypes can be found here: https://www.w3schools.com/sql/sql_datatypes.asp
func TYPE(key string, dataType string) *DataType {
	return &DataType{key: toAlphaNumeric(key) + " " + dataType, valType: "custom"}
}

//* String Data Types

// CHAR is a String DataType
//
//	size: 0 to 255 (default: 1)
func CHAR(key string, size ...uint8) *DataType {
	t := toAlphaNumeric(key) + " CHAR"

	if len(size) != 0 {
		t += "(" + strconv.FormatUint(uint64(size[0]), 10) + ")"
	}

	return &DataType{key: t, valType: "string"}
}

// VARCHAR is a String DataType
//
//	size: 0 to 65535
func VARCHAR(key string, size ...uint16) *DataType {
	t := toAlphaNumeric(key) + " VARCHAR"

	if len(size) != 0 {
		t += "(" + strconv.FormatUint(uint64(size[0]), 10) + ")"
	}

	return &DataType{key: t, valType: "string"}
}

// BINARY is a String DataType
//
//	size: 0 to 255 (default: 1)
func BINARY(key string, size ...uint8) *DataType {
	t := toAlphaNumeric(key) + " BINARY"

	if len(size) != 0 {
		t += "(" + strconv.FormatUint(uint64(size[0]), 10) + ")"
	}

	return &DataType{key: t, valType: "string"}
}

// VARBINARY is a String DataType
//
//	size: 0 to 65535
func VARBINARY(key string, size ...uint16) *DataType {
	t := toAlphaNumeric(key) + " VARBINARY"

	if len(size) != 0 {
		t += "(" + strconv.FormatUint(uint64(size[0]), 10) + ")"
	}

	return &DataType{key: t, valType: "string"}
}

// TINYBLOB (Binary Large Objects) is a String DataType
//
//	size: 255
func TINYBLOB(key string) *DataType {
	t := toAlphaNumeric(key) + " TINYBLOB"

	return &DataType{key: t, valType: "string"}
}

// TINYTEXT is a String DataType
//
//	size: 255
func TINYTEXT(key string) *DataType {
	t := toAlphaNumeric(key) + " TINYTEXT"

	return &DataType{key: t, valType: "string"}
}

// TEXT is a String DataType
//
//	size: 0 to 65535
func TEXT(key string, size ...uint16) *DataType {
	t := toAlphaNumeric(key) + " TEXT"

	if len(size) != 0 {
		t += "(" + strconv.FormatUint(uint64(size[0]), 10) + ")"
	}

	return &DataType{key: t, valType: "string"}
}

// BLOB (Binary Large Objects) is a String DataType
//
//	size: 0 to 65535
func BLOB(key string, size ...uint16) *DataType {
	t := toAlphaNumeric(key) + " BLOB"

	if len(size) != 0 {
		t += "(" + strconv.FormatUint(uint64(size[0]), 10) + ")"
	}

	return &DataType{key: t, valType: "string"}
}

// MEDIUMTEXT is a String DataType
//
//	size: 16777215
func MEDIUMTEXT(key string) *DataType {
	t := toAlphaNumeric(key) + " MEDIUMTEXT"

	return &DataType{key: t, valType: "string"}
}

// MEDIUMBLOB (Binary Large Objects) is a String DataType
//
//	size: 16777215
func MEDIUMBLOB(key string) *DataType {
	t := toAlphaNumeric(key) + " MEDIUMBLOB"

	return &DataType{key: t, valType: "string"}
}

// LONGTEXT is a String DataType
//
//	size: 4294967295
func LONGTEXT(key string) *DataType {
	t := toAlphaNumeric(key) + " LONGTEXT"

	return &DataType{key: t, valType: "string"}
}

// LONGBLOB (Binary Large Objects) is a String DataType
//
//	size: 4294967295
func LONGBLOB(key string) *DataType {
	t := toAlphaNumeric(key) + " LONGBLOB"

	return &DataType{key: t, valType: "string"}
}

// ENUM is a String DataType
//
// A string object that can have only one value, chosen from a list of possible values.
// You can list up to 65535 values in an ENUM list. If a value is inserted that is not
// in the list, a blank value will be inserted. The values are sorted in the order you
// enter them.
func ENUM(key string, val ...string) *DataType {
	t := toAlphaNumeric(key) + " ENUM"

	if len(val) != 0 {
		t += "("
		for i, v := range val {
			t += toAlphaNumeric(v)

			if i >= 65535 {
				break
			}

			if i != len(val)-1 {
				t += ", "
			}
		}
		t += ")"
	}

	return &DataType{key: t, valType: "string"}
}

// SET is a String DataType
//
// A string object that can have 0 or more values, chosen from a list of possible values.
// You can list up to 64 values in a SET list.
func SET(key string, val ...string) *DataType {
	t := toAlphaNumeric(key) + " SET"

	if len(val) != 0 {
		t += "("
		for i, v := range val {
			t += toAlphaNumeric(v)

			if i >= 64 {
				break
			}

			if i != len(val)-1 {
				t += ", "
			}
		}
		t += ")"
	}

	return &DataType{key: t, valType: "string"}
}

//* Numeric Data Types

// BIT is a Numeric DataType
//
//	size: 1 to 64 (default: 1)
func BIT(key string, size ...uint8) *DataType {
	t := toAlphaNumeric(key) + " BIT"

	if len(size) != 0 {
		if size[0] > 64 {
			size[0] = 64
		}

		t += "(" + strconv.FormatUint(uint64(size[0]), 10) + ")"
	}

	return &DataType{key: t, valType: "numeric"}
}

// BOOL is a Numeric DataType
//
//	0 = false | 1 = true
func BOOL(key string) *DataType {
	t := toAlphaNumeric(key) + " BOOL"

	return &DataType{key: t, valType: "numeric"}
}

// TINYINT is a Numeric DataType
//
//	size: -128 to 127 | 0 to 255
func TINYINT(key string) *DataType {
	t := toAlphaNumeric(key) + " TINYINT"

	return &DataType{key: t, valType: "numeric"}
}

// SMALLINT is a Numeric DataType
//
//	size: -32768 to 32767 | 0 to 65535
func SMALLINT(key string) *DataType {
	t := toAlphaNumeric(key) + " SMALLINT"

	return &DataType{key: t, valType: "numeric"}
}

// MEDIUMINT is a Numeric DataType
//
//	size: -8388608 to 8388607 | 0 to 16777215
func MEDIUMINT(key string) *DataType {
	t := toAlphaNumeric(key) + " MEDIUMINT"

	return &DataType{key: t, valType: "numeric"}
}

// INT is a Numeric DataType
//
//	size: -2147483648 to 2147483647 | 0 to 4294967295
func INT(key string) *DataType {
	t := toAlphaNumeric(key) + " INT"

	return &DataType{key: t, valType: "numeric"}
}

// BIGINT is a Numeric DataType
//
//	size: -9223372036854775808 to 9223372036854775807 | 0 to 18446744073709551615
func BIGINT(key string) *DataType {
	t := toAlphaNumeric(key) + " BIGINT"

	return &DataType{key: t, valType: "numeric"}
}

// FLOAT is a Numeric DataType
//
// A floating point number. MySQL uses the p value to determine whether to use FLOAT or
// DOUBLE for the resulting data type. If p is from 0 to 24, the data type becomes FLOAT().
// If p is from 25 to 53, the data type becomes DOUBLE().
func FLOAT(key string, p ...uint8) *DataType {
	t := toAlphaNumeric(key) + " FLOAT"

	if len(p) != 0 {
		if p[0] > 53 {
			p[0] = 53
		}

		t += "(" + strconv.FormatUint(uint64(p[0]), 10) + ")"
	}

	return &DataType{key: t, valType: "numeric"}
}

// DOUBLE is a Numeric DataType
//
// A normal-size floating point number. The total number of digits is specified in size.
// The number of digits after the decimal point is specified in the d parameter.
func DOUBLE(key string, sizeD ...uint8) *DataType {
	t := toAlphaNumeric(key) + " DOUBLE"

	if len(sizeD) != 0 {
		t += "(" + strconv.FormatUint(uint64(sizeD[0]), 10)

		if len(sizeD) > 1 {
			t += ", " + strconv.FormatUint(uint64(sizeD[1]), 10)
		}

		t += ")"
	}

	return &DataType{key: t, valType: "numeric"}
}

// DECIMAL is a Numeric DataType
//
// An exact fixed-point number. The total number of digits is specified in size. The number
// of digits after the decimal point is specified in the d parameter. The maximum number for
// size is 65. The maximum number for d is 30. The default value for size is 10. The default
// value for d is 0.
func DECIMAL(key string, sizeD ...uint8) *DataType {
	t := toAlphaNumeric(key) + " DECIMAL"

	if len(sizeD) != 0 {
		if sizeD[0] > 65 {
			sizeD[0] = 65
		}

		t += "(" + strconv.FormatUint(uint64(sizeD[0]), 10)

		if len(sizeD) > 1 {
			if sizeD[1] > 30 {
				sizeD[1] = 30
			}

			t += ", " + strconv.FormatUint(uint64(sizeD[1]), 10)
		}

		t += ")"
	}

	return &DataType{key: t, valType: "numeric"}
}

//* Date and Time Data Types

// DATE is a DateTime DataType
//
// A date. Format: YYYY-MM-DD. The supported range is from '1000-01-01' to '9999-12-31'.
func DATE(key string) *DataType {
	t := toAlphaNumeric(key) + " DATE"

	return &DataType{key: t, valType: "datetime"}
}

// DATETIME is a DateTime DataType
//
// A date and time combination. Format: YYYY-MM-DD hh:mm:ss. The supported range is from
// '1000-01-01 00:00:00' to '9999-12-31 23:59:59'. Adding DEFAULT and ON UPDATE in the column
// definition to get automatic initialization and updating to the current date and time.
func DATETIME(key string, fsp ...string) *DataType {
	t := toAlphaNumeric(key) + " DATETIME"

	if len(fsp) != 0 {
		t += "('" + sqlEscapeQuote(fsp[0]) + "')"
	}

	return &DataType{key: t, valType: "datetime"}
}

// TIMESTAMP is a DateTime DataType
//
// A timestamp. TIMESTAMP values are stored as the number of seconds since the Unix epoch
// ('1970-01-01 00:00:00' UTC). Format: YYYY-MM-DD hh:mm:ss. The supported range is from
// '1970-01-01 00:00:01' UTC to '2038-01-09 03:14:07' UTC. Automatic initialization and
// updating to the current date and time can be specified using DEFAULT CURRENT_TIMESTAMP
// and ON UPDATE CURRENT_TIMESTAMP in the column definition.
func TIMESTAMP(key string, fsp ...string) *DataType {
	t := toAlphaNumeric(key) + " TIMESTAMP"

	if len(fsp) != 0 {
		t += "('" + sqlEscapeQuote(fsp[0]) + "')"
	}

	return &DataType{key: t, valType: "datetime"}
}

// TIME is a DateTime DataType
//
// A time. Format: hh:mm:ss. The supported range is from '-838:59:59' to '838:59:59'.
func TIME(key string, fsp ...string) *DataType {
	t := toAlphaNumeric(key) + " TIME"

	if len(fsp) != 0 {
		t += "('" + sqlEscapeQuote(fsp[0]) + "')"
	}

	return &DataType{key: t, valType: "datetime"}
}

// YEAR is a DateTime DataType
//
// A year in four-digit format. Values allowed in four-digit format: 1901 to 2155, and 0000.
func YEAR(key string) *DataType {
	t := toAlphaNumeric(key) + " YEAR"

	return &DataType{key: t, valType: "datetime"}
}
