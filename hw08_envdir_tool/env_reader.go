package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

// Environment environment variables.
type Environment map[string]string

type nameValue struct {
	Name  string
	Value string
	Error error
}

var (
	// ErrCharacter forbidden character '=' in file name.
	ErrCharacter = fmt.Errorf("forbidden character '=' in file name")
)

func readValueFromFile(dir string, name string, wg *sync.WaitGroup, result chan nameValue) {
	defer wg.Done()
	if strings.Contains(name, "=") {
		result <- nameValue{
			Name:  name,
			Error: ErrCharacter,
		}
	}
	file, err := os.Open(dir + name)
	if err != nil {
		result <- nameValue{
			Name:  name,
			Error: err,
		}
		return
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
		result <- nameValue{
			Name:  name,
			Value: value,
			Error: err,
		}
	} else {
		result <- nameValue{
			Name:  name,
			Value: "",
			Error: err,
		}
	}
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	dir = filepath.ToSlash(dir)
	if dir[len(dir)-1] != '/' {
		dir += string(os.PathSeparator)
	}
	res := Environment{}
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("error read dir: %w", err)
	}
	readerWg := sync.WaitGroup{}
	valuesCh := make(chan nameValue)
	errorsList := make([]error, 0)
	readerWg.Add(1)
	go func() {
		defer readerWg.Done()
		for nv := range valuesCh {
			if nv.Error != nil {
				errorsList = append(errorsList, nv.Error)
				continue
			}
			if old, ok := res[nv.Name]; !ok || len(old) != 0 {
				res[nv.Name] = nv.Value
			}
		}
	}()

	writersWg := sync.WaitGroup{}
	for _, f := range files {
		if f.IsDir() || !f.Mode().IsRegular() {
			continue
		}
		writersWg.Add(1)
		go readValueFromFile(dir, f.Name(), &writersWg, valuesCh)
	}
	writersWg.Wait()

	close(valuesCh)
	readerWg.Wait()

	if len(errorsList) > 0 {
		err = fmt.Errorf("errors ReadDir: %v", errorsList) // nolint
	}
	return res, err
}
