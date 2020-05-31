package main

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestReadDir(t *testing.T) {
	t.Run("default tests", func(t *testing.T) {
		path := "./testdata/env"
		ioutil.WriteFile(path+"/UNSET1", []byte("ERROR"), 0644)
		defer os.Remove(path + "/UNSET1")

		info, err := ReadDir(path)
		if err != nil {
			t.Fatal(err)
		}
		if len(info) != 5 {
			t.Fatal("error count")
		}
		if v, ok := info["BAR"]; !ok || v != "bar" {
			t.Fatal("error get value BAR")
		}
		if v, ok := info["FOO"]; !ok || v != "   foo\nwith new line" {
			t.Fatal("error get value FOO")
		}
		if v, ok := info["UNSET"]; !ok || len(v) > 0 {
			t.Fatal("error get value UNSET")
		}
		if v, ok := info["HELLO"]; !ok || v != `"hello"` {
			t.Fatal("error get value UNSET")
		}
		if v, ok := info["UNSET1"]; !ok || v != `ERROR` {
			t.Fatal("error get value UNSET1")
		}
	})
}
