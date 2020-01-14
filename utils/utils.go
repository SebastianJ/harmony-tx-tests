package utils

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

// ParseYaml - parses yaml into a specific type
func ParseYaml(path string, entity interface{}) error {
	yamlData, err := ReadFileToString(path)

	if err != nil {
		return err
	}

	err = yaml.Unmarshal([]byte(yamlData), entity)

	if err != nil {
		return err
	}

	return nil
}

// RandomItemFromMap - select a random item from a map
func RandomItemFromMap(itemMap map[string]string) (string, string) {
	var keys []string

	for key := range itemMap {
		keys = append(keys, key)
	}

	randKey := RandomItemFromSlice(keys)
	randItem := itemMap[randKey]

	return randKey, randItem
}

// RandomItemFromSlice - select a random item from a slice
func RandomItemFromSlice(items []string) string {
	rand.Seed(time.Now().Unix())
	item := items[rand.Intn(len(items))]

	return item
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

// GlobFiles - find a set of files matching a specific pattern
func GlobFiles(pattern string) ([]string, error) {
	files, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}

	return files, nil
}
