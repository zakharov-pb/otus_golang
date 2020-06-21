package main

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	t.Run("default tests", func(t *testing.T) {
		path := "./testdata/env"
		ioutil.WriteFile(path+"/UNSET1", []byte("ERROR"), 0644)
		defer os.Remove(path + "/UNSET1")

		info, err := ReadDir(path)
		if err != nil {
			t.Fatalf("error ReadDir: %v", err)
		}
		expectedEnv := Environment{
			"BAR":    "bar",
			"FOO":    "   foo\nwith new line",
			"UNSET":  "",
			"HELLO":  `"hello"`,
			"UNSET1": "ERROR",
		}
		require.Equal(t, expectedEnv, info)
	})
}
