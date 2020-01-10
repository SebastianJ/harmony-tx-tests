package utils

import (
	"fmt"
	"math/rand"
	"path/filepath"
	"regexp"
	"time"

	"github.com/harmony-one/go-sdk/pkg/address"
)

var (
	keyRegex = regexp.MustCompile(`\/_tx_gen_NEW_FUNDED_ACC_(?P<id>\d+)\/UTC--\d{4}-\d{2}-\d{2}T\d{2}-\d{2}-\d{2}\.\d+Z--(?P<address>[a-z0-9]+)`)
)

// RandomItemFromMap - select a random item from a map
func RandomItemFromMap(itemMap map[string]string) (string, string) {
	var keys []string

	for key, _ := range itemMap {
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

// ParseKeyDetailsFromKeyFile - parse key details from a given key path
func ParseKeyDetailsFromKeyFile(path string) (keyDetails map[string]string) {
	match := keyRegex.FindStringSubmatch(path)
	keyDetails = make(map[string]string)

	for i, name := range keyRegex.SubexpNames() {
		if i > 0 && i <= len(match) {
			keyDetails[name] = match[i]
		}
	}

	addr := address.Parse(keyDetails["address"])
	bech32 := address.ToBech32(addr)

	if bech32 != "" {
		keyDetails["address"] = bech32
	}

	return keyDetails
}

// IdentifyKeyFiles - identify pem files from a specified path
func IdentifyKeyFiles(path string) ([]string, error) {
	pattern := fmt.Sprintf("%s/_tx_gen_NEW_FUNDED_ACC_*/UTC*", path)

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
