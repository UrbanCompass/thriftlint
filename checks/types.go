package checks

import (
	"github.com/alecthomas/go-thrift/parser"

	"github.com/wy90021/thriftlint"
)

// CheckTypeReferences checks that types referenced in Thrift files are actually imported
// and exist.
func CheckTypeReferences() thriftlint.Check {
	return thriftlint.MakeCheck("types", func(file *parser.Thrift, t *parser.Type) (messages thriftlint.Messages) {
		if !thriftlint.BuiltinThriftTypes[t.Name] && !thriftlint.BuiltinThriftCollections[t.Name] &&
			thriftlint.Resolve(t.Name, file) == nil {
			messages.Error(t, "unknown type %q", t.Name)
		}
		return
	})
}
