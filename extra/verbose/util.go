package verbose

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"encoding/hex"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unicode"

	bundialect "github.com/uptrace/bun/dialect"
	"github.com/uptrace/bun/schema"
)

// Sprintf will interpolate SQL args into a query string containing prepared
// statement parameters. It returns an error if an argument cannot be properly
// represented in SQL. This function may be vulnerable to SQL injection and
// should be used for logging purposes only.
func Sprintf(dialect schema.Dialect, query string, args []any) (string, error) {
	if len(args) == 0 {
		return query, nil
	}
	buf := bufPool.Get().(*bytes.Buffer)
	buf.Reset()
	defer bufPool.Put(buf)
	buf.Grow(len(query))
	namedIndices := make(map[string]int)
	for i, arg := range args {
		switch arg := arg.(type) {
		case sql.NamedArg:
			namedIndices[arg.Name] = i
		}
	}
	runningArgsIndex := 0
	mustWriteCharAt := -1
	insideStringOrIdentifier := false
	var openingQuote rune
	var paramName []rune
	for i, char := range query {
		// do we unconditionally write in the current char?
		if mustWriteCharAt == i {
			buf.WriteRune(char)
			continue
		}
		// are we currently inside a string or identifier?
		if insideStringOrIdentifier {
			buf.WriteRune(char)
			switch openingQuote {
			case '\'', '"', '`':
				// does the current char terminate the current string or identifier?
				if char == openingQuote {
					// is the next char the same as the current char, which
					// escapes it and prevents it from terminating the current
					// string or identifier?
					if i+1 < len(query) && rune(query[i+1]) == openingQuote {
						mustWriteCharAt = i + 1
					} else {
						insideStringOrIdentifier = false
					}
				}
			case '[':
				// does the current char terminate the current string or identifier?
				if char == ']' {
					// is the next char the same as the current char, which
					// escapes it and prevents it from terminating the current
					// string or identifier?
					if i+1 < len(query) && query[i+1] == ']' {
						mustWriteCharAt = i + 1
					} else {
						insideStringOrIdentifier = false
					}
				}
			}
			continue
		}
		// does the current char mark the start of a new string or identifier?
		if char == '\'' || char == '"' || (char == '`' && dialect.Name() == bundialect.MySQL) || (char == '[' && dialect.Name() == bundialect.MSSQL) {
			insideStringOrIdentifier = true
			openingQuote = char
			buf.WriteRune(char)
			continue
		}
		// are we currently inside a parameter name?
		if len(paramName) > 0 {
			// does the current char terminate the current parameter name?
			if char != '_' && !unicode.IsLetter(char) && !unicode.IsDigit(char) {
				paramValue, err := lookupParam(dialect, args, paramName, namedIndices, runningArgsIndex)
				if err != nil {
					return buf.String(), err
				}
				buf.WriteString(paramValue)
				buf.WriteRune(char)
				if len(paramName) == 1 && paramName[0] == '?' {
					runningArgsIndex++
				}
				paramName = paramName[:0]
			} else {
				paramName = append(paramName, char)
			}
			continue
		}
		// does the current char mark the start of a new parameter name?
		if (char == '$' && (dialect.Name() == bundialect.SQLite || dialect.Name() == bundialect.PG)) ||
			(char == ':' && dialect.Name() == bundialect.SQLite) ||
			(char == '@' && (dialect.Name() == bundialect.SQLite || dialect.Name() == bundialect.MSSQL)) {
			paramName = append(paramName, char)
			continue
		}
		// is the current char the anonymous '?' parameter?
		if char == '?' && dialect.Name() != bundialect.PG {
			// for sqlite, just because we encounter a '?' doesn't mean it
			// is an anonymous param. sqlite also supports using '?' for
			// ordinal params (e.g. ?1, ?2, ?3) or named params (?foo,
			// ?bar, ?baz). Hence we treat it as an ordinal/named param
			// first, and handle the edge case later when it isn't.
			if dialect.Name() == bundialect.SQLite {
				paramName = append(paramName, char)
				continue
			}
			if runningArgsIndex >= len(args) {
				return buf.String(), fmt.Errorf("too few args provided, expected more than %d", runningArgsIndex+1)
			}
			paramValue, err := Sprint(dialect, args[runningArgsIndex])
			if err != nil {
				return buf.String(), err
			}
			buf.WriteString(paramValue)
			runningArgsIndex++
			continue
		}
		// if all the above questions answer false, we just write the current
		// char in and continue
		buf.WriteRune(char)
	}
	// flush the paramName buffer (to handle edge case where the query ends with a parameter name)
	if len(paramName) > 0 {
		paramValue, err := lookupParam(dialect, args, paramName, namedIndices, runningArgsIndex)
		if err != nil {
			return buf.String(), err
		}
		buf.WriteString(paramValue)
	}
	if insideStringOrIdentifier {
		return buf.String(), fmt.Errorf("unclosed string or identifier")
	}
	return buf.String(), nil
}

