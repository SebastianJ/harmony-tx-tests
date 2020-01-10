package utils

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

// IdentifyKeyFiles - identify pem files from a specified path
func IdentifyKeyFiles(path string) ([]string, error) {
	pattern := fmt.Sprintf("%s/_tx_gen_NEW_FUNDED_ACC_*/*", path)

	fmt.Println("Key file pattern is now: ", pattern)

	keys, err := globFiles(pattern)

	if err != nil {
		return nil, err
	}

	return keys, err
}

func globFiles(pattern string) ([]string, error) {
	files, err := filepath.Glob(pattern)

	if err != nil {
		return nil, err
	}

	return files, nil
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// ReadFileToString - check if a file exists, proceed to read it to memory if it does
func ReadFileToString(filePath string) (string, error) {
	if fileExists(filePath) {
		data, err := ioutil.ReadFile(filePath)

		if err != nil {
			return "", err
		}

		return string(data), nil
	} else {
		return "", nil
	}
}
