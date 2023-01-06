package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"golang.org/x/sys/windows"

	"github.com/hymkor/go-termgap"
)

//go:generate go run generate.go

type WidthChecker struct {
	handle windows.Handle
}

func NewWidthChecker() (*WidthChecker, error) {
	handle, err := windows.GetStdHandle(windows.STD_ERROR_HANDLE)
	if err != nil {
		return nil, err
	}
	return &WidthChecker{handle: handle}, nil
}

func (wc *WidthChecker) X() (int, error) {
	var buffer windows.ConsoleScreenBufferInfo

	err := windows.GetConsoleScreenBufferInfo(wc.handle, &buffer)
	if err != nil {
		return 0, err
	}
	return int(buffer.CursorPosition.X), nil
}

func (wc *WidthChecker) Test(c rune) (int, error) {
	fmt.Fprintf(os.Stderr, "\r%c", c)

	width, err := wc.X()
	if err != nil {
		return 0, err
	}
	os.Stderr.Write([]byte{'\r'})
	return width, nil
}

func main1() error {
	wc, err := NewWidthChecker()
	if err != nil {
		return err
	}

	jsonPath, err := termgap.DatabasePath()
	if err != nil {
		return err
	}

	data := map[rune]int{}
	for _, rng := range table {
		for c := rng[0]; c <= rng[1]; c++ {
			w, err := wc.Test(c)
			if err != nil {
				return err
			}
			data[c] = w
		}
	}

	jsonData, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(jsonPath, jsonData, 0666)
}

func main() {
	if err := main1(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
