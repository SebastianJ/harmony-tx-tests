package testing

import (
	"github.com/SebastianJ/harmony-tx-tests/accounts"
	"github.com/SebastianJ/harmony-tx-tests/config"
	"github.com/SebastianJ/harmony-tx-tests/transactions"
)

// Teardown - return any sent tokens (minus a gas cost) and remove the account from the keystore
func Teardown(accountName string, fromAddress string, fromShardID uint32, toAddress string, toShardID uint32, amount float64, gasPrice int64, confirmationWaitTime int) {
	if amount > 0.0 {
		amount = (amount - config.Configuration.Network.GasCost)
	}

	transactions.SendTransaction(fromAddress, fromShardID, toAddress, toShardID, amount, -1, gasPrice, "", confirmationWaitTime)
	accounts.RemoveAccount(accountName)
}
