package thriftlint

import (
	"github.com/alecthomas/go-thrift/parser"
	"github.com/stretchr/testify/require"

	"testing"
)

func TestResolve(t *testing.T) {
	parsed, err := parser.Parse("test.thrift", []byte(`
struct Struct {};
const string CONST = "";
enum Enum { CASE = 1; };
exception Exception {};
service Service {};
typedef Service Typedef;
union Union {};
`))
	require.NoError(t, err)
	ast := parsed.(*parser.Thrift)
	ast.Imports = map[string]*parser.Thrift{
		"pkg": ast,
	}

	actual := Resolve("Struct", ast)
	require.NotNil(t, actual)
	require.Equal(t, ast.Structs["Struct"], actual)
	actual = Resolve("pkg.Struct", ast)
	require.NotNil(t, actual)
	require.Equal(t, ast.Structs["Struct"], actual)

	actual = Resolve("CONST", ast)
	require.NotNil(t, actual)
	require.Equal(t, ast.Constants["CONST"], actual)
	actual = Resolve("pkg.CONST", ast)
	require.NotNil(t, actual)
	require.Equal(t, ast.Constants["CONST"], actual)

	actual = Resolve("Enum", ast)
	require.NotNil(t, actual)
	require.Equal(t, ast.Enums["Enum"], actual)
	actual = Resolve("pkg.Enum", ast)
	require.NotNil(t, actual)
	require.Equal(t, ast.Enums["Enum"], actual)

	actual = Resolve("Service", ast)
	require.NotNil(t, actual)
	require.Equal(t, ast.Services["Service"], actual)
	actual = Resolve("pkg.Service", ast)
	require.NotNil(t, actual)
	require.Equal(t, ast.Services["Service"], actual)

	actual = Resolve("Typedef", ast)
	require.NotNil(t, actual)
	require.Equal(t, ast.Typedefs["Typedef"], actual)
	actual = Resolve("pkg.Typedef", ast)
	require.NotNil(t, actual)
	require.Equal(t, ast.Typedefs["Typedef"], actual)

	actual = Resolve("Union", ast)
	require.NotNil(t, actual)
	require.Equal(t, ast.Unions["Union"], actual)
	actual = Resolve("pkg.Union", ast)
	require.NotNil(t, actual)
	require.Equal(t, ast.Unions["Union"], actual)
}
