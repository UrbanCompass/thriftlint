package thriftlint

import (
	"io"
	"os"
	"path/filepath"

	"github.com/alecthomas/go-thrift/parser"
)

// Parse a set of .thrift source files into their corresponding ASTs.
func Parse(includeDirs []string, sources []string) (map[string]*parser.Thrift, error) {
	p := parser.New()
	p.Filesystem = &includeFilesystem{IncludeDirs: includeDirs}

	var files map[string]*parser.Thrift
	for _, path := range sources {
		path, err := filepath.Abs(path)
		if err != nil {
			return nil, err
		}
		files, _, err = p.ParseFile(path)
		if err != nil {
			return nil, err
		}
	}
	for _, file := range files {
		file.Imports = map[string]*parser.Thrift{}
		for symbol, path := range file.Includes {
			file.Imports[symbol] = files[path]
		}
	}
	return files, nil
}

// A go-thrift/parser.Filesystem implementation that searches include dirs when attempting to open
// sources.
type includeFilesystem struct {
	IncludeDirs []string
}

func (i *includeFilesystem) Open(filename string) (io.ReadCloser, error) {
	if filepath.IsAbs(filename) {
		return os.Open(filename)
	}
	for _, d := range i.IncludeDirs {
		path := filepath.Join(d, filename)
		r, err := os.Open(path)
		if os.IsNotExist(err) {
			continue
		} else if err != nil {
			return nil, err
		}
		return r, nil
	}
	return nil, os.ErrNotExist
}

func (i *includeFilesystem) Abs(dir, path string) (string, error) {
	if filepath.IsAbs(path) {
		return path, nil
	}
	for _, d := range i.IncludeDirs {
		p, err := filepath.Abs(filepath.Join(d, path))
		if err != nil {
			continue
		}
		if _, err := os.Stat(p); err == nil {
			return filepath.Abs(p)
		}
	}
	return filepath.Abs(filepath.Join(dir, path))
}
