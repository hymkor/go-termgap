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
	"unicode/utf16"
)

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
		start, err := strconv.ParseInt(ranges[0], 16, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s:%s\n", err.Error(), text)
			continue
		}
		if utf16.IsSurrogate(rune(start)) || start >= 0xFFFF {
			continue
		}
		end := start
		if len(ranges) >= 2 {
			end, err = strconv.ParseInt(ranges[1], 16, 64)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s:%s\n", err.Error(), text)
				continue
			}
		}
		if end <= 0x1F {
			continue
		}
		if err := f(rune(start), rune(end), pair[1]); err != nil {
			return err
		}
	}
	return nil
}

func parseEmojiDataTxt(in io.Reader, f func(start, end rune) error) error {
	sc := bufio.NewScanner(in)
	for sc.Scan() {
		text := sc.Text()
		if len(text) <= 0 || text[0] == '#' {
			continue
		}
		field := strings.Fields(text)
		ranges := strings.Split(field[0], "..")
		start, err := strconv.ParseInt(ranges[0], 16, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s:%s\n", err.Error(), text)
			continue
		}
		if utf16.IsSurrogate(rune(start)) || start >= 0xFFFF {
			continue
		}
		end := start
		if len(ranges) >= 2 {
			end, err = strconv.ParseInt(ranges[1], 16, 64)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%s:%s\n", err.Error(), text)
				continue
			}
			if end <= 0x1F {
				continue
			}
		}
		if err := f(rune(start), rune(end)); err != nil {
			return err
		}
	}
	return nil
}

func main1() error {
	fd, err := os.Create("table.go")
	if err != nil {
		return err
	}
	defer fd.Close()

	fmt.Fprint(fd, "package main\n\nvar table = [][2]rune{\n")

	resp1, err := http.Get("https://unicode.org/Public/12.1.0/ucd/EastAsianWidth.txt")
	if err != nil {
		return err
	}
	defer resp1.Body.Close()

	lastStart := rune(-1)
	lastEnd := rune(-1)
	err = parseEastAsianWidthTxt(resp1.Body, func(start, end rune, typ string) error {
		if lastEnd+1 == start {
			lastEnd = end
			return nil
		}
		if lastStart > 0 {
			fmt.Fprintf(fd, "\t{%d, %d},\n", lastStart, lastEnd)
		}
		lastStart = start
		lastEnd = end
		return nil
	})
	if err != nil {
		return err
	}
	if lastStart > 0 {
		fmt.Fprintf(fd, "\t{%d, %d},\n", lastStart, lastEnd)
	}
	fmt.Fprintf(fd, "}\n")
	return nil
}

func main() {
	if err := main1(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
