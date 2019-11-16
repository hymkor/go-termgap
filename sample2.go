// +build run

package main

import (
	"fmt"
	"github.com/zetamatta/go-termgap/hybrid"
)

func main() {
	fmt.Printf("[A]'s width=%d\n", hybrid.RuneWidth('A'))
	fmt.Printf("[\u2727]'s width=%d\n", hybrid.RuneWidth('\u2727'))
}
