package checks

import (
	"reflect"

	"github.com/wy90021/thriftlint"
)

var (
	// Map of (parent, self) AST node types to the expected indentation for that type in that
	// context.
	expectedIndentation = map[string]int{
		indentationContext(thriftlint.ThriftType, thriftlint.ServiceType):  1,
		indentationContext(thriftlint.ServiceType, thriftlint.MethodType):  3,
		indentationContext(thriftlint.ThriftType, thriftlint.EnumType):     1,
		indentationContext(thriftlint.EnumType, thriftlint.EnumValueType):  3,
		indentationContext(thriftlint.ThriftType, thriftlint.StructType):   8,
		indentationContext(thriftlint.StructType, thriftlint.FieldType):    3,
		indentationContext(thriftlint.ThriftType, thriftlint.TypedefType):  1,
		indentationContext(thriftlint.ThriftType, thriftlint.ConstantType): 1,
	}
)

func indentationContext(parent, self reflect.Type) string {
	return parent.Name() + ":" + self.Name()
}

// CheckIndentation checks indentation is a multiple of 2.
func CheckIndentation() thriftlint.Check {
	return thriftlint.MakeCheck("indentation", func(parent, self interface{}) (messages thriftlint.Messages) {
		context := indentationContext(reflect.TypeOf(parent), reflect.TypeOf(self))
		pos := thriftlint.Pos(self)
		if expected, ok := expectedIndentation[context]; ok && expected != pos.Col {
			messages.Warning(self, "should be indented to column %d not %d", expected, pos.Col)
		}
		return
	})
}
