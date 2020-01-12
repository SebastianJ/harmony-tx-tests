package testing

import (
	"fmt"

	"github.com/SebastianJ/harmony-tx-tests/accounts"
	"github.com/SebastianJ/harmony-tx-tests/transactions"
)

// FundAccounts - funds a given set of accounts in a given set of shards using a set of source accounts
func FundAccounts(sources []string, shards []uint32, count int, amount float64, prefix string, gasPrice int64, passphrase string, node string, confirmationWaitTime int) (map[string]string, error) {
	fundedAccounts := make(map[string]string)
	
	for _, source := range sources {
		for i := 0; i < count; i++ {
			accountName, accountAddress, err := fundAccount(source, shards, amount, prefix, i, gasPrice, passphrase, node, confirmationWaitTime)

			if err != nil {
				return nil, err
			}

			if accountName != "" && accountAddress != "" {
				fundedAccounts[accountName] = accountAddress
			}
		}
	}

	return fundedAccounts, nil
}

func fundAccount(source string, shards []uint32, amount float64, prefix string, index int, gasPrice int64, passphrase string, node string, confirmationWaitTime int) (string, string, error) {
	accountName := fmt.Sprintf("%s_%d", prefix, index)

	// Remove the account just to make sure that we're starting using a clean slate
	accounts.RemoveAccount(accountName)
	
	err := accounts.GenerateAccount(accountName, passphrase)

	if err != nil {
		return "", "", err
	}

	sourceAddress := accounts.FindAccountAddressByName(source)
	toAddress, err := accounts.GenerateAccountAndReturnAddress(accountName, passphrase)

	for _, shard := range shards {
		success := performFundingTransaction(sourceAddress, 0, toAddress, shard, amount, gasPrice, passphrase, node, confirmationWaitTime, 10)

		if !success {
			return "", "", fmt.Errorf("failed to fund account %s on shard %d with amount %f", toAddress, shard, amount)
		}
	}

	return accountName, toAddress, nil
}

func performFundingTransaction(fromAddress string, fromShardID uint32, toAddress string, toShardID uint32, amount float64, gasPrice int64, passphrase string, node string, confirmationWaitTime int, attempts int) bool {
	success := false

	for ok := true; ok; ok = !success {
		attempts--

		if attempts > 0 {
			rawTx, err := transactions.SendTransaction(fromAddress, fromShardID, toAddress, toShardID, amount, gasPrice, "", passphrase, node, confirmationWaitTime)

			if err != nil {
				success = false
			} else {
				success = transactions.IsTransactionSuccessful(rawTx)
			}
		} else {
			ok = false
		}
	}

	return success
}
