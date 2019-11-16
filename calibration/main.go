// +build run

package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"unicode/utf16"

	"golang.org/x/sys/windows"
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

func test(wc *WidthChecker, c rune, typeStr string, out map[rune]int) error {
	width, err := wc.Test(c)
	if err != nil {
		return err
	}
	switch typeStr {
	case "Na", "N":
		if width == 1 {
			return nil
		}
	case "W", "F", "A":
		if width == 2 {
			return nil
		}
	}
	out[rune(c)] = width
	return nil
}

func tryUrl(wc *WidthChecker, url string, out map[rune]int) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	sc := bufio.NewScanner(resp.Body)
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
			for i := start; i <= end; i++ {
				err := test(wc, rune(i), pair[1], out)
				if err != nil {
					return err
				}
			}
		} else {
			mid, err := strconv.ParseInt(ranges[0], 16, 64)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s:%s\n", err.Error(), text)
				continue
			}
			if utf16.IsSurrogate(rune(mid)) || mid >= 0xFFFF {
				continue
			}
			err = test(wc, rune(mid), pair[1], out)
			if err != nil {
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
	err = tryUrl(wc, "https://unicode.org/Public/12.1.0/ucd/EastAsianWidth.txt", data)
	if err != nil {
		return err
	}
	jsonData, err := json.MarshalIndent(data, "", "\t")
	if err != nil {
		return err
	}
	return ioutil.WriteFile("termgap.json", jsonData, 0666)
}

func main() {
	if err := main1(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
