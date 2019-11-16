// +build run

package main

import (
	"fmt"
	"os"

	"github.com/zetamatta/go-termgap"
)

func main1() error {
	db, err := termgap.New()
	if err != nil {
		return err
	}
	w, err := db.RuneWidth('\u2727')
	if err != nil {
		return err
	}
	fmt.Printf("[\u2727]'s width=%d.\n", w)
	return nil
}

func main() {
	if err := main1(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
