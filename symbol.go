package thriftlint

import (
	"bytes"
	"fmt"
	"go/doc"
	"reflect"
	"strings"
	"unicode"
)

var (
	// BuiltinThriftTypes is a map of the basic builtin Thrift types. Useful in templates.
	BuiltinThriftTypes = map[string]bool{
		"bool":   true,
		"byte":   true,
		"i16":    true,
		"i32":    true,
		"i64":    true,
		"double": true,
		"string": true,
	}

	// BuiltinThriftCollections is the set of builtin collection types in Thrift.
	BuiltinThriftCollections = map[string]bool{
		"map":    true,
		"list":   true,
		"set":    true,
		"binary": true,
	}
)

// Scanner over a sequence of runes.
type scanner struct {
	runes  []rune
	cursor int
}

func (s *scanner) peek() rune {
	if s.cursor >= len(s.runes) {
		return -1
	}
	return s.runes[s.cursor]
}

func (s *scanner) next() rune {
	r := s.peek()
	if r != -1 {
		s.cursor++
	}
	return r
}

func (s *scanner) reverse() {
	s.cursor--
}

func consumeLower(scan *scanner) string {
	out := ""
	for unicode.IsLower(scan.peek()) || unicode.IsNumber(scan.peek()) {
		out += string(scan.next())
	}
	return out
}

// ([A-Z]+)(?:[A-Z][a-z]|$)
func consumeMostUpper(scan *scanner) string {
	out := ""
	for unicode.IsUpper(scan.peek()) || unicode.IsNumber(scan.peek()) {
		r := scan.next()
		if unicode.IsLower(scan.peek()) && !commonInitialisms[out+string(r)] {
			scan.reverse()
			break
		}
		out += string(r)
	}
	return out
}

func title(s string) string {
	return strings.ToUpper(s[0:1]) + strings.ToLower(s[1:])
}

// From https://github.com/golang/lint/blob/master/lint.go
var commonInitialisms = map[string]bool{
	"API":   true,
	"ASCII": true,
	"CPU":   true,
	"CSS":   true,
	"DB":    true,
	"DNS":   true,
	"EOF":   true,
	"GUID":  true,
	"HTML":  true,
	"HTTP":  true,
	"HTTPS": true,
	"ID":    true,
	"IP":    true,
	"JSON":  true,
	"LHS":   true,
	"MD5":   true,
	"MLS":   true,
	"OK":    true,
	"QPS":   true,
	"RAM":   true,
	"RHS":   true,
	"RPC":   true,
	"SHA":   true,
	"SLA":   true,
	"SMTP":  true,
	"SQL":   true,
	"SSH":   true,
	"TCP":   true,
	"TLS":   true,
	"TTL":   true,
	"UDP":   true,
	"UI":    true,
	"UID":   true,
	"URI":   true,
	"URL":   true,
	"UTC":   true,
	"UTF8":  true,
	"UUID":  true,
	"VM":    true,
	"XML":   true,
	"XSRF":  true,
	"XSS":   true,
}

func Comment(v interface{}) []string {
	comment := reflect.Indirect(reflect.ValueOf(v)).FieldByName("Comment").Interface().(string)
	if comment == "" {
		return nil
	}
	w := bytes.NewBuffer(nil)
	doc.ToText(w, comment, "", "", 80)
	return strings.Split(strings.TrimSpace(w.String()), "\n")
}

func IsInitialism(s string) bool {
	return commonInitialisms[strings.ToUpper(s)]
}

// UpperCamelCase converts a symbol to CamelCase
func UpperCamelCase(s string) string {
	parts := []string{}
	for _, part := range SplitSymbol(s) {
		if part == "" {
			parts = append(parts, "_")
			continue
		}
		if part == "s" && len(parts) > 0 {
			parts[len(parts)-1] += part
		} else {
			if commonInitialisms[strings.ToUpper(part)] {
				part = strings.ToUpper(part)
			} else {
				part = title(part)
			}
			parts = append(parts, part)
		}
	}
	return strings.Join(parts, "")
}

// LowerCamelCase converts a symbol to lowerCamelCase
func LowerCamelCase(s string) string {
	first := true
	parts := []string{}
	for _, part := range SplitSymbol(s) {
		if part == "" {
			parts = append(parts, "_")
			continue
		}
		if first {
			parts = append(parts, strings.ToLower(part))
			first = false
		} else {
			// Merge trailing s
			if part == "s" && len(parts) > 0 {
				parts[len(parts)-1] += part
			} else {
				if commonInitialisms[strings.ToUpper(part)] {
					part = strings.ToUpper(part)
				} else {
					part = title(part)
				}
				parts = append(parts, part)
			}
		}
	}
	return strings.Join(parts, "")
}

// LowerSnakeCase converts a symbol to snake_case
func LowerSnakeCase(s string) string {
	parts := []string{}
	for _, part := range SplitSymbol(s) {
		if part == "" {
			parts = append(parts, "_")
			continue
		}
		parts = append(parts, strings.ToLower(part))
	}
	return strings.Join(parts, "_")
}

// UpperSnakeCase converts a symbol to UPPER_SNAKE_CASE
func UpperSnakeCase(s string) string {
	parts := []string{}
	for _, part := range SplitSymbol(s) {
		if part == "" {
			parts = append(parts, "_")
			continue
		}
		parts = append(parts, strings.ToUpper(part))
	}
	return strings.Join(parts, "_")
}

// SplitSymbol splits an arbitrary symbol into parts. It accepts symbols in snake case and camel
// case, and correctly supports all-caps substrings.
//
// eg. "some_snake_case_symbol" would become ["some", "snake", "case", "symbol"]
// and "someCamelCaseSymbol" would become ["some", "Camel", "Case", "Symbol"]
func SplitSymbol(s string) []string {
	// This is painful. See TestSplitSymbol for examples of what this does.
	out := []string{}
	scan := &scanner{runes: []rune(s)}
	for scan.peek() != -1 {
		part := ""
		r := scan.peek()
		switch {
		case unicode.IsLower(r):
			part = consumeLower(scan)
		case unicode.IsUpper(r):
			scan.next()
			// [A-Z][a-z]+
			if unicode.IsLower(scan.peek()) {
				part += string(r)
				part += consumeLower(scan)
			} else {
				scan.reverse()
				part += consumeMostUpper(scan)
			}
		case unicode.IsNumber(r):
			for unicode.IsNumber(scan.peek()) {
				part += string(scan.next())
			}
		case r == '_':
			scan.next()
			if len(out) == 0 {
				break
			}
			continue
		default:
			panic(fmt.Sprintf("unsupported character %q in %q", r, s))
		}
		out = append(out, part)
	}
	return out
}

// Extract the suffix from a . separated string (ie. namespace).
// Useful for getting the package reference from a files namespace.
func DotSuffix(pkg string) string {
	parts := strings.Split(pkg, ".")
	return parts[len(parts)-1]
}
