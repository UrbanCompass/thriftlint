package thriftlint

import (
	"strings"

	"github.com/alecthomas/go-thrift/parser"
	"github.com/stretchr/testify/require"

	"testing"
)

func TestSplitSymbol(t *testing.T) {
	actual := SplitSymbol("someSnakeCaseAPI")
	require.Equal(t, []string{"some", "Snake", "Case", "API"}, actual)
	actual = SplitSymbol("someSnakeCase")
	require.Equal(t, []string{"some", "Snake", "Case"}, actual)
	actual = SplitSymbol("SomeCamelCase")
	require.Equal(t, []string{"Some", "Camel", "Case"}, actual)
	actual = SplitSymbol("SomeCamelCaseAPI")
	require.Equal(t, []string{"Some", "Camel", "Case", "API"}, actual)
	actual = SplitSymbol("some_underscore_case")
	require.Equal(t, []string{"some", "underscore", "case"}, actual)
	actual = SplitSymbol("SOME_UNDERSCORE_CASE")
	require.Equal(t, []string{"SOME", "UNDERSCORE", "CASE"}, actual)
	actual = SplitSymbol("ListingDBService")
	require.Equal(t, []string{"Listing", "DB", "Service"}, actual)
	actual = SplitSymbol("some1")
	require.Equal(t, []string{"some1"}, actual)
	actual = SplitSymbol("_id")
	require.Equal(t, []string{"", "id"}, actual)
	actual = SplitSymbol("listingIdSHAs")
	require.Equal(t, []string{"listing", "Id", "SHA", "s"}, actual)
	actual = SplitSymbol("listingIdSHAsToAdd")
	require.Equal(t, []string{"listing", "Id", "SHA", "s", "To", "Add"}, actual)
	actual = SplitSymbol("APIv3ProtocolTestService")
	require.Equal(t, []string{"API", "v3", "Protocol", "Test", "Service"}, actual)
}

func TestUpperCamelCase(t *testing.T) {
	actual := UpperCamelCase("listingIdSHAs")
	require.Equal(t, "ListingIDSHAs", actual)
	actual = UpperCamelCase("listingIdSHAsToAdd")
	require.Equal(t, "ListingIDSHAsToAdd", actual)
}

func TestLowerCamelCase(t *testing.T) {
	actual := LowerCamelCase("listingIdSHAs")
	require.Equal(t, "listingIDSHAs", actual)
	actual = LowerCamelCase("listingIdSHAsToAdd")
	require.Equal(t, "listingIDSHAsToAdd", actual)
}

func TestComment(t *testing.T) {
	enum := &parser.Enum{
		Comment: strings.Repeat("hello ", 30),
	}
	actual := Comment(enum)
	expected := []string{
		"hello hello hello hello hello hello hello hello hello hello hello hello hello",
		"hello hello hello hello hello hello hello hello hello hello hello hello hello",
		"hello hello hello hello",
	}
	require.Equal(t, expected, actual)
	enum = &parser.Enum{}
	actual = Comment(enum)
	require.Nil(t, actual)
}
