package testcases

import (
	"fmt"
	"sync"

	"github.com/SebastianJ/harmony-tx-tests/accounts"
	"github.com/SebastianJ/harmony-tx-tests/balances"
	"github.com/SebastianJ/harmony-tx-tests/config"
	"github.com/SebastianJ/harmony-tx-tests/funding"
	"github.com/SebastianJ/harmony-tx-tests/testing"
	"github.com/SebastianJ/harmony-tx-tests/transactions"
)

// Sbs6TestCase - defines the common properties for the SBS6 test case
var Sbs6TestCase testing.TestCase = testing.TestCase{
	Scenario: "Same Beacon Shard",
	Name:     "SBS6",
	Goal:     "Multiple accounts",
	Priority: 0,
	Expected: true,
	Verbose:  true,
	Parameters: testing.TestCaseParameters{
		FromShardID:          0,
		ToShardID:            0,
		SenderCount:          10,
		ReceiverCount:        1,
		Data:                 "",
		Amount:               1.0,
		GasPrice:             1,
		Count:                1,
		ConfirmationWaitTime: 16,
	},
}

// RunSbs6TestCase - Same Beacon Shard single account transfer A1 - A10 -> B, Shard 0 -> 0, Amount 1, Tx Data nil, expects: successful token transfer from A1 - A10 to B within 2 blocks time 16s
func RunSbs6TestCase() {
	testing.Title(Sbs6TestCase.Name, "header", Sbs6TestCase.Verbose)

	senderAccounts, _ := generateAndFundSenderAccounts()

	responses := make(chan bool, Sbs6TestCase.Parameters.SenderCount)
	var waitGroup sync.WaitGroup

	for _, senderAccount := range senderAccounts {
		waitGroup.Add(1)
		go performSingleAccountTest(senderAccount, responses, &waitGroup)
	}

	waitGroup.Wait()
	close(responses)

	successfulCount := 0
	for response := range responses {
		if response {
			successfulCount++
		}
	}

	Sbs6TestCase.Result = (successfulCount == Sbs6TestCase.Parameters.SenderCount)

	testing.Title(Sbs6TestCase.Name, "footer", Sbs6TestCase.Verbose)

	Results = append(Results, Sbs6TestCase)
}

func performSingleAccountTest(senderAccount accounts.Account, responses chan<- bool, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()

	testing.Log(Sbs6TestCase.Name, fmt.Sprintf("Using sender account: %s and address: %s", senderAccount.Name, senderAccount.Address), Sbs6TestCase.Verbose)

	receiverAccountName := fmt.Sprintf("%s_Receiver", senderAccount.Name)
	testing.Log(Sbs6TestCase.Name, fmt.Sprintf("Generating a new receiver account: %s", receiverAccountName), Sbs6TestCase.Verbose)
	toAddress, err := accounts.GenerateAccountAndReturnAddress(receiverAccountName)

	senderStartingBalance, _ := balances.GetShardBalance(senderAccount.Address, Sbs6TestCase.Parameters.FromShardID)
	receiverStartingBalance, _ := balances.GetShardBalance(toAddress, Sbs6TestCase.Parameters.ToShardID)

	testing.Log(Sbs6TestCase.Name, fmt.Sprintf("Generated a new receiver account: %s, address: %s", receiverAccountName, toAddress), Sbs6TestCase.Verbose)
	testing.Log(Sbs6TestCase.Name, fmt.Sprintf("Using source account %s (address: %s) and sink account %s (address : %s)", senderAccount.Name, senderAccount.Address, receiverAccountName, toAddress), Sbs6TestCase.Verbose)
	testing.Log(Sbs6TestCase.Name, fmt.Sprintf("Source account %s (address: %s) has a starting balance of %f in shard %d before the test", senderAccount.Name, senderAccount.Address, senderStartingBalance, Sbs6TestCase.Parameters.FromShardID), Sbs6TestCase.Verbose)
	testing.Log(Sbs6TestCase.Name, fmt.Sprintf("Sink account %s (address: %s) has a starting balance of %f in shard %d before the test", receiverAccountName, toAddress, receiverStartingBalance, Sbs6TestCase.Parameters.ToShardID), Sbs6TestCase.Verbose)
	testing.Log(Sbs6TestCase.Name, fmt.Sprintf("Will let the transaction wait up to %d seconds to try to get finalized within 2 blocks", Sbs6TestCase.Parameters.ConfirmationWaitTime), Sbs6TestCase.Verbose)
	testing.Log(Sbs6TestCase.Name, "Sending transaction...", Sbs6TestCase.Verbose)

	rawTx, err := transactions.SendTransaction(senderAccount.Address, Sbs6TestCase.Parameters.FromShardID, toAddress, Sbs6TestCase.Parameters.ToShardID, Sbs6TestCase.Parameters.Amount, -1, Sbs6TestCase.Parameters.GasPrice, Sbs6TestCase.Parameters.Data, Sbs6TestCase.Parameters.ConfirmationWaitTime)
	testCaseTx := testing.ConvertToTestCaseTransaction(senderAccount.Address, toAddress, rawTx, Sbs6TestCase.Parameters, err)
	Sbs6TestCase.Transactions = append(Sbs6TestCase.Transactions, testCaseTx)

	testing.Log(Sbs6TestCase.Name, fmt.Sprintf("Sent %f token(s) from %s to %s - transaction hash: %s, tx successful: %t", Sbs6TestCase.Parameters.Amount, config.Configuration.Funding.Account.Address, toAddress, testCaseTx.TransactionHash, testCaseTx.Success), Sbs6TestCase.Verbose)

	senderEndingBalance, _ := balances.GetShardBalance(senderAccount.Address, Sbs6TestCase.Parameters.FromShardID)
	receiverEndingBalance, _ := balances.GetShardBalance(toAddress, Sbs6TestCase.Parameters.ToShardID)

	testing.Log(Sbs6TestCase.Name, fmt.Sprintf("Source account %s (address: %s) has an ending balance of %f in shard %d after the test", senderAccount.Name, senderAccount.Address, senderEndingBalance, Sbs6TestCase.Parameters.FromShardID), Sbs6TestCase.Verbose)
	testing.Log(Sbs6TestCase.Name, fmt.Sprintf("Sink account %s (address: %s) has an ending balance of %f in shard %d after the test", receiverAccountName, toAddress, receiverEndingBalance, Sbs6TestCase.Parameters.ToShardID), Sbs6TestCase.Verbose)
	testing.Log(Sbs6TestCase.Name, "Performing test teardown (returning funds and removing sink account)", Sbs6TestCase.Verbose)

	testing.Teardown(receiverAccountName, toAddress, Sbs6TestCase.Parameters.FromShardID, config.Configuration.Funding.Account.Address, Sbs6TestCase.Parameters.ToShardID, Sbs6TestCase.Parameters.Amount, Sbs6TestCase.Parameters.GasPrice, 0)
	testing.Teardown(senderAccount.Name, senderAccount.Address, Sbs6TestCase.Parameters.FromShardID, config.Configuration.Funding.Account.Address, Sbs6TestCase.Parameters.ToShardID, Sbs6TestCase.Parameters.Amount, Sbs6TestCase.Parameters.GasPrice, 0)

	responses <- (testCaseTx.Success && ((receiverStartingBalance)+Sbs1TestCase.Parameters.Amount == receiverEndingBalance))
}

