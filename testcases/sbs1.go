package testcases

import (
	"fmt"

	"github.com/SebastianJ/harmony-tx-tests/accounts"
	"github.com/SebastianJ/harmony-tx-tests/balances"
	"github.com/SebastianJ/harmony-tx-tests/transactions"
	"github.com/SebastianJ/harmony-tx-tests/utils"
)

// Common test parameters are defined here - e.g. the test case name, expected result of the test and the required parameters to run the test case
var testCase TestCase = TestCase{
	Scenario: "Same Beacon Shard",
	Name:     "SBS1",
	Goal:     "Single account",
	Priority: 0,
	Expected: true,
	Verbose:  true,
	Parameters: TestCaseParameters{
		FromShardID:          0,
		ToShardID:            0,
		Data:                 "",
		Amount:               1.0,
		GasPrice:             1,
		Count:                1,
		ConfirmationWaitTime: 16,
	},
}

// Sbs1TestCase - Same Beacon Shard single account transfer A -> B, Shard 0 -> 0, Amount 1, Tx Data nil, expects: successful token transferred from A to B within 2 blocks time 16s
func Sbs1TestCase(accs map[string]string, passphrase string, node string) TestCase {
	keyName, fromAddress := utils.RandomItemFromMap(accs)

	TestTitle(testCase.Name, "header", testCase.Verbose)
	TestLog(testCase.Name, fmt.Sprintf("Using source/sender key: %s and address: %s", keyName, fromAddress), testCase.Verbose)

	sinkAccountName := fmt.Sprintf("%s_sink", keyName)
	TestLog(testCase.Name, fmt.Sprintf("Generating a new receiver/sink account: %s", sinkAccountName), testCase.Verbose)
	toAddress, err := accounts.GenerateAccountAndReturnAddress(sinkAccountName, passphrase)

	senderStartingBalance, _ := balances.GetShardBalance(fromAddress, testCase.Parameters.FromShardID, node)
	receiverStartingBalance, _ := balances.GetShardBalance(toAddress, testCase.Parameters.ToShardID, node)

	TestLog(testCase.Name, fmt.Sprintf("Generated a new receiver/sink account: %s, address: %s", sinkAccountName, toAddress), testCase.Verbose)
	TestLog(testCase.Name, fmt.Sprintf("Using source account %s (address: %s) and sink account %s (address : %s)", keyName, fromAddress, sinkAccountName, toAddress), testCase.Verbose)
	TestLog(testCase.Name, fmt.Sprintf("Source account %s (address: %s) has a starting balance of %f in shard %d before the test", keyName, fromAddress, senderStartingBalance, testCase.Parameters.FromShardID), testCase.Verbose)
	TestLog(testCase.Name, fmt.Sprintf("Sink account %s (address: %s) has a starting balance of %f in shard %d before the test", sinkAccountName, toAddress, receiverStartingBalance, testCase.Parameters.ToShardID), testCase.Verbose)
	TestLog(testCase.Name, fmt.Sprintf("Will let the transaction wait up to %d seconds to try to get finalized within 2 blocks", testCase.Parameters.ConfirmationWaitTime), testCase.Verbose)
	TestLog(testCase.Name, "Sending transaction...", testCase.Verbose)

	rawTx, err := transactions.SendTransaction(fromAddress, testCase.Parameters.FromShardID, toAddress, testCase.Parameters.ToShardID, testCase.Parameters.Amount, testCase.Parameters.GasPrice, testCase.Parameters.Data, passphrase, node, testCase.Parameters.ConfirmationWaitTime)
	testCaseTx := ConvertToTestCaseTransaction(fromAddress, toAddress, rawTx, testCase.Parameters, err)
	testCase.Transactions = append(testCase.Transactions, testCaseTx)

	TestLog(testCase.Name, fmt.Sprintf("Sent %f token(s) from %s to %s - transaction hash: %s, tx successful: %t", testCase.Parameters.Amount, fromAddress, toAddress, testCaseTx.TransactionHash, testCaseTx.Success), testCase.Verbose)

	senderEndingBalance, _ := balances.GetShardBalance(fromAddress, testCase.Parameters.FromShardID, node)
	receiverEndingBalance, _ := balances.GetShardBalance(toAddress, testCase.Parameters.ToShardID, node)

	TestLog(testCase.Name, fmt.Sprintf("Source account %s (address: %s) has an ending balance of %f in shard %d after the test", keyName, fromAddress, senderEndingBalance, testCase.Parameters.FromShardID), testCase.Verbose)
	TestLog(testCase.Name, fmt.Sprintf("Sink account %s (address: %s) has an ending balance of %f in shard %d after the test", sinkAccountName, toAddress, receiverEndingBalance, testCase.Parameters.ToShardID), testCase.Verbose)
	TestLog(testCase.Name, "Performing test teardown (returning funds and removing sink account)", testCase.Verbose)
	TestTitle(testCase.Name, "footer", testCase.Verbose)

	Teardown(sinkAccountName, toAddress, testCase.Parameters.FromShardID, fromAddress, testCase.Parameters.ToShardID, testCase.Parameters.Amount, testCase.Parameters.GasPrice, passphrase, node, 0)

	testCase.Result = (testCaseTx.Success && ((receiverStartingBalance)+testCase.Parameters.Amount == receiverEndingBalance))

	return testCase
}