// Sprint is the equivalent of Sprintf but for converting a single value into
// its SQL representation.
func Sprint(dialect schema.Dialect, v any) (string, error) {
	const (
		timestamp             = "2006-01-02 15:04:05"
		timestampWithTimezone = "2006-01-02 15:04:05.9999999-07:00"
	)
	switch v := v.(type) {
	case nil:
		return "NULL", nil
	case bool:
		if v {
			if dialect.Name() == bundialect.MSSQL {
				return "1", nil
			}
			return "TRUE", nil
		}
		if dialect.Name() == bundialect.MSSQL {
			return "0", nil
		}
		return "FALSE", nil
	case []byte:
		switch dialect.Name() {
		case bundialect.PG:
			// https://www.postgresql.org/docs/current/datatype-binary.html
			// (see 8.4.1. bytea Hex Format)
			return `'\x` + hex.EncodeToString(v) + `'`, nil
		case bundialect.MSSQL:
			return `0x` + hex.EncodeToString(v), nil
		default:
			return `x'` + hex.EncodeToString(v) + `'`, nil
		}
	case string:
		str := v
		i := strings.IndexAny(str, "\r\n")
		if i < 0 {
			return `'` + strings.ReplaceAll(str, `'`, `''`) + `'`, nil
		}
		var b strings.Builder
		if dialect.Name() == bundialect.MySQL || dialect.Name() == bundialect.MSSQL {
			b.WriteString("CONCAT(")
		}
		for i >= 0 {
			if str[:i] != "" {
				b.WriteString(`'` + strings.ReplaceAll(str[:i], `'`, `''`) + `'`)
				if dialect.Name() == bundialect.MSSQL || dialect.Name() == bundialect.MSSQL {
					b.WriteString(", ")
				} else {
					b.WriteString(" || ")
				}
			}
			switch str[i] {
			case '\r':
				if dialect.Name() == bundialect.PG {
					b.WriteString("CHR(13)")
				} else {
					b.WriteString("CHAR(13)")
				}
			case '\n':
				if dialect.Name() == bundialect.PG {
					b.WriteString("CHR(10)")
				} else {
					b.WriteString("CHAR(10)")
				}
			}
			if str[i+1:] != "" {
				if dialect.Name() == bundialect.MySQL || dialect.Name() == bundialect.MSSQL {
					b.WriteString(", ")
				} else {
					b.WriteString(" || ")
				}
			}
			str = str[i+1:]
			i = strings.IndexAny(str, "\r\n")
		}
		if str != "" {
			b.WriteString(`'` + strings.ReplaceAll(str, `'`, `''`) + `'`)
		}
		if dialect.Name() == bundialect.MSSQL || dialect.Name() == bundialect.MSSQL {
			b.WriteString(")")
		}
		return b.String(), nil
	case time.Time:
		if dialect.Name() == bundialect.PG || dialect.Name() == bundialect.MSSQL {
			return `'` + v.Format(timestampWithTimezone) + `'`, nil
		}
		return `'` + v.UTC().Format(timestamp) + `'`, nil
	case int:
		return strconv.FormatInt(int64(v), 10), nil
	case int8:
		return strconv.FormatInt(int64(v), 10), nil
	case int16:
		return strconv.FormatInt(int64(v), 10), nil
	case int32:
		return strconv.FormatInt(int64(v), 10), nil
	case int64:
		return strconv.FormatInt(v, 10), nil
	case uint:
		return strconv.FormatUint(uint64(v), 10), nil
	case uint8:
		return strconv.FormatUint(uint64(v), 10), nil
	case uint16:
		return strconv.FormatUint(uint64(v), 10), nil
	case uint32:
		return strconv.FormatUint(uint64(v), 10), nil
	case uint64:
		return strconv.FormatUint(v, 10), nil
	case float32:
		return strconv.FormatFloat(float64(v), 'g', -1, 64), nil
	case float64:
		return strconv.FormatFloat(v, 'g', -1, 64), nil
	case sql.NamedArg:
		return Sprint(dialect, v.Value)
	case sql.NullBool:
		if !v.Valid {
			return "NULL", nil
		}
		if v.Bool {
			if dialect.Name() == bundialect.MSSQL {
				return "1", nil
			}
			return "TRUE", nil
		}
		if dialect.Name() == bundialect.MSSQL {
			return "0", nil
		}
		return "FALSE", nil
	case sql.NullFloat64:
		if !v.Valid {
			return "NULL", nil
		}
		return strconv.FormatFloat(v.Float64, 'g', -1, 64), nil
	case sql.NullInt64:
		if !v.Valid {
			return "NULL", nil
		}
		return strconv.FormatInt(v.Int64, 10), nil
	case sql.NullInt32:
		if !v.Valid {
			return "NULL", nil
		}
		return strconv.FormatInt(int64(v.Int32), 10), nil
	case sql.NullString:
		if !v.Valid {
			return "NULL", nil
		}
		return Sprint(dialect, v.String)
	case sql.NullTime:
		if !v.Valid {
			return "NULL", nil
		}
		if dialect.Name() == bundialect.PG || dialect.Name() == bundialect.MSSQL {
			return `'` + v.Time.Format(timestampWithTimezone) + `'`, nil
		}
		return `'` + v.Time.UTC().Format(timestamp) + `'`, nil
	case driver.Valuer:
		vv, err := v.Value()
		if err != nil {
			return "", fmt.Errorf("error when calling Value(): %w", err)
		}
		switch vv.(type) {
		case int64, float64, bool, []byte, string, time.Time, nil:
			return Sprint(dialect, vv)
		default:
			return "", fmt.Errorf("invalid driver.Value type %T (must be one of int64, float64, bool, []byte, string, time.Time, nil)", vv)
		}
	}
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Pointer {
		rv = rv.Elem()
		if !rv.IsValid() {
			return "NULL", nil
		}
	}
	switch v := rv.Interface().(type) {
	case bool, []byte, string, time.Time, int, int8, int16, int32, int64, uint,
		uint8, uint16, uint32, uint64, float32, float64, sql.NamedArg,
		sql.NullBool, sql.NullFloat64, sql.NullInt64, sql.NullInt32,
		sql.NullString, sql.NullTime, driver.Valuer:
		return Sprint(dialect, v)
	default:
		return "", fmt.Errorf("%T has no SQL representation", v)
	}
}