func generateAndFundSenderAccounts() (senderAccounts []accounts.Account, err error) {
	networkHandler, err := transactions.NetworkHandler(Sbs6TestCase.Parameters.FromShardID)
	if err != nil {
		return nil, err
	}

	nonce := -1
	receivedNonce, err := transactions.CurrentNonce(config.Configuration.Funding.Account.Address, networkHandler)
	if err != nil {
		return nil, err
	}
	nonce = int(receivedNonce)

	amount := (Sbs6TestCase.Parameters.Amount + config.Configuration.Network.GasCost)

	var waitGroup sync.WaitGroup

	senderAccountsChannel := make(chan accounts.Account, Sbs6TestCase.Parameters.SenderCount)

	for i := 0; i < Sbs6TestCase.Parameters.SenderCount; i++ {
		waitGroup.Add(1)
		go generateAndFundSenderAccount(i, amount, nonce, senderAccountsChannel, &waitGroup)
		nonce++
	}

	waitGroup.Wait()
	close(senderAccountsChannel)

	for senderAccount := range senderAccountsChannel {
		senderAccounts = append(senderAccounts, senderAccount)
	}

	return senderAccounts, nil
}

func generateAndFundSenderAccount(index int, amount float64, nonce int, senderAccountsChannel chan<- accounts.Account, waitGroup *sync.WaitGroup) {
	defer waitGroup.Done()

	accountName, address := generateSenderAccount(index)
	funding.PerformFundingTransaction(config.Configuration.Funding.Account.Address, Sbs6TestCase.Parameters.FromShardID, address, Sbs6TestCase.Parameters.ToShardID, amount, nonce, config.Configuration.Funding.GasPrice, config.Configuration.Funding.ConfirmationWaitTime, config.Configuration.Funding.Attempts)

	senderAccountsChannel <- accounts.Account{Name: accountName, Address: address}
}

func generateSenderAccount(index int) (string, string) {
	accountName := fmt.Sprintf("TestCase_%s_Sender_%d", Sbs6TestCase.Name, index)
	testing.Log(Sbs6TestCase.Name, fmt.Sprintf("Generating a new sender account: %s", accountName), Sbs6TestCase.Verbose)
	toAddress, _ := accounts.GenerateAccountAndReturnAddress(accountName)
	return accountName, toAddress
}
