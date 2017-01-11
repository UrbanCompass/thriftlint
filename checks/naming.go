package checks

import (
	"reflect"
	"regexp"
	"strings"

	"github.com/UrbanCompass/thriftlint"
)

var (
	upperCamelCaseRegex = `^[_A-Z][a-z]*([A-Z][0-9a-z]*)*$`
	lowerCamelCaseRegex = `^[_a-z]+([A-Z0-9a-z]*)*$`
	upperSnakeCaseRegex = `^[A-Z_]+([A-Z0-9]+_?)*$`
)

var (
	// CheckNamesDefaults is a map of Thrift AST node type to a regular expression for
	// validating names of that type.
	CheckNamesDefaults = map[reflect.Type]string{
		thriftlint.ServiceType:   upperCamelCaseRegex,
		thriftlint.EnumType:      upperCamelCaseRegex,
		thriftlint.StructType:    upperCamelCaseRegex,
		thriftlint.EnumValueType: upperSnakeCaseRegex,
		thriftlint.FieldType:     lowerCamelCaseRegex,
		thriftlint.MethodType:    lowerCamelCaseRegex,
		thriftlint.ConstantType:  upperSnakeCaseRegex,
	}

	// CheckNamesDefaultBlacklist is names that should never be used for symbols.
	CheckNamesDefaultBlacklist = map[string]bool{
		"class": true,
		"int":   true,
	}
)

// CheckNames checks Thrift symbols comply with a set of regular expressions.
//
// If mathces or blacklist are nil, global defaults will be used.
func CheckNames(matches map[reflect.Type]string, blacklist map[string]bool) thriftlint.Check {
	if matches == nil {
		matches = CheckNamesDefaults
	}
	if blacklist == nil {
		blacklist = CheckNamesDefaultBlacklist
	}
	regexes := map[reflect.Type]*regexp.Regexp{}
	for t, p := range matches {
		regexes[t] = regexp.MustCompile(p)
	}
	return thriftlint.MakeCheck("naming", func(v interface{}) (messages thriftlint.Messages) {
		rv := reflect.Indirect(reflect.ValueOf(v))
		nameField := rv.FieldByName("Name")
		if !nameField.IsValid() {
			return nil
		}
		name := nameField.Interface().(string)
		// Special-case DEPRECATED_ fields.
		checker, ok := regexes[rv.Type()]
		if !ok || strings.HasPrefix(name, "DEPRECATED_") {
			return nil
		}
		if blacklist[name] {
			messages.Warning(v, "%q is a disallowed name", name)
		}
		if ok := checker.MatchString(name); !ok {
			messages.Warning(v, "name of %s %q should match %q", strings.ToLower(rv.Type().Name()),
				name, checker.String())
		}
		return
	})
}
