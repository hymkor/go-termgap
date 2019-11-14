// +build run

package main

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"

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

func (wc *WidthChecker) X() (int16, error) {
	var buffer windows.ConsoleScreenBufferInfo

	err := windows.GetConsoleScreenBufferInfo(wc.handle, &buffer)
	if err != nil {
		return 0, err
	}
	return buffer.CursorPosition.X, nil
}

func (wc *WidthChecker) Test(c int64) (int16, error) {
	fmt.Fprintf(os.Stderr, "\r%c", rune(c))

	width, err := wc.X()
	if err != nil {
		return 0, err
	}
	os.Stderr.Write([]byte{'\r'})
	return width, nil
}

func test(wc *WidthChecker, c int64, typeStr string, out io.Writer) error {
	if typeStr == "N" {
		return nil
	}
	width, err := wc.Test(c)
	if err != nil {
		return err
	}
	switch typeStr {
	case "Na":
		if width == 1 {
			return nil
		}
	case "W", "F", "A":
		if width == 2 {
			return nil
		}
	}
	fmt.Fprintf(out, "\t'\\u%04X': %d,\n", c, width)
	return nil
}

func tryUrl(wc *WidthChecker, url string, out io.Writer) error {
	fmt.Fprintf(out, "\t// diff from %s\n", url)
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
		field1 := strings.Split(field[0], ";")
		if len(field1) < 2 {
			continue
		}
		ranges := strings.Split(field1[0], "..")
		if len(ranges) >= 2 {
			start, err := strconv.ParseInt(ranges[0], 16, 64)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s:%s\n", err.Error(), text)
				continue
			}
			if start >= 0xFFFF {
				continue
			}
			end, err := strconv.ParseInt(ranges[1], 16, 64)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s:%s\n", err.Error(), text)
				continue
			}
			for i := start; i <= end; i++ {
				err := test(wc, i, field1[1], out)
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
			err = test(wc, mid, field1[1], out)
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
	fd, err := os.Create("table.go")
	if err != nil {
		return err
	}
	defer fd.Close()
	fmt.Fprintln(fd, "package lunewhydos")
	fmt.Fprintln(fd, "")
	fmt.Fprintln(fd, "var table = map[rune]int{")

	err = tryUrl(wc, "https://unicode.org/Public/12.1.0/ucd/EastAsianWidth.txt", fd)
	if err != nil {
		return err
	}
	fmt.Fprintln(fd, "}")
	return nil
}

func main() {
	if err := main1(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
