package accounts

import (
	"fmt"
	"os"
	"path"

	"github.com/harmony-one/go-sdk/pkg/account"
	"github.com/harmony-one/go-sdk/pkg/address"
	"github.com/harmony-one/go-sdk/pkg/common"
	"github.com/harmony-one/go-sdk/pkg/store"
	homedir "github.com/mitchellh/go-homedir"
)

// GenerateAccount - generates a new account using the specified name and passphrase
func GenerateAccount(name string, passphrase string) error {
	accountExists := store.DoesNamedAccountExist(name)

	if !accountExists {
		acc := account.Creation{
			Name:            name,
			Passphrase:      passphrase,
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
func ImportAccount(keyFile string, keyName string, passphrase string, keyDetails map[string]string) error {
	accountExists := store.DoesNamedAccountExist(keyName)

	if accountExists {
		fmt.Println(fmt.Sprintf("Keyfile with id: %s, name: %s and address %s already exists in your keystore - proceeding...", keyDetails["id"], keyName, keyDetails["address"]))
	} else {
		fmt.Println(fmt.Sprintf("Proceeding to import keyfile with id: %s, name: %s and address %s", keyDetails["id"], keyName, keyDetails["address"]))

		importedKeyName, err := account.ImportKeyStore(keyFile, keyName, passphrase)

		if importedKeyName == "" || err != nil {
			return err
		}

		fmt.Println(fmt.Sprintf("Successfully imported keyfile with id: %s, name: %s and address %s to the keystore!", keyDetails["id"], keyName, keyDetails["address"]))
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

// RemoveAccount - removes an account from the keystore
func RemoveAccount(name string) {
	uDir, _ := homedir.Dir()
	hmyCLIDir := path.Join(uDir, common.DefaultConfigDirName, common.DefaultConfigAccountAliasesDirName)
	accountDir := fmt.Sprintf("%s/%s", hmyCLIDir, name)
	os.RemoveAll(accountDir)
}
