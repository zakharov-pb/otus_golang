package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"testing"
)

const (
	dirCode  = "./testCode"
	pathCode = "./testCode/testCode.go"
	pathCMD  = "./testCode/testEnv"
	pathEnv  = "./testdata/env"
)

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

func clearEnv(e Environment) {
	for k := range e {
		os.Unsetenv(k)
	}
}

func TestRunCmd(t *testing.T) {
	if !createTestBinary() {
		t.Fatalf("create test code")
	}
	defer os.RemoveAll(dirCode)
	code := RunCmd(nil, nil)
	if code != 0 {
		t.Fatalf("error when cmd is nil")
	}
	e, err := ReadDir(pathEnv)
	if err != nil {
		t.Fatalf("error ReadDir %v", err)
	}
	code = RunCmd([]string{pathCMD}, e)
	if code != 0 {
		t.Fatalf("error call RunCmd")
	}
	clearEnv(e)
	code = RunCmd([]string{pathCMD}, nil)
	if code != 100 {
		t.Fatalf("error when Environment is nil")
	}
}
