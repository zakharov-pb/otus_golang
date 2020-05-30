package main

import (
	"errors"
	"fmt"
	"io"
	"os"
)

const (
	minBufferSize = 100
	maxBufferSize = 1024
	numberPart    = 3
)

var (
	// ErrUnsupportedFile unsupported file
	ErrUnsupportedFile = errors.New("unsupported file")
	// ErrOffsetExceedsFileSize offset exceeds file size
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func getFileSize(f *os.File) (int64, error) {
	fi, err := f.Stat()
	if err != nil {
		return 0, fmt.Errorf("error Copy %w", err)
	}
	if fi.IsDir() {
		return 0, ErrUnsupportedFile
	}
	if !fi.Mode().IsRegular() {
		if limit <= 0 {
			return 0, ErrUnsupportedFile
		}
		return limit, nil
	}
	return fi.Size(), nil
}

// Copy function copy files
func Copy(fromPath string, toPath string, offset, limit int64) error {
	fmt.Printf("\nFROM: %s\nTO: %s\n", fromPath, toPath)
	out, err := os.Open(fromPath)
	if err != nil {
		return fmt.Errorf("error open file %s  %w", fromPath, err)
	}
	defer out.Close()
	in, err := os.Create(toPath)
	if err != nil {
		return fmt.Errorf("error create file %s %w", to, err)
	}
	defer in.Close()
	size, err := getFileSize(out)
	if err != nil {
		return fmt.Errorf("error get size %s %w", fromPath, err)
	}
	if size < offset {
		return ErrOffsetExceedsFileSize
	}
	if limit <= 0 || (limit+offset) > size {
		limit = size - offset
	}
	var bufferSize int
	if limit < minBufferSize {
		bufferSize = int(limit)
	} else {
		bufferSize = int(limit / numberPart)
		if bufferSize > maxBufferSize {
			bufferSize = maxBufferSize
		}
	}
	buffer := make([]byte, bufferSize)
	p := ProgressBar{}
	progress, err := p.Run("COPY", limit)
	if err != nil {
		return fmt.Errorf("error run progress bar %w", err)
	}
	defer p.Stop()
	_, err = out.Seek(offset, io.SeekStart)
	if err != nil {
		return fmt.Errorf("error set offset %w", err)
	}
	for cur := int64(0); cur < limit; {
		count, err := out.Read(buffer)
		if err == io.EOF {
			return fmt.Errorf("error read data %w", err)
		}
		cur += int64(count)
		if limit-cur < 0 {
			count -= int(cur - limit)
		}
		_, err = in.Write(buffer[:count])
		if err != nil {
			return fmt.Errorf("error write data %w", err)
		}
		progress <- cur
	}
	return nil
}
