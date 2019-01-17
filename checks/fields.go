package checks

import (
	"sort"

	"github.com/alecthomas/go-thrift/parser"

	"github.com/wy90021/thriftlint"
)

type sortedFields []*parser.Field

func (s sortedFields) Len() int           { return len(s) }
func (s sortedFields) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s sortedFields) Less(i, j int) bool { return s[i].ID < s[j].ID }

// CheckStructFieldOrder ensures that struct field IDs are present in-order in the file.
func CheckStructFieldOrder() thriftlint.Check {
	return thriftlint.MakeCheck("field.order", func(s *parser.Struct) (messages thriftlint.Messages) {
		fields := sortedFields(s.Fields)
		sort.Sort(fields)
		for i := 0; i < len(fields)-1; i++ {
			a := fields[i]
			b := fields[i+1]
			if a.Pos.Line > b.Pos.Line {
				messages.Warning(fields[i], "field %d and %d of %s are out of order", a.ID, b.ID, s.Name)
			}
		}
		return
	})
}
