package thriftlint

import (
	"reflect"

	"github.com/alecthomas/go-thrift/parser"
)

// Annotation returns the annotation value associated with "key" from the .Annotations field of a
// go-thrift AST node.
//
// This will panic if node is not a struct with an Annotations field of the correct type.
func Annotation(node interface{}, key, dflt string) string {
	annotations := reflect.Indirect(reflect.ValueOf(node)).
		FieldByName("Annotations").
		Interface().([]*parser.Annotation)
	for _, annotation := range annotations {
		if annotation.Name == key {
			return annotation.Value
		}
	}
	return dflt
}

// AnnotationExists checks if an annotation is present at all.
func AnnotationExists(node interface{}, key string) bool {
	return Annotation(node, key, "\000") != "\000"
}
