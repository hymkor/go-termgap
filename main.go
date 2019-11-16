package termgap

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Database struct {
	database map[rune]int
}

func DatabasePath() (string, error) {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return "", err
	}
	cacheDir = filepath.Join(cacheDir, "nyaos_org")
	if err = os.MkdirAll(cacheDir, 0777); err != nil {
		return "", err
	}
	return filepath.Join(cacheDir, "termgap.json"), nil
}

func New() (*Database, error) {
	jsonPath, err := DatabasePath()
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadFile(jsonPath)
	if err != nil {
		return nil, err
	}
	database := map[rune]int{}
	err = json.Unmarshal(data, &database)
	if err != nil {
		return nil, err
	}
	return &Database{
		database: database,
	}, nil
}

func (d *Database) RuneWidth(n rune) (int, error) {
	w, ok := d.database[n]
	if !ok {
		return 0, errors.New("not found")
	}
	return w, nil
}
