package hybrid

import (
	"testing"
)

func TestRuneWidth(t *testing.T) {
	w := RuneWidth('\u2727')
	println("hybrid.RuneWidth=", w)
}
