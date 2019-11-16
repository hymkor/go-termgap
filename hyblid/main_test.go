package hyblid

import (
	"testing"
)

func TestRuneWidth(t *testing.T) {
	w := RuneWidth('\u2727')
	println("hyblid.RuneWidth=", w)
}
