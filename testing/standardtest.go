package testing

import (
	"fmt"

	"github.com/SebastianJ/harmony-tx-tests/accounts"
	"github.com/SebastianJ/harmony-tx-tests/balances"
	"github.com/SebastianJ/harmony-tx-tests/config"
	"github.com/SebastianJ/harmony-tx-tests/transactions"
)

// RunStandardTestCase - executes a standard/simple test case
func RunStandardTestCase(testCase TestCase) {
	Title(testCase.Name, "header", testCase.Verbose)
	Log(testCase.Name, fmt.Sprintf("Using sender address: %s", config.Configuration.Funding.Account.Address), testCase.Verbose)

	sinkAccountName := fmt.Sprintf("%s_receiver", config.Configuration.Funding.Account.Address)
	Log(testCase.Name, fmt.Sprintf("Generating a new receiver account: %s", sinkAccountName), testCase.Verbose)
	toAddress, err := accounts.GenerateAccountAndReturnAddress(sinkAccountName)

	senderStartingBalance, _ := balances.GetShardBalance(config.Configuration.Funding.Account.Address, testCase.Parameters.FromShardID)
	receiverStartingBalance, _ := balances.GetShardBalance(toAddress, testCase.Parameters.ToShardID)

	Log(testCase.Name, fmt.Sprintf("Generated a new receiver account: %s, address: %s", sinkAccountName, toAddress), testCase.Verbose)
	Log(testCase.Name, fmt.Sprintf("Using sender address: %s and receiver address : %s", config.Configuration.Funding.Account.Address, toAddress), testCase.Verbose)
	Log(testCase.Name, fmt.Sprintf("Sender address: %s has a starting balance of %f in shard %d before the test", config.Configuration.Funding.Account.Address, senderStartingBalance, testCase.Parameters.FromShardID), testCase.Verbose)
	Log(testCase.Name, fmt.Sprintf("Receiver account %s (address: %s) has a starting balance of %f in shard %d before the test", sinkAccountName, toAddress, receiverStartingBalance, testCase.Parameters.ToShardID), testCase.Verbose)
	Log(testCase.Name, fmt.Sprintf("Will let the transaction wait up to %d seconds to try to get finalized within 2 blocks", testCase.Parameters.ConfirmationWaitTime), testCase.Verbose)
	Log(testCase.Name, "Sending transaction...", testCase.Verbose)

	rawTx, err := transactions.SendTransaction(config.Configuration.Funding.Account.Address, testCase.Parameters.FromShardID, toAddress, testCase.Parameters.ToShardID, testCase.Parameters.Amount, testCase.Parameters.Nonce, testCase.Parameters.GasPrice, testCase.Parameters.Data, testCase.Parameters.ConfirmationWaitTime)
	testCaseTx := ConvertToTestCaseTransaction(config.Configuration.Funding.Account.Address, toAddress, rawTx, testCase.Parameters, err)
	testCase.Transactions = append(testCase.Transactions, testCaseTx)

	Log(testCase.Name, fmt.Sprintf("Sent %f token(s) from %s to %s - transaction hash: %s, tx successful: %t", testCase.Parameters.Amount, config.Configuration.Funding.Account.Address, toAddress, testCaseTx.TransactionHash, testCaseTx.Success), testCase.Verbose)

	senderEndingBalance, _ := balances.GetShardBalance(config.Configuration.Funding.Account.Address, testCase.Parameters.FromShardID)
	receiverEndingBalance, _ := balances.GetShardBalance(toAddress, testCase.Parameters.ToShardID)

	Log(testCase.Name, fmt.Sprintf("Sender address: %s has an ending balance of %f in shard %d after the test", config.Configuration.Funding.Account.Address, senderEndingBalance, testCase.Parameters.FromShardID), testCase.Verbose)
	Log(testCase.Name, fmt.Sprintf("Receiver address: %s has an ending balance of %f in shard %d after the test", toAddress, receiverEndingBalance, testCase.Parameters.ToShardID), testCase.Verbose)
	Log(testCase.Name, "Performing test teardown (returning funds and removing sink account)", testCase.Verbose)
	Title(testCase.Name, "footer", testCase.Verbose)

	Teardown(sinkAccountName, toAddress, testCase.Parameters.FromShardID, config.Configuration.Funding.Account.Address, testCase.Parameters.ToShardID, testCase.Parameters.Amount, testCase.Parameters.GasPrice, 0)

	testCase.Result = testCaseTx.Success

	//testCase.Result = (testCaseTx.Success && ((receiverStartingBalance)+testCase.Parameters.Amount == receiverEndingBalance))

	Results = append(Results, testCase)
}
