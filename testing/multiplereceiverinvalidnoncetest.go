package testing

import (
	"fmt"
	"sync"

	"github.com/SebastianJ/harmony-tx-tests/accounts"
	"github.com/SebastianJ/harmony-tx-tests/balances"
	"github.com/SebastianJ/harmony-tx-tests/config"
	"github.com/SebastianJ/harmony-tx-tests/funding"
	"github.com/SebastianJ/harmony-tx-tests/transactions"
)

// RunMultipleReceiverInvalidNonceTestCase - runs a tests where multiple receiver wallets receive txs with the exact same nonce
func RunMultipleReceiverInvalidNonceTestCase(testCase TestCase) TestCase {
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

	Log(testCase.Name, fmt.Sprintf("Funding sender account: %s, address: %s", senderAccount.Name, senderAccount.Address), testCase.Verbose)
	funding.PerformFundingTransaction(config.Configuration.Funding.Account.Address, testCase.Parameters.FromShardID, senderAccount.Address, testCase.Parameters.ToShardID, fundingAmount, -1, config.Configuration.Funding.GasPrice, config.Configuration.Funding.ConfirmationWaitTime, config.Configuration.Funding.Attempts)
	
	nameTemplate := fmt.Sprintf("TestCase_%s_Receiver_", testCase.Name)
	receiverAccounts := accounts.AsyncGenerateMultipleTypedAccounts(nameTemplate, testCase.Parameters.ReceiverCount)

	networkHandler, _ := transactions.NetworkHandler(testCase.Parameters.FromShardID)
	nonce := -1
	receivedNonce, _ := transactions.CurrentNonce(senderAccount.Address, networkHandler)
	nonce = int(receivedNonce)

	Log(testCase.Name, fmt.Sprintf("Current nonce for sender account: %s, address: %s is %d", senderAccount.Name, senderAccount.Address, nonce), testCase.Verbose)

	txs := make(chan TestCaseTransaction, testCase.Parameters.ReceiverCount)
	var waitGroup sync.WaitGroup

	for _, receiverAccount := range receiverAccounts {
		waitGroup.Add(1)
		go performSingleReceiverAccountTest(testCase, senderAccount, receiverAccount, nonce, txs, &waitGroup)
	}

	waitGroup.Wait()
	close(txs)

	successfulCount := 0
	for tx := range txs {
		testCase.Transactions = append(testCase.Transactions, tx)
		if tx.Success {
			successfulCount++
		}
	}

	Log(testCase.Name, fmt.Sprintf("A total of %d transaction(s) were successful", successfulCount), testCase.Verbose)

	txsSuccessful := (successfulCount == testCase.Parameters.ReceiverCount)
	testCase.Result = !txsSuccessful && successfulCount == 1

	Teardown(senderAccount.Name, senderAccount.Address, testCase.Parameters.FromShardID, config.Configuration.Funding.Account.Address, testCase.Parameters.ToShardID, testCase.Parameters.Amount, testCase.Parameters.GasPrice, 0)

	for _, receiverAccount := range receiverAccounts {
		Teardown(receiverAccount.Name, receiverAccount.Address, testCase.Parameters.FromShardID, config.Configuration.Funding.Account.Address, testCase.Parameters.ToShardID, testCase.Parameters.Amount, testCase.Parameters.GasPrice, 0)
	}

	Title(testCase.Name, "footer", testCase.Verbose)

	return testCase
}

func performSingleReceiverAccountTest(testCase TestCase, senderAccount accounts.Account, receiverAccount accounts.Account, nonce int, responses chan<- TestCaseTransaction, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()

	senderStartingBalance, _ := balances.GetShardBalance(senderAccount.Address, testCase.Parameters.FromShardID)

	Log(testCase.Name, fmt.Sprintf("Generated a new receiver account: %s, address: %s", receiverAccount.Name, receiverAccount.Address), testCase.Verbose)
	Log(testCase.Name, fmt.Sprintf("Using sender account %s (address: %s) and receiver account %s (address : %s)", senderAccount.Name, senderAccount.Address, receiverAccount.Name, receiverAccount.Address), testCase.Verbose)
	Log(testCase.Name, fmt.Sprintf("Sender account %s (address: %s) has a starting balance of %f in shard %d before the test", senderAccount.Name, senderAccount.Address, senderStartingBalance, testCase.Parameters.FromShardID), testCase.Verbose)
	Log(testCase.Name, fmt.Sprintf("Will let the transaction wait up to %d seconds to try to get finalized within 2 blocks", testCase.Parameters.ConfirmationWaitTime), testCase.Verbose)
	Log(testCase.Name, "Sending transaction...", testCase.Verbose)

	rawTx, err := transactions.SendTransaction(senderAccount.Address, testCase.Parameters.FromShardID, receiverAccount.Address, testCase.Parameters.ToShardID, testCase.Parameters.Amount, nonce, testCase.Parameters.GasPrice, testCase.Parameters.Data, testCase.Parameters.ConfirmationWaitTime)
	testCaseTx := ConvertToTestCaseTransaction(senderAccount.Address, receiverAccount.Address, rawTx, testCase.Parameters, err)

	Log(testCase.Name, fmt.Sprintf("Sent %f token(s) from %s to %s - transaction hash: %s, tx successful: %t", testCase.Parameters.Amount, config.Configuration.Funding.Account.Address, receiverAccount.Address, testCaseTx.TransactionHash, testCaseTx.Success), testCase.Verbose)

	senderEndingBalance, _ := balances.GetShardBalance(senderAccount.Address, testCase.Parameters.FromShardID)

	Log(testCase.Name, fmt.Sprintf("Sender account %s (address: %s) has an ending balance of %f in shard %d after the test", senderAccount.Name, senderAccount.Address, senderEndingBalance, testCase.Parameters.FromShardID), testCase.Verbose)
	Log(testCase.Name, "Performing test teardown (returning funds and removing sender account)", testCase.Verbose)

	responses <- testCaseTx
}
