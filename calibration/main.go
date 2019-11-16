package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"unicode/utf16"

	"golang.org/x/sys/windows"

	"github.com/zetamatta/go-termgap"
)

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

func parseEastAsianWidthTxt(in io.Reader, f func(start, end rune, typ string) error) error {
	sc := bufio.NewScanner(in)
	for sc.Scan() {
		text := sc.Text()
		if len(text) <= 0 || text[0] == '#' {
			continue
		}
		field := strings.Fields(text)
		pair := strings.Split(field[0], ";")
		if len(pair) < 2 {
			continue
		}
		ranges := strings.Split(pair[0], "..")
		if len(ranges) >= 2 {
			start, err := strconv.ParseInt(ranges[0], 16, 64)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s:%s\n", err.Error(), text)
				continue
			}
			if utf16.IsSurrogate(rune(start)) || start >= 0xFFFF {
				continue
			}
			end, err := strconv.ParseInt(ranges[1], 16, 64)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s:%s\n", err.Error(), text)
				continue
			}
			if end <= 0x1F {
				continue
			}
			if err := f(rune(start), rune(end), pair[1]); err != nil {
				return err
			}
		} else {
			mid, err := strconv.ParseInt(ranges[0], 16, 64)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s:%s\n", err.Error(), text)
				continue
			}
			if utf16.IsSurrogate(rune(mid)) || mid >= 0xFFFF || mid <= 0x1F {
				continue
			}
			if err := f(rune(mid), rune(mid), pair[1]); err != nil {
				return err
			}
		}
	}
	return nil
}

func main1() error {
	wc, err := NewWidthChecker()
	if err != nil {
		return err
	}
	data := map[rune]int{}

	resp, err := http.Get("https://unicode.org/Public/12.1.0/ucd/EastAsianWidth.txt")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	err = parseEastAsianWidthTxt(resp.Body, func(start, end rune, typ string) error {
		for c := start; c <= end; c++ {
			width, err := wc.Test(c)
			if err != nil || width <= 0 {
				return err
			}
			data[rune(c)] = width
		}
		return nil
	})
	if err != nil {
		return err
	}

	jsonPath, err := termgap.DatabasePath()
	if err != nil {
		return err
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
