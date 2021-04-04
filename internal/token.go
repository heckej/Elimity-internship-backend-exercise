package internal

import (
	"io/ioutil"
)

// ReadTokenFromFile reads and returns the token from the file at the given filePath.
//
// The given filePath must refer to an existing file.
func ReadTokenFromFile(filePath string) (string, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
