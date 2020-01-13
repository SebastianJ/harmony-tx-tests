package accounts

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"regexp"

	"github.com/SebastianJ/harmony-tx-tests/balances"
	"github.com/SebastianJ/harmony-tx-tests/config"
	"github.com/SebastianJ/harmony-tx-tests/utils"
	"github.com/harmony-one/go-sdk/pkg/account"
	"github.com/harmony-one/go-sdk/pkg/address"
	"github.com/harmony-one/go-sdk/pkg/common"
	"github.com/harmony-one/go-sdk/pkg/store"
	homedir "github.com/mitchellh/go-homedir"
)

var (
	addressRegex = regexp.MustCompile(`"address":"(?P<address>[a-z0-9]+)"`)
)

// Account - represents a simple keystore account
type Account struct {
	Name    string
	Address string
}

// LoadSourceAccounts - loads the source accounts and imports them to the keystore when necessary
func LoadSourceAccounts() (accs []string, err error) {
	accountMapping, err := IdentifyKeys()

	if err != nil {
		return nil, err
	}

	fmt.Println(fmt.Sprintf("Found a total of %d keys", len(accs)))

	for path, address := range accountMapping {
		fmt.Println(fmt.Sprintf("Keyfile path: %s, address: %s", path, address))

		err := ImportAccount(path, address)

		if err != nil {
			return nil, err
		}
	}

	hasFunds, missingFunds, err := balances.FilterMinimumBalanceAccounts(accountMapping, config.Configuration.Funding.MinimumFunds)

	if err != nil {
		return nil, err
	}

	for path, address := range missingFunds {
		fmt.Println(fmt.Sprintf("Keyfile path: %s, address: %s doesn't hold any funds - removing the account and the key file...", path, address))
		RemoveAccount(address)
		os.RemoveAll(path)
	}

	for _, address := range hasFunds {
		accs = append(accs, address)
	}

	return accs, nil
}

// DoesNamedAccountExist - wrapper around store.DoesNamedAccountExist(name)
func DoesNamedAccountExist(name string) bool {
	return store.DoesNamedAccountExist(name)
}

// GenerateAccount - generates a new account using the specified name and passphrase
func GenerateAccount(name string) error {
	accountExists := store.DoesNamedAccountExist(name)

	if !accountExists {
		acc := account.Creation{
			Name:            name,
			Passphrase:      config.Configuration.Account.Passphrase,
			Mnemonic:        "",
			HdAccountNumber: nil,
			HdIndexNumber:   nil,
		}

		err := account.CreateNewLocalAccount(&acc)

		return err
	}

	return nil
}

// ImportAccount - imports an existing keystore
func ImportAccount(keyFile string, keyName string) error {
	accountExists := store.DoesNamedAccountExist(keyName)

	if accountExists {
		fmt.Println(fmt.Sprintf("Keyfile name: %s already exists in your keystore - proceeding...", keyName))
	} else {
		fmt.Println(fmt.Sprintf("Proceeding to import keyfile with path %s and name %s", keyFile, keyName))

		importedKeyName, err := account.ImportKeyStore(keyFile, keyName, config.Configuration.Account.Passphrase)

		if importedKeyName == "" || err != nil {
			return err
		}

		fmt.Println(fmt.Sprintf("Successfully imported keyfile with path: %s and name %s to the keystore!", keyFile, keyName))
	}

	return nil
}

// FindAccountAddressByName - finds the account address associated with a given key store name
func FindAccountAddressByName(targetName string) string {
	for _, name := range store.LocalAccounts() {
		if name == targetName {
			ks := store.FromAccountName(name)
			allAccounts := ks.Accounts()
			for _, account := range allAccounts {
				return address.ToBech32(account.Address)
			}
		}
	}

	return ""
}

// DoesAddressExistInKeystore - checks if a given address exists in the keystore
func DoesAddressExistInKeystore(targetAddress string) bool {
	exists := false

	for _, name := range store.LocalAccounts() {
		ks := store.FromAccountName(name)
		allAccounts := ks.Accounts()
		for _, account := range allAccounts {
			if targetAddress == address.ToBech32(account.Address) {
				return true
			}
		}
	}

	return exists
}

// GenerateAccountAndReturnAddress - Generate a new keystore account and return its address
func GenerateAccountAndReturnAddress(name string) (string, error) {
	err := GenerateAccount(name)

	if err != nil {
		return "", err
	}

	address := FindAccountAddressByName(name)

	return address, nil
}

// RemoveAccount - removes an account from the keystore
func RemoveAccount(name string) {
	uDir, _ := homedir.Dir()
	hmyCLIDir := path.Join(uDir, common.DefaultConfigDirName, common.DefaultConfigAccountAliasesDirName)
	accountDir := fmt.Sprintf("%s/%s", hmyCLIDir, name)
	os.RemoveAll(accountDir)
}

// IdentifyKeys - identifies key files
func IdentifyKeys() (map[string]string, error) {
	var path string = config.Configuration.Account.KeysPath
	var files []string
	keys := make(map[string]string)

	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})

	if err != nil {
		return nil, err
	}

	for _, file := range files {
		keyData, err := utils.ReadFileToString(file)

		if err == nil {
			keyDetails, err := parseKeyJson(keyData)

			if err == nil {
				if address, ok := keyDetails["address"]; ok {
					if address.(string) != "" {
						keys[file] = address.(string)
					}
				}
			}
		}
	}

	return keys, nil
}

func parseKeyJson(data string) (map[string]interface{}, error) {
	var rawData interface{}
	err := json.Unmarshal([]byte(data), &rawData)

	if err != nil {
		return nil, err
	}

	jsonData := rawData.(map[string]interface{})
	ethAddress := jsonData["address"].(string)

	addr := address.Parse(ethAddress)
	bech32 := address.ToBech32(addr)

	if bech32 != "" {
		jsonData["address"] = bech32
	}

	return jsonData, nil
}
