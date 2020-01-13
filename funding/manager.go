package funding

import (
	"fmt"
	"sync"

	"github.com/SebastianJ/harmony-tx-tests/config"
	"github.com/SebastianJ/harmony-tx-tests/accounts"
	"github.com/SebastianJ/harmony-tx-tests/balances"
	"github.com/SebastianJ/harmony-tx-tests/transactions"
	sdkBalances "github.com/SebastianJ/harmony-sdk/balances"
)

// SetupFundingAccount - sets up the initial funding account
func SetupFundingAccount(accs []string) (err error) {
	if config.Configuration.Funding.Account.Address == "" {
		if accounts.DoesNamedAccountExist(config.Configuration.Funding.Account.Name) {
			if resolvedAccName := accounts.FindAccountAddressByName(config.Configuration.Funding.Account.Name); resolvedAccName != "" {
				config.Configuration.Funding.Account.Address = resolvedAccName
			}
		} else {
			config.Configuration.Funding.Account.Address, err = accounts.GenerateAccountAndReturnAddress(config.Configuration.Funding.Account.Name)

			if err != nil {
				return err
			}
		}
	}

	totalBalance, err := sdkBalances.GetTotalBalance(config.Configuration.Funding.Account.Address, config.Configuration.Network.Name)
	if err != nil {
		return err
	}

	fmt.Println(fmt.Sprintf("The current balance of the funding account %s / %s is: %f", config.Configuration.Funding.Account.Name, config.Configuration.Funding.Account.Address, totalBalance))

	if totalBalance <= config.Configuration.Funding.MinimumFunds {
		var waitGroup sync.WaitGroup

		for _, address := range accs {
			for shard := 0; shard < config.Configuration.Network.Shards; shard++ {
				availableShardBalance, err := balances.GetShardBalance(address, uint32(shard))
				amount := availableShardBalance - 0.001

				if err == nil && availableShardBalance > 0.0 {
					waitGroup.Add(1)
					go asyncPerformFundingTransaction(address, uint32(shard), config.Configuration.Funding.Account.Address, uint32(shard), amount, config.Configuration.Funding.GasPrice, config.Configuration.Funding.ConfirmationWaitTime, config.Configuration.Funding.Attempts, &waitGroup)
				}
			}
		}

		waitGroup.Wait()
	}

	return nil
}

// FundAccounts - funds a given set of accounts in a given set of shards using a set of source accounts
func FundAccounts(sources []string, count int, amount float64, prefix string, gasPrice int64, confirmationWaitTime int) (map[string]string, error) {
	fundedAccounts := make(map[string]string)
	
	for _, source := range sources {
		for i := 0; i < count; i++ {
			accountName, accountAddress, err := fundAccount(source, amount, prefix, i, gasPrice, confirmationWaitTime)

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

func fundAccount(source string, amount float64, prefix string, index int, gasPrice int64, confirmationWaitTime int) (string, string, error) {
	accountName := fmt.Sprintf("%s_%d", prefix, index)

	// Remove the account just to make sure that we're starting using a clean slate
	accounts.RemoveAccount(accountName)
	
	err := accounts.GenerateAccount(accountName)

	if err != nil {
		return "", "", err
	}

	sourceAddress := accounts.FindAccountAddressByName(source)
	toAddress, err := accounts.GenerateAccountAndReturnAddress(accountName)

	for shard := 0; shard < config.Configuration.Network.Shards; shard++ {
		success := performFundingTransaction(sourceAddress, 0, toAddress, uint32(shard), amount, gasPrice, confirmationWaitTime, 10)

		if !success {
			return "", "", fmt.Errorf("failed to fund account %s on shard %d with amount %f", toAddress, shard, amount)
		}
	}

	return accountName, toAddress, nil
}

func asyncPerformFundingTransaction(fromAddress string, fromShardID uint32, toAddress string, toShardID uint32, amount float64, gasPrice int64, confirmationWaitTime int, attempts int, waitGroup *sync.WaitGroup) {
	performFundingTransaction(fromAddress, fromShardID, toAddress, toShardID, amount, gasPrice, confirmationWaitTime, attempts)

	defer waitGroup.Done()
}

func performFundingTransaction(fromAddress string, fromShardID uint32, toAddress string, toShardID uint32, amount float64, gasPrice int64, confirmationWaitTime int, attempts int) bool {
	success := false

	for ok := true; ok; ok = !success {
		attempts--

		if ok && attempts > 0 {
			fmt.Println(fmt.Sprintf("Attempting funding transaction from %s (shard: %d) to %s (shard: %d) of amount %f!", fromAddress, fromShardID, toAddress, toShardID, amount))

			rawTx, err := transactions.SendTransaction(fromAddress, fromShardID, toAddress, toShardID, amount, gasPrice, "", confirmationWaitTime)

			fmt.Println(fmt.Sprintf("rawTx: %v", rawTx))

			if err != nil {
				success = false
				fmt.Println(fmt.Sprintf("Failed to perform funding transaction from %s (shard: %d) to %s (shard: %d) of amount %f - error: %s", fromAddress, fromShardID, toAddress, toShardID, amount, err.Error))
			} else {
				success = transactions.IsTransactionSuccessful(rawTx)
				if success {
					break
					fmt.Println(fmt.Sprintf("Successfully performed funding transaction from %s (shard: %d) to %s (shard: %d) of amount %f", fromAddress, fromShardID, toAddress, toShardID, amount))
				} else {
					fmt.Println(fmt.Sprintf("Failed to perform funding transaction from %s (shard: %d) to %s (shard: %d) of amount %f", fromAddress, fromShardID, toAddress, toShardID, amount))
				}
			}
		} else {
			ok = false
		}
	}

	return success
}
