package balances

import (
	"fmt"
	"github.com/SebastianJ/harmony-tx-sender/balances"
)

// GetTotalBalance - gets the total balance across all shards for a given address
func GetTotalBalance(address string, node string) (float64, error) {
	shardBalances, err := balances.CheckAllShardBalances(node, address)

	if err != nil {
		return -1.0, err
	}

	totalBalance := 0.0

	for _, balance := range shardBalances {
		totalBalance += balance
	}

	return totalBalance, nil
}

// OutputBalanceStatusForAddresses - outputs which keys/accounts that hold funds and which don't
func OutputBalanceStatusForAddresses(accounts map[string]string, node string) error {
	hasFunds := make(map[string]string)
	missingFunds := make(map[string]string)

	for keyName, address := range accounts {
		shardBalances, err := balances.CheckAllShardBalances(node, address)

		if err != nil {
			return err
		}

		totalBalance := 0.0

		for shardID, balance := range shardBalances {
			fmt.Println(fmt.Sprintf("Balance in shard %d is %f", shardID, balance))
			totalBalance += balance
		}

		if totalBalance > 0.0 {
			fmt.Println(fmt.Sprintf("Keyfile with name: %s and address: %s holds a total of %f ONE", keyName, address, totalBalance))
			hasFunds[keyName] = address
		} else {
			fmt.Println(fmt.Sprintf("Keyfile with name: %s and address: %s doesn't hold any funds!", keyName, address))
			missingFunds[keyName] = address
		}
	}

	fmt.Println("\nThe following keys hold funds:")
	for keyName, address := range hasFunds {
		fmt.Println(fmt.Sprintf("%s / %s", keyName, address))
	}

	fmt.Println("\nThe following keys don't hold any funds:")
	for keyName, address := range missingFunds {
		fmt.Println(fmt.Sprintf("%s / %s", keyName, address))
	}

	return nil
}
