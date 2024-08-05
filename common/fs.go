package common

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/AspieSoft/go-regex-re2/v2"
	"gopkg.in/yaml.v3"
)

// JoinPath joins multiple file types with safety from backtracking
func JoinPath(path ...string) (string, error) {
	resPath, err := filepath.Abs(string(path[0]))
	if err != nil {
		return "", err
	}
	for i := 1; i < len(path); i++ {
		p := filepath.Join(resPath, string(path[i]))
		if p == resPath || !strings.HasPrefix(p, resPath) {
			return "", errors.New("path leaked outside of root")
		}
		resPath = p
	}
	return resPath, nil
}

// Copy lets you copy files from the src to the dst
func CopyFile(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

// ReadYaml loads a yaml file into a struct
//
// this method will read the buffer, and normalize names so
// '-' and '_' characters are optional, and everything is lowercase
//
// this method is useful for loading a config file
func ReadYaml(path string, out interface{}) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	b = regex.Comp(`(?m)^(\s*(?:-\s+|))([\w_\-]+):`).RepFunc(b, func(data func(int) []byte) []byte {
		return regex.JoinBytes(data(1), bytes.ReplaceAll(bytes.ReplaceAll(bytes.ToLower(data(2)), []byte{'-'}, []byte{}), []byte{'_'}, []byte{}), ':')
	})
	return yaml.Unmarshal(b, out)
}

// ReadJson loads a json file into a struct
//
// this method will read the buffer, and normalize names so
// '-' and '_' characters are optional, and everything is lowercase
//
// this method is useful for loading a config file
func ReadJson(path string, out interface{}) error {
	b, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	b = regex.Comp(`(?s)"([\w_-]+)"\s*:`).RepFunc(b, func(data func(int) []byte) []byte {
		return regex.JoinBytes('"', bytes.ReplaceAll(bytes.ReplaceAll(bytes.ToLower(data(1)), []byte{'-'}, []byte{}), []byte{'_'}, []byte{}), '"', ':')
	})
	return json.Unmarshal(b, out)
}

// ReadConfig loads a config file into a struct
//
// this method will read the buffer, and normalize names so
// '-' and '_' characters are optional, and everything is lowercase
//
// this method will try different file types in the following order:
//
//	[yml, yaml, json]
//
// you can specify the first file type to try, by adding a .ext of that file type to the path
//
// by accepting moltiple file types, the user can choose what type of file they want to use for their config file
func ReadConfig(path string, out interface{}) error {
	t := "yaml"
	var b []byte
	var err error = io.EOF

	// path .ext prioritize
	if strings.HasSuffix(path, ".yml") {
		t = "yaml"
		path = strings.TrimSuffix(path, ".yml")

		b, err = os.ReadFile(path + ".yml")
		if err != nil {
			b, err = os.ReadFile(path + ".yaml")
		}
	} else if strings.HasSuffix(path, ".yaml") {
		t = "yaml"
		path = strings.TrimSuffix(path, ".yaml")

		b, err = os.ReadFile(path + ".yaml")
		if err != nil {
			b, err = os.ReadFile(path + ".yml")
		}
	} else if strings.HasSuffix(path, ".json") {
		t = "json"
		path = strings.TrimSuffix(path, ".json")

		b, err = os.ReadFile(path + ".json")
	}

	// try yml
	if err != nil {
		t = "yaml"
		b, err = os.ReadFile(path + ".yml")
	}

	// try yaml
	if err != nil {
		t = "yaml"
		b, err = os.ReadFile(path + ".yaml")
	}

	// try json
	if err != nil {
		t = "json"
		b, err = os.ReadFile(path + ".json")
	}

	if err != nil {
		return io.EOF
	}

	switch t {
	case "yaml":
		b = regex.Comp(`(?m)^(\s*(?:-\s+|))([\w_\-]+):`).RepFunc(b, func(data func(int) []byte) []byte {
			return regex.JoinBytes(data(1), bytes.ReplaceAll(bytes.ReplaceAll(bytes.ToLower(data(2)), []byte{'-'}, []byte{}), []byte{'_'}, []byte{}), ':')
		})
		return yaml.Unmarshal(b, out)
	case "json":
		b = regex.Comp(`(?s)"([\w_-]+)"\s*:`).RepFunc(b, func(data func(int) []byte) []byte {
			return regex.JoinBytes('"', bytes.ReplaceAll(bytes.ReplaceAll(bytes.ToLower(data(1)), []byte{'-'}, []byte{}), []byte{'_'}, []byte{}), '"', ':')
		})
		return json.Unmarshal(b, out)
	default:
		return io.EOF
	}
}
