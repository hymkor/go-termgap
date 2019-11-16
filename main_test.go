package termgap

import (
	"testing"
)

func TestRuneWidth(t *testing.T) {
	gap, err := New()
	if err != nil {
		t.Fatalf("termgap.json not found.(%s)", err.Error())
		return
	}
	w, err := gap.RuneWidth('A')
	if err != nil {
		t.Fatalf("RuneWidth('A') -> %s", err.Error())
		return
	}
	if w != 1 {
		t.Fatalf("RuneWidth('A') -> %d", w)
	}

	w, err = gap.RuneWidth('\u82a0') // hiragana's a
	if err != nil {
		t.Fatalf("RuneWidth('\u82a0') -> %s", err.Error())
		return
	}
	if w != 2 {
		t.Fatalf("RuneWidth('\u82a0') -> %d", w)
	}
}
