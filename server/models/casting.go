package models

import (
	"database/sql"
	"reflect"
	"time"
)

func floatToBool(val float64) bool {
	if val == float64(0) {
		return false
	} else {
		return true
	}
}

func floatToInt(val float64) int {
	return int(val)
}

func nilToEmptyString(val interface{}) string {
	if val == nil {
		return ""
	} else {
		return val.(string)
	}
}

func nilToInt(val interface{}) int {
	if val == nil {
		return 0
	} else {
		return val.(int)
	}
}

func toNullString(val interface{}) sql.NullString {
	if val == nil {
		return sql.NullString{"", false}
	} else {
		return sql.NullString{val.(string), true}
	}
}

func toNullInt32(val interface{}) sql.NullInt32 {
	if val == nil {
		return sql.NullInt32{0, false}
	} else {
		return sql.NullInt32{val.(int32), true}
	}
}

func nullInt32ToJsonFormat(val interface{}) string {
	if reflect.TypeOf(val) == nil {
		return "null"
	} else {
		return "0"
	}
}

func boolToInt(val interface{}) int {
	if val == false {
		return 0
	} else {
		return 1
	}
}

func nullStringToJsonFormat(val interface{}) string {
	if reflect.TypeOf(val) == nil {
		return "null"
	} else {
		return "\"" + val.(string) + "\""
	}
}

func sqlNullStringToJsonFormat(val sql.NullString) string {
	if !val.Valid {
		return "null"
	} else {
		return "\"" + val.String + "\""
	}
}

func stringToMssqlDateTime(val string) (string, error) {
	formatOld := "06/01/02 15:04:05"
	formatNew := "2006-01-02 15:04:05"
	res, err := time.Parse(formatOld, val)
	return res.Format(formatNew), err
}

func sqlNullStringToMssqlDateTime(val sql.NullString) (string, error) {
	formatOld := "06/01/02 15:04:05"
	formatNew := "2006-01-02 15:04:05"
	if !val.Valid {
		return "null", nil
	} else {
		res, err := time.Parse(formatOld, val.String)
		return "'" + res.Format(formatNew) + "'", err
	}
}

func nullTimeToJsonFormat(val interface{}) string {
	if reflect.TypeOf(val) == nil {
		return "null"
	} else {
		return "\"" + val.(time.Time).Format("2006-01-02 15:04:05") + "\""
	}
}
