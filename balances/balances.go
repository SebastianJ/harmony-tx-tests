package balances

import (
	"fmt"

	sdkBalances "github.com/SebastianJ/harmony-sdk/balances"
	"github.com/SebastianJ/harmony-tx-tests/config"
)

// GetShardBalance - gets the balance for a given address and shard
func GetShardBalance(address string, shardID uint32) (float64, error) {
	return sdkBalances.GetShardBalance(address, shardID, config.Configuration.Network.Name)
}

// FilterMinimumBalanceAccounts - Filters out a list of accounts without any balance
func FilterMinimumBalanceAccounts(accounts map[string]string, minimumBalance float64) (map[string]string, map[string]string, error) {
	hasFunds := make(map[string]string)
	missingFunds := make(map[string]string)

	for keyName, address := range accounts {
		totalBalance, err := sdkBalances.GetTotalBalance(address, config.Configuration.Network.Name)

		if err != nil {
			return nil, nil, err
		}

		if totalBalance > minimumBalance {
			hasFunds[keyName] = address
		} else {
			missingFunds[keyName] = address
		}
	}

	return hasFunds, missingFunds, nil
}

// OutputBalanceStatusForAddresses - outputs balance status
func OutputBalanceStatusForAddresses(accounts map[string]string, minimumBalance float64) {
	hasFunds, missingFunds, err := FilterMinimumBalanceAccounts(accounts, minimumBalance)

	if err == nil {
		fmt.Println(fmt.Sprintf("\nThe following keys hold sufficient funds >%f:", minimumBalance))
		for keyName, address := range hasFunds {
			fmt.Println(fmt.Sprintf("%s / %s", keyName, address))
		}

		fmt.Println(fmt.Sprintf("\nThe following keys don't hold sufficient funds of >%f:", minimumBalance))
		for keyName, address := range missingFunds {
			fmt.Println(fmt.Sprintf("%s / %s", keyName, address))
		}
	}
}
