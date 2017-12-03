package thriftlint

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/alecthomas/go-thrift/parser"
	yaml "gopkg.in/yaml.v2"
)

// Severity of a linter message.
type Severity int

// Message severities.
const (
	Warning Severity = iota
	Error
)

func (s Severity) String() string {
	if s == Warning {
		return "warning"
	}
	return "error"
}

// Message represents a single linter message.
type Message struct {
	// File that resulted in the message.
	File *parser.Thrift
	// ID of the Checker that generated this message.
	Checker  string
	Severity Severity
	Object   interface{}
	Message  string
}

// Messages is the set of messages each check should return.
//
// Typically it will be used like so:
//
// func MyCheck(...) (messages Messages) {
//   messages.Warning(t, "some warning")
// }
type Messages []*Message

// Warning adds a warning-level message to the Messages.
func (w *Messages) Warning(object interface{}, msg string, args ...interface{}) Messages {
	message := &Message{Severity: Warning, Object: object, Message: fmt.Sprintf(msg, args...)}
	*w = append(*w, message)
	return *w
}

// Warning adds an error-level message to the Messages.
func (w *Messages) Error(object interface{}, msg string, args ...interface{}) Messages {
	message := &Message{Severity: Error, Object: object, Message: fmt.Sprintf(msg, args...)}
	*w = append(*w, message)
	return *w
}

type Standard struct {
	Linters map[string]interface{} `yaml:"linters"`
}

// Checks is a convenience wrapper around a slice of Checks.
type Checks []Check

// CloneAndDisable returns a copy of this Checks slice with all checks matching prefix disabled.
func (c Checks) CloneAndDisable(prefixes ...string) Checks {
	out := Checks{}
skip:
	for _, check := range c {
		id := check.ID()
		for _, prefix := range prefixes {
			if prefix == id || strings.HasPrefix(id, prefix+".") {
				continue skip
			}
		}
		out = append(out, check)
	}
	return out
}

// Has returns true if the Checks slice contains any checks matching prefix.
func (c Checks) Has(prefix string) bool {
	for _, check := range c {
		id := check.ID()
		if prefix == id || strings.HasPrefix(id, prefix+".") {
			return true
		}
	}
	return false
}

func (c Checks) ApplyStandard(nameOrFilepath string) (checkers Checks, err error) {
	if nameOrFilepath == "default" {
		checkers = c
		return
	}

	f, err := os.Open(nameOrFilepath)
	if err != nil {
		return
	}
	buf := bytes.NewBuffer(nil)
	io.Copy(buf, f)

	standard := Standard{}
	if err = yaml.Unmarshal(buf.Bytes(), &standard); err != nil {
		return
	}

	checkers = make(Checks, 0)
	for _, checker := range c {
		if config, ok := standard.Linters[checker.ID()]; ok {
			checker.Config(config)
			checkers = append(checkers, checker)
		}
	}
	return
}

// Check implementations are used by the linter to check AST nodes.
type Check interface {
	// ID of the Check. Must be unique across all checks.
	//
	// IDs may be hierarchical, separated by a period. eg. "enum", "enum.values"
	ID() string

	// Checker returns the checking function.
	//
	// The checking function has the signature "func(...) Messages", where "..." is a sequence of
	// Thrift AST types that are matched against the current node's ancestors as the linter walks
	// the AST of each file.  "..." may also be "interface{}" in which case the checker function
	// will be called for each node in the AST.
	//
	// For example, the function:
	//
	//     func (s *parser.Struct, f *parser.Field) (messages Messages)
	//
	// Will match all each struct field, but not union fields.
	Checker() interface{}

	// Config receives a yaml unmarshalled interface{}
	Config(interface{}) error
}

// MakeCheck creates a stateless Check type from an ID and a checker function.
func MakeCheck(id string, checker interface{}) Check {
	return &statelessCheck{
		id:      id,
		checker: checker,
	}
}

type statelessCheck struct {
	id      string
	checker interface{}
}

func (s *statelessCheck) ID() string {
	return s.id
}

func (s *statelessCheck) Checker() interface{} {
	return s.checker
}

func (s *statelessCheck) Config(interface{}) error {
	return nil
}
