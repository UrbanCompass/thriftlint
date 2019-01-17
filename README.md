# An extensible linter for Thrift [![](https://godoc.org/github.com/wy90021/thriftlint?status.svg)](http://godoc.org/github.com/wy90021/thriftlint)

This is an extensible linter for [Thrift](https://thrift.apache.org/). It
includes a set of common lint checks, but allows for both customisation of those
checks, and creation of new ones by implementing the
[Check](https://godoc.org/github.com/wy90021/thriftlint#Check) interface.

For an example of how to build your own linter utility, please refer to the
[thrift-lint source](https://github.com/wy90021/thriftlint/tree/master/cmd/thrift-lint).

## Example checker

Here is an example of a checker utilising the `MakeCheck()` convenience
function to ensure that fields are present in the file in the same order as
their field numbers:

```go
func CheckStructFieldOrder() thriftlint.Check {
  return thriftlint.MakeCheck("field.order", func(s *parser.Struct) (messages thriftlint.Messages) {
    fields := sortedFields(s.Fields)
    sort.Sort(fields)
    for i := 0; i < len(fields)-1; i++ {
      a := fields[i]
      b := fields[i+1]
      if a.Pos.Line > b.Pos.Line {
        messages.Warning(fields[i], "field %d and %d are out of order", a.ID, b.ID)
      }
    }
    return
  })
}
```

## thrift-lint tool

A binary is included that can be used to perform basic linting with the builtin checks:

```
$ go get github.com/wy90021/thriftlint/cmd/thrift-lint
$ thrift-lint --help
usage: thrift-lint [<flags>] <sources>...

A linter for Thrift.

For details, please refer to https://github.com/wy90021/thriftlint

Flags:
      --help                Show context-sensitive help (also try --help-long
                            and --help-man).
  -I, --include=DIR ...     Include directories to search.
      --debug               Enable debug logging.
      --disable=LINTER ...  Linters to disable.
      --list                List linter checks.
      --errors              Only show errors.

Args:
  <sources>  Thrift sources to lint.
```
