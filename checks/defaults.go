package checks

import (
	"github.com/UrbanCompass/thriftlint"

	"github.com/alecthomas/go-thrift/parser"
)

// CheckDefaultValues checks that default values are not provided.
func CheckDefaultValues() thriftlint.Check {
	return thriftlint.MakeCheck("defaults", func(field *parser.Field) (messages thriftlint.Messages) {
		if field.Default != nil {
			messages.Warning(field, "default values are not allowed")
		}
		return
	})
}
