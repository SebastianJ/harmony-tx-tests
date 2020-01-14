package testing

import (
	"fmt"

	"github.com/SebastianJ/harmony-tx-tests/balances"
	"github.com/SebastianJ/harmony-tx-tests/config"
	"github.com/SebastianJ/harmony-tx-tests/transactions"
)

// RunSameAccountTestCase - executes a test case where the sender and receiver address is the same
func RunSameAccountTestCase(testCase TestCase) {
	Title(testCase.Name, "header", testCase.Verbose)

	senderStartingBalance, _ := balances.GetShardBalance(config.Configuration.Funding.Account.Address, testCase.Parameters.FromShardID)

	Log(testCase.Name, fmt.Sprintf("Using account %s (address: %s) for a self transfer", config.Configuration.Funding.Account.Name, config.Configuration.Funding.Account.Address), testCase.Verbose)
	Log(testCase.Name, fmt.Sprintf("Source account %s (address: %s) has a starting balance of %f in shard %d before the test", config.Configuration.Funding.Account.Name, config.Configuration.Funding.Account.Address, senderStartingBalance, testCase.Parameters.FromShardID), testCase.Verbose)
	Log(testCase.Name, fmt.Sprintf("Will let the transaction wait up to %d seconds to try to get finalized within 2 blocks", testCase.Parameters.ConfirmationWaitTime), testCase.Verbose)
	Log(testCase.Name, "Sending transaction...", testCase.Verbose)

	rawTx, err := transactions.SendTransaction(config.Configuration.Funding.Account.Address, testCase.Parameters.FromShardID, config.Configuration.Funding.Account.Address, testCase.Parameters.ToShardID, testCase.Parameters.Amount, testCase.Parameters.Nonce, testCase.Parameters.GasPrice, testCase.Parameters.Data, testCase.Parameters.ConfirmationWaitTime)
	testCaseTx := ConvertToTestCaseTransaction(config.Configuration.Funding.Account.Address, config.Configuration.Funding.Account.Address, rawTx, testCase.Parameters, err)
	testCase.Transactions = append(testCase.Transactions, testCaseTx)

	Log(testCase.Name, fmt.Sprintf("Sent %f token(s) from %s to %s - transaction hash: %s, tx successful: %t", testCase.Parameters.Amount, config.Configuration.Funding.Account.Address, config.Configuration.Funding.Account.Address, testCaseTx.TransactionHash, testCaseTx.Success), testCase.Verbose)

	senderEndingBalance, _ := balances.GetShardBalance(config.Configuration.Funding.Account.Address, testCase.Parameters.FromShardID)

	Log(testCase.Name, fmt.Sprintf("Account %s (address: %s) has an ending balance of %f in shard %d after the test", config.Configuration.Funding.Account.Name, config.Configuration.Funding.Account.Address, senderEndingBalance, testCase.Parameters.FromShardID), testCase.Verbose)
	Title(testCase.Name, "footer", testCase.Verbose)

	// We should end up with a lesser amount compared to the initial amount since we pay a gas fee
	testCase.Result = (testCaseTx.Success && (senderStartingBalance <= senderEndingBalance))

	Results = append(Results, testCase)
}
