package main

import (
	"errors"
	"os"
	"os/exec"
	"testing"
)

func compareFiles(f1, f2 string) error {
	const sttyResultCount = 2
	cmd := exec.Command("cmp", f1, f2)
	cmd.Stdin = os.Stdin
	_, err := cmd.Output()
	if err != nil {
		return err
	}
	return nil
}

func testCopyFile(from, to, reference string, offset, limit int64) error {
	err := Copy(from, to, offset, limit)
	if err != nil {
		return err
	}
	return compareFiles(to, reference)
}

func TestCopy(t *testing.T) {
	t.Run("copy files", func(t *testing.T) {
		err := testCopyFile("./testdata/input.txt", "./out.txt", "./testdata/out_offset0_limit0.txt", 0, 0)
		if err != nil {
			t.Fatal(err)
		}
		err = testCopyFile("./testdata/input.txt", "./out.txt", "./testdata/out_offset0_limit10.txt", 0, 10)
		if err != nil {
			t.Fatal(err)
		}
		err = testCopyFile("./testdata/input.txt", "./out.txt", "./testdata/out_offset0_limit1000.txt", 0, 1000)
		if err != nil {
			t.Fatal(err)
		}
		err = testCopyFile("./testdata/input.txt", "./out.txt", "./testdata/out_offset0_limit10000.txt", 0, 10000)
		if err != nil {
			t.Fatal(err)
		}
		err = testCopyFile("./testdata/input.txt", "./out.txt", "./testdata/out_offset100_limit1000.txt", 100, 1000)
		if err != nil {
			t.Fatal(err)
		}
		err = testCopyFile("./testdata/input.txt", "./out.txt", "./testdata/out_offset6000_limit1000.txt", 6000, 1000)
		if err != nil {
			t.Fatal(err)
		}
		err = testCopyFile("./testdata/input.txt", "./out.txt", "./testdata/out_offset0_limit1000.txt", 10000, 1000)
		if err == nil {
			t.Fatal("Error offset")
		}
	})
	t.Run("copy dir", func(t *testing.T) {
		err := Copy("./testdata", "./out.txt", 0, 0)
		if err == nil {
			t.Fatal("Error copy folder")
		}
		err = errors.Unwrap(err)
		if err != ErrUnsupportedFile {
			t.Fatal("Get error ErrUnsupportedFile")
		}
	})
}
