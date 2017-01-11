package thriftlint

import (
	"strings"

	"github.com/alecthomas/go-thrift/parser"
)

// Resolve a symbol within a file to its type.
func Resolve(symbol string, file *parser.Thrift) interface{} {
	parts := strings.SplitN(symbol, ".", 2)
	name := symbol
	target := file
	if len(parts) == 2 {
		target = file.Imports[parts[0]]
		if target == nil {
			return nil
		}
		name = parts[1]
		if target == nil {
			return nil
		}
	}
	if t, ok := target.Constants[name]; ok {
		return t
	}
	if t, ok := target.Enums[name]; ok {
		return t
	}
	if t, ok := target.Exceptions[name]; ok {
		return t
	}
	if t, ok := target.Services[name]; ok {
		return t
	}
	if t, ok := target.Structs[name]; ok {
		return t
	}
	if t, ok := target.Typedefs[name]; ok {
		return t
	}
	if t, ok := target.Unions[name]; ok {
		return t
	}
	return nil
}
