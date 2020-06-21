package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"sync"
)

const (
	reservedAmount = 9
	minSize        = 50
)

var (
	// ErrProgressBarIsRun call error in method Run.
	ErrProgressBarIsRun = errors.New("ProgressBar is already running")
	// ErrSTTY invalid stty answer.
	ErrSTTY = errors.New("error getting terminal width: Invalid stty answer")
)

// ProgressBar struct for displaying a process.
type ProgressBar struct {
	progressCh chan int64
	doneCh     chan struct{}
	wg         sync.WaitGroup
	line       []byte
}

func getMaxSize() (size int64, err error) {
	const sttyResultCount = 2
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	out, err := cmd.Output()
	if err != nil {
		return 0, fmt.Errorf("error getting terminal width: %w", err)
	}
	items := strings.Split(string(out), " ")
	if len(items) < sttyResultCount {
		return 0, ErrSTTY
	}
	s, err := strconv.Atoi(items[1][:len(items[1])-1])
	if err != nil {
		return 0, fmt.Errorf("error getting terminal width: %w", err)
	}
	size = int64(s)
	return
}

func (p *ProgressBar) fillLine(progress int64) {
	var i int64
	count := int64(len(p.line))
	for i = 0; i < count; i++ {
		switch {
		case i == progress:
			p.line[i] = '>'
		case i < progress:
			p.line[i] = '='
		default:
			p.line[i] = ' '
		}
	}
}

// Run start rendering, progress - channel for transferring the state of the process.
func (p *ProgressBar) Run(label string, max int64) (progress chan<- int64, err error) {
	if p.progressCh != nil {
		return nil, ErrProgressBarIsRun
	}
	width, err := getMaxSize()
	if err != nil || width < minSize {
		width = minSize
	}
	lbl := []rune(label)
	if int64(len(lbl)) > width {
		lbl = lbl[:width-reservedAmount]
		newSize := len(lbl)
		lbl[newSize-1] = '.'
		lbl[newSize-2] = '.'
		lbl[newSize-3] = '.'
		label = string(lbl)
	}
	countLine := width - reservedAmount - int64(len(lbl))
	if countLine < 0 {
		countLine = 0
	}
	p.line = make([]byte, countLine)
	p.progressCh = make(chan int64)
	p.doneCh = make(chan struct{})
	p.wg.Add(1)
	go func() {
		defer func() {
			fmt.Println()
			p.wg.Done()
		}()
		p.fillLine(0)
		fmt.Printf("\r[%s]   0%% %s", p.line, label)
		for {
			select {
			case <-p.doneCh:
				return
			case newVal := <-p.progressCh:
				if newVal < 0 {
					newVal = 0
				}
				if newVal >= max {
					p.fillLine(int64(len(p.line)))
					fmt.Printf("\r[%s] 100%% %s", p.line, label)
					return
				}
				p.fillLine((newVal * int64(len(p.line))) / max)
				fmt.Printf("\r[%s] %3d%% %s", p.line, (newVal*100)/max, label)
			}
		}
	}()
	return p.progressCh, nil
}

// Stop stop work.
func (p *ProgressBar) Stop() {
	if p.doneCh == nil {
		return
	}
	close(p.doneCh)
	p.wg.Wait()
	p.doneCh = nil
	p.progressCh = nil
}
