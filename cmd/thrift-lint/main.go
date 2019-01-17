package main

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/alecthomas/kingpin.v3-unstable"

	"github.com/wy90021/thriftlint"
	"github.com/wy90021/thriftlint/checks"
)

var (
	includeDirsFlag = kingpin.Flag("include", "Include directories to search.").Short('I').PlaceHolder("DIR").ExistingDirs()
	debugFlag       = kingpin.Flag("debug", "Enable debug logging.").Bool()
	disableFlag     = kingpin.Flag("disable", "Linters to disable.").PlaceHolder("LINTER").Strings()
	listFlag        = kingpin.Flag("list", "List linter checks.").Bool()
	errorFlag       = kingpin.Flag("errors", "Only show errors.").Bool()
	sourcesArgs     = kingpin.Arg("sources", "Thrift sources to lint.").Required().ExistingFiles()
)

func main() {
	kingpin.CommandLine.Help = `A linter for Thrift.

For details, please refer to https://github.com/wy90021/thriftlint
`
	kingpin.Parse()
	checkers := thriftlint.Checks{
		checks.CheckIndentation(),
		checks.CheckNames(nil, nil),
		checks.CheckOptional(),
		checks.CheckDefaultValues(),
		checks.CheckEnumSequence(),
		checks.CheckMapKeys(),
		checks.CheckTypeReferences(),
		checks.CheckStructFieldOrder(),
	}
	checkers = append(checkers, checks.CheckAnnotations(nil, checkers))

	if *listFlag {
		for _, linter := range checkers {
			fmt.Printf("%s\n", linter.ID())
		}
		return
	}

	options := []thriftlint.Option{
		thriftlint.WithIncludeDirs(*includeDirsFlag...),
		thriftlint.Disable(*disableFlag...),
	}
	if *debugFlag {
		logger := log.New(os.Stdout, "debug: ", 0)
		options = append(options, thriftlint.WithLogger(logger))
	}
	linter, err := thriftlint.New(checkers, options...)
	kingpin.FatalIfError(err, "")
	messages, err := linter.Lint(*sourcesArgs)
	kingpin.FatalIfError(err, "")
	status := 0
	for _, msg := range messages {
		if *errorFlag && msg.Severity != thriftlint.Error {
			continue
		}
		pos := thriftlint.Pos(msg.Object)
		fmt.Fprintf(os.Stderr, "%s:%d:%d:%s: %s (%s)\n", msg.File.Filename, pos.Line, pos.Col,
			msg.Severity, msg.Message, msg.Checker)
		status |= 1 << uint(msg.Severity)
	}
	os.Exit(status)
}
