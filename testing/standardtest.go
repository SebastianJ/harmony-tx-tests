package testing

import (
	"fmt"

	"github.com/SebastianJ/harmony-tx-tests/accounts"
	"github.com/SebastianJ/harmony-tx-tests/balances"
	"github.com/SebastianJ/harmony-tx-tests/config"
	"github.com/SebastianJ/harmony-tx-tests/funding"
	"github.com/SebastianJ/harmony-tx-tests/transactions"
)

// RunStandardTestCase - executes a standard/simple test case
func RunStandardTestCase(testCase TestCase) TestCase {
	Title(testCase.Name, "header", testCase.Verbose)

	senderAccountName := fmt.Sprintf("TestCase_%s_Sender", testCase.Name)
	Log(testCase.Name, fmt.Sprintf("Generating a new sender account: %s", senderAccountName), testCase.Verbose)
	senderAccount := accounts.GenerateTypedAccount(senderAccountName)

	fundingAccountBalance, _ := balances.GetShardBalance(config.Configuration.Funding.Account.Address, testCase.Parameters.FromShardID)
	fundingAmount := 1.0
	if (float64(testCase.Parameters.ReceiverCount) * testCase.Parameters.Amount) <= fundingAccountBalance {
		fundingAmount = (float64(testCase.Parameters.ReceiverCount) * testCase.Parameters.Amount)
	} else {
		fundingAmount = (float64(testCase.Parameters.ReceiverCount) * fundingAmount)
	}

	if testCase.Parameters.Data == "" && testCase.Parameters.DataSize > 0 {
		testCase.Parameters.Data = transactions.GenerateTxData(testCase.Parameters.DataSize)
	}

	fmt.Println(fmt.Sprintf("Tx data is now %s", testCase.Parameters.Data))
	fmt.Println(fmt.Sprintf("Tx data size is %d", testCase.Parameters.DataSize))

	Log(testCase.Name, fmt.Sprintf("Funding sender account: %s, address: %s", senderAccount.Name, senderAccount.Address), testCase.Verbose)
	funding.PerformFundingTransaction(config.Configuration.Funding.Account.Address, testCase.Parameters.FromShardID, senderAccount.Address, testCase.Parameters.ToShardID, fundingAmount, -1, config.Configuration.Funding.GasPrice, config.Configuration.Funding.ConfirmationWaitTime, config.Configuration.Funding.Attempts)

	receiverAccountName := fmt.Sprintf("TestCase_%s_Receiver", testCase.Name)
	Log(testCase.Name, fmt.Sprintf("Generating a new receiver account: %s", receiverAccountName), testCase.Verbose)
	receiverAccount := accounts.GenerateTypedAccount(receiverAccountName)

	senderStartingBalance, _ := balances.GetShardBalance(senderAccount.Address, testCase.Parameters.FromShardID)
	receiverStartingBalance, _ := balances.GetShardBalance(receiverAccount.Address, testCase.Parameters.ToShardID)

	Log(testCase.Name, fmt.Sprintf("Using sender account %s, address: %s and receiver account %s, address : %s", senderAccount.Name, senderAccount.Address, receiverAccount.Name, receiverAccount.Address), testCase.Verbose)
	Log(testCase.Name, fmt.Sprintf("Sender account %s, address: %s has a starting balance of %f in shard %d before the test", senderAccount.Name, senderAccount.Address, senderStartingBalance, testCase.Parameters.FromShardID), testCase.Verbose)
	Log(testCase.Name, fmt.Sprintf("Receiver account %s, address: %s has a starting balance of %f in shard %d before the test", receiverAccount.Name, receiverAccount.Address, receiverStartingBalance, testCase.Parameters.ToShardID), testCase.Verbose)
	Log(testCase.Name, fmt.Sprintf("Will let the transaction wait up to %d seconds to try to get finalized within 2 blocks", testCase.Parameters.ConfirmationWaitTime), testCase.Verbose)
	Log(testCase.Name, "Sending transaction...", testCase.Verbose)

	rawTx, err := transactions.SendTransaction(senderAccount.Address, testCase.Parameters.FromShardID, receiverAccount.Address, testCase.Parameters.ToShardID, testCase.Parameters.Amount, testCase.Parameters.Nonce, testCase.Parameters.GasPrice, testCase.Parameters.Data, testCase.Parameters.ConfirmationWaitTime)
	testCaseTx := ConvertToTestCaseTransaction(senderAccount.Address, receiverAccount.Address, rawTx, testCase.Parameters, err)
	testCase.Transactions = append(testCase.Transactions, testCaseTx)

	Log(testCase.Name, fmt.Sprintf("Sent %f token(s) from %s to %s - transaction hash: %s, tx successful: %t", testCase.Parameters.Amount, senderAccount.Address, receiverAccount.Address, testCaseTx.TransactionHash, testCaseTx.Success), testCase.Verbose)

	senderEndingBalance, _ := balances.GetShardBalance(senderAccount.Address, testCase.Parameters.FromShardID)
	receiverEndingBalance, _ := balances.GetShardBalance(receiverAccount.Address, testCase.Parameters.ToShardID)

	Log(testCase.Name, fmt.Sprintf("Sender address: %s has an ending balance of %f in shard %d after the test", senderAccount.Address, senderEndingBalance, testCase.Parameters.FromShardID), testCase.Verbose)
	Log(testCase.Name, fmt.Sprintf("Receiver address: %s has an ending balance of %f in shard %d after the test", receiverAccount.Address, receiverEndingBalance, testCase.Parameters.ToShardID), testCase.Verbose)
	Log(testCase.Name, "Performing test teardown (returning funds and removing sink account)", testCase.Verbose)
	Title(testCase.Name, "footer", testCase.Verbose)

	Teardown(senderAccount.Name, senderAccount.Address, testCase.Parameters.FromShardID, config.Configuration.Funding.Account.Address, testCase.Parameters.ToShardID, testCase.Parameters.Amount, testCase.Parameters.GasPrice, 0)
	Teardown(receiverAccount.Name, receiverAccount.Address, testCase.Parameters.FromShardID, config.Configuration.Funding.Account.Address, testCase.Parameters.ToShardID, testCase.Parameters.Amount, testCase.Parameters.GasPrice, 0)

	testCase.Result = testCaseTx.Success

	//testCase.Result = (testCaseTx.Success && ((receiverStartingBalance)+testCase.Parameters.Amount == receiverEndingBalance))

	return testCase
}
