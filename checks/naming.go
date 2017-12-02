package checks

import (
	"reflect"
	"regexp"
	"strings"

	"github.com/UrbanCompass/thriftlint"
	yaml "gopkg.in/yaml.v2"
)

type NamingStyle struct {
	Name    string
	Pattern *regexp.Regexp
}

var (
	upperCamelCaseStyle = NamingStyle{
		Name:    "title case",
		Pattern: regexp.MustCompile(`^_?([A-Z][0-9a-z]*)*$`),
	}
	lowerCamelCaseStyle = NamingStyle{
		Name:    "camel case",
		Pattern: regexp.MustCompile(`^_?[a-z][A-Z0-9a-z]*$`),
	}
	upperSnakeCaseStyle = NamingStyle{
		Name:    "upper snake case",
		Pattern: regexp.MustCompile(`^_?[A-Z][A-Z0-9]*(_[A-Z0-9]+)*$`),
	}
	lowerSnakeCaseStyle = NamingStyle{
		Name:    "lower snake case",
		Pattern: regexp.MustCompile(`^_?[a-z][a-z0-9]*(_[a-z0-9]+)*$`),
	}

	// CheckNamesDefaults is a map of Thrift AST node type to a regular expression for
	// validating names of that type.
	CheckNamesDefaults = map[reflect.Type]NamingStyle{
		thriftlint.ServiceType:   upperCamelCaseStyle,
		thriftlint.EnumType:      upperCamelCaseStyle,
		thriftlint.StructType:    upperCamelCaseStyle,
		thriftlint.EnumValueType: upperSnakeCaseStyle,
		thriftlint.FieldType:     lowerCamelCaseStyle,
		thriftlint.MethodType:    lowerCamelCaseStyle,
		thriftlint.ConstantType:  upperSnakeCaseStyle,
	}

	// CheckNamesDefaultBlacklist is names that should never be used for symbols.
	CheckNamesDefaultBlacklist = map[string]bool{
		"class": true,
		"int":   true,
	}
)

// CheckNames checks Thrift symbols comply with a set of regular expressions.
//
// If matches or blacklist are nil, global defaults will be used.
func CheckNames(
	matches map[reflect.Type]NamingStyle,
	blacklist map[string]bool,
) thriftlint.Check {
	if matches == nil {
		matches = CheckNamesDefaults
	}
	if blacklist == nil {
		blacklist = CheckNamesDefaultBlacklist
	}
	return &NamingChecker{
		matches:   matches,
		blacklist: blacklist,
	}
}

type NamingChecker struct {
	matches   map[reflect.Type]NamingStyle
	blacklist map[string]bool
}

func (c *NamingChecker) ID() string {
	return "naming"
}

func (c *NamingChecker) Checker() interface{} {
	return func(v interface{}) (messages thriftlint.Messages) {
		rv := reflect.Indirect(reflect.ValueOf(v))
		nameField := rv.FieldByName("Name")
		if !nameField.IsValid() {
			return nil
		}
		name := nameField.Interface().(string)
		// Special-case DEPRECATED_ fields.
		checker, ok := c.matches[rv.Type()]
		if !ok || strings.HasPrefix(name, "DEPRECATED_") {
			return nil
		}
		if c.blacklist[name] {
			messages.Warning(v, "%q is a disallowed name", name)
		}
		if ok := checker.Pattern.MatchString(name); !ok {
			messages.Warning(v, "name of %s %q should be %s", strings.ToLower(rv.Type().Name()),
				name, checker.Name)
		}
		return
	}
}

type NamingConf struct {
	Matches   map[string]string `yaml:"matches"`
	Blacklist []string          `yaml:"blacklist"`
}

func convert(src interface{}, dest interface{}) (err error) {
	b, err := yaml.Marshal(src)
	if err != nil {
		return
	}
	yaml.Unmarshal(b, dest)
	return
}

var reflectTypeMap = map[string]reflect.Type{
	"service":    thriftlint.ServiceType,
	"enum":       thriftlint.EnumType,
	"enum.value": thriftlint.EnumValueType,
	"field":      thriftlint.FieldType,
	"method":     thriftlint.MethodType,
	"constant":   thriftlint.ConstantType,
}

var namingStyleMap = map[string]NamingStyle{
	"upperCamel": upperCamelCaseStyle,
	"upperSnake": upperSnakeCaseStyle,
	"lowerCamel": lowerCamelCaseStyle,
	"lowerSnake": lowerSnakeCaseStyle,
}

func (c *NamingChecker) Reset() {
	c.matches = make(map[reflect.Type]NamingStyle)
	c.blacklist = make(map[string]bool)
}

func (c *NamingChecker) Config(config interface{}) (err error) {
	c.Reset()
	var conf NamingConf
	if err = convert(config, &conf); err != nil {
		return
	}
	for key, val := range conf.Matches {
		rtype, ok1 := reflectTypeMap[key]
		style, ok2 := namingStyleMap[val]
		if ok1 && ok2 {
			c.matches[rtype] = style
		}
	}
	for _, key := range conf.Blacklist {
		c.blacklist[key] = true
	}
	return
}