// lookupParam returns the SQL representation of a paramName (inside the args
// slice).
func lookupParam(dialect schema.Dialect, args []any, paramName []rune, namedIndices map[string]int, runningArgsIndex int) (paramValue string, err error) {
	var maybeNum string
	if paramName[0] == '@' && dialect.Name() == bundialect.MSSQL && len(paramName) >= 2 && (paramName[1] == 'p' || paramName[1] == 'P') {
		maybeNum = string(paramName[2:])
	} else {
		maybeNum = string(paramName[1:])
	}

	// is paramName an anonymous parameter?
	if maybeNum == "" {
		if paramName[0] != '?' {
			return "", fmt.Errorf("parameter name missing")
		}
		paramValue, err = Sprint(dialect, args[runningArgsIndex])
		if err != nil {
			return "", err
		}
		return paramValue, nil
	}

	// is paramName an ordinal paramater?
	ordinal, err := strconv.Atoi(maybeNum)
	if err == nil {
		index := ordinal - 1
		if index < 0 || index >= len(args) {
			return "", fmt.Errorf("args index %d out of bounds", ordinal)
		}
		paramValue, err = Sprint(dialect, args[index])
		if err != nil {
			return "", err
		}
		return paramValue, nil
	}

	// if we reach here, we know that the paramName is not an ordinal parameter
	// i.e. it is a named parameter
	if dialect.Name() == bundialect.PG || dialect.Name() == bundialect.MySQL {
		return "", fmt.Errorf("%s does not support %s named parameter", dialect, string(paramName))
	}
	index, ok := namedIndices[string(paramName[1:])]
	if !ok {
		return "", fmt.Errorf("named parameter %s not provided", string(paramName))
	}
	if index < 0 || index >= len(args) {
		return "", fmt.Errorf("args index %d out of bounds", ordinal)
	}
	paramValue, err = Sprint(dialect, args[index])
	if err != nil {
		return "", err
	}
	return paramValue, nil
}
