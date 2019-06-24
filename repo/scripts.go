package repo

import (
	"io/ioutil"
)

// GetSQLScript returns an SQL script obtained by the given path.
func GetSQLScript(path string) (string, error) {
	rawScript, err := ioutil.ReadFile(path)

	if err != nil {
		return "", err
	}

	script := string(rawScript)

	return script, nil
}
