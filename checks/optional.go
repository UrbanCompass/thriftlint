package checks

import (
	"github.com/wy90021/thriftlint"

	"github.com/alecthomas/go-thrift/parser"
)

// CheckOptional ensures that all Thrift fields are optional, as is generally accepted best
// practice for Thrift.
func CheckOptional() thriftlint.Check {
	return thriftlint.MakeCheck("optional", func(s *parser.Struct, f *parser.Field) (messages thriftlint.Messages) {
		if f.Type.Name != "list" && f.Type.Name != "set" && f.Type.Name != "map" && !f.Optional {
			messages.Warning(f, "%s must be optional", f.Name)
		}
		return
	})
}
