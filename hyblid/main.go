package hyblid

import (
	"github.com/mattn/go-runewidth"
	"github.com/zetamatta/go-termgap"
)

var runeWidth func(rune) int

func RuneWidth(ch rune) int {
	if runeWidth == nil {
		db, err := termgap.New()
		if err == nil {
			runeWidth = func(ch rune) int {
				w, err := db.RuneWidth(ch)
				if err == nil {
					return w
				} else {
					return runewidth.RuneWidth(ch)
				}
			}
		} else {
			runeWidth = runewidth.RuneWidth
		}
	}
	return runeWidth(ch)
}
