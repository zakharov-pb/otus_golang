package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

// Environment environment variables.
type Environment map[string]string

var (
	// ErrCharacter forbidden character '=' in file name.
	ErrCharacter = fmt.Errorf("forbidden character '=' in file name")
)

func readValueFromFile(dir string, name string) (string, string, error) {
	if strings.Contains(name, "=") {
		return "", "", ErrCharacter
	}
	file, err := os.Open(filepath.Join(dir, name))
	if err != nil {
		return "", "", fmt.Errorf("error open file %s: %w", name, err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	if scanner.Scan() {
		value := strings.TrimRightFunc(
			strings.ReplaceAll(scanner.Text(), "\x00", "\n"),
			func(r rune) bool {
				if r == ' ' || r == '\t' {
					return true
				}
				return false
			})
		return name, value, nil
	}
	return name, "", nil
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	res := Environment{}
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("error read dir: %w", err)
	}
	for _, f := range files {
		if f.IsDir() || !f.Mode().IsRegular() {
			continue
		}
		name, value, err := readValueFromFile(dir, f.Name())
		if err != nil {
			return nil, fmt.Errorf("error read file %s: %w", name, err)
		}
		if v, ok := res[name]; ok && len(v) == 0 {
			continue
		}
		res[name] = value
	}
	return res, err
}
