package thriftlint

import (
	"testing"

	"github.com/alecthomas/go-thrift/parser"
	"github.com/stretchr/testify/require"
)

func TestCallCheckerValidation(t *testing.T) {
	failf := func(*parser.Thrift) {}
	require.Panics(t, func() { callChecker(failf, []interface{}{}) })
	okf := func(*parser.Thrift) Messages { return nil }
	require.NotPanics(t, func() { callChecker(okf, []interface{}{&parser.Thrift{}}) })
}

func TestCallChecker(t *testing.T) {
	okfuncs := []interface{}{
		func(*parser.Thrift, *parser.Struct, *parser.Field) Messages {
			return Messages{}
		},
		func(*parser.Struct, *parser.Field) Messages { return Messages{} },
		func(*parser.Thrift, *parser.Field) Messages { return Messages{} },
		func(*parser.Field) Messages { return Messages{} },
		func(self interface{}) Messages { return Messages{} },
		func(parent, self interface{}) Messages { return Messages{} },
	}
	ancestors := []interface{}{&parser.Thrift{}, &parser.Struct{}, &parser.Field{}}
	for _, okf := range okfuncs {
		out := callChecker(okf, ancestors)
		require.NotNil(t, out)
	}

	badfuncs := []interface{}{
		func(*parser.Thrift) Messages { return Messages{} },
		func(*parser.Struct) Messages { return Messages{} },
		func(*parser.Field, *parser.Struct) Messages { return Messages{} },
	}
	for _, badf := range badfuncs {
		out := callChecker(badf, ancestors)
		require.Nil(t, out)
	}
}
