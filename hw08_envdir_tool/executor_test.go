package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"testing"
)

const (
	dirCode  = "./testCode"
	pathCode = "./testCode/testCode.go"
	pathCMD  = "./testCode/testEnv"
	pathEnv  = "./testdata/env"
)

var ignoredEnv map[string]string

func init() {
	ignoredEnv = make(map[string]string)
	for _, env := range os.Environ() {
		nameValue := strings.SplitN(env, "=", 2)
		ignoredEnv[nameValue[0]] = nameValue[1]
	}
}

func TestRunInvalidParams(t *testing.T) {
	clearEnv()
	code := RunCmd(nil, nil)
	if code != 0 {
		t.Fatalf("error when cmd is nil")
	}
}

func TestRunCmdWithoutEnv(t *testing.T) {
	clearEnv()

	if !createTestBinary() {
		t.Fatalf("create test code")
	}
	defer os.RemoveAll(dirCode)

	code := RunCmd([]string{pathCMD}, nil)
	if code != 100 {
		t.Fatalf("error when Environment is nil")
	}
}

func TestRunCmd(t *testing.T) {
	if !createTestBinary() {
		t.Fatalf("create test code")
	}
	defer os.RemoveAll(dirCode)

	e, err := ReadDir(pathEnv)
	if err != nil {
		t.Fatalf("error ReadDir %v", err)
	}
	code := RunCmd([]string{pathCMD}, e)
	if code != 0 {
		t.Fatalf("error call RunCmd")
	}
}

func createTestBinary() bool {
	code := []byte(`
	package main

	import (
		"os"
		"fmt"
	)
	
	func main() {
		bar, ok := os.LookupEnv("BAR")
		if !ok || bar != "bar" {
			fmt.Println("ERROR")
			os.Exit(100)
		}
		fmt.Println("OK")
	}
	`)
	err := os.MkdirAll(dirCode, 0777)
	if err != nil {
		return false
	}
	err = ioutil.WriteFile(pathCode, code, 0777)
	if err != nil {
		fmt.Println(err)
		return false
	}
	cmd := exec.Command("go", "build", "-o", "testEnv")
	cmd.Dir = dirCode
	err = cmd.Run()
	if err != nil {
		return false
	}
	return true
}

func clearEnv() {
	for _, env := range os.Environ() {
		nameValue := strings.SplitN(env, "=", 2)
		if _, ok := ignoredEnv[nameValue[0]]; ok {
			continue
		}
		os.Unsetenv(nameValue[0])
	}
}
