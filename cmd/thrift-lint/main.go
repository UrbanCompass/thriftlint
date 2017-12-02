package main

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/alecthomas/kingpin.v3-unstable"

	"github.com/UrbanCompass/thriftlint"
	"github.com/UrbanCompass/thriftlint/checks"
)

var (
	includeDirsFlag = kingpin.Flag("include", "Include directories to search.").Short('I').PlaceHolder("DIR").ExistingDirs()
	debugFlag       = kingpin.Flag("debug", "Enable debug logging.").Bool()
	disableFlag     = kingpin.Flag("disable", "Linters to disable.").PlaceHolder("LINTER").Strings()
	listFlag        = kingpin.Flag("list", "List linter checks.").Bool()
	errorFlag       = kingpin.Flag("errors", "Only show errors.").Bool()
	standard        = kingpin.Flag("standard", "The name or path of the coding standard to use").Default("default").String()

	sourcesArgs = kingpin.Arg("sources", "Thrift sources to lint.").Required().ExistingFiles()
)

func main() {
	kingpin.CommandLine.Help = `A linter for Thrift.

For details, please refer to https://github.com/UrbanCompass/thriftlint
`

	kingpin.Parse()

	fmt.Printf("using standard: %s\n", *standard)
	checkers, err := checks.AllCheckers.ApplyStandard(*standard)
	kingpin.FatalIfError(err, "")

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
