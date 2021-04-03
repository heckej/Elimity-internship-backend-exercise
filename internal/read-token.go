package internal

import (
	"fmt"
	"io/ioutil"
)

// ReadTokenFromFile reads and returns the token from the file at the given filePath.
//
// The given filePath must refer to an existing file.
func ReadTokenFromFile(filePath string) (string, error) {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println("File reading error", err)
		return "", err
	}
	fmt.Println("Contents of file:", string(data))
	return string(data), nil
}
