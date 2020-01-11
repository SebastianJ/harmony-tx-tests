package testcases

import (
	"fmt"

	"github.com/SebastianJ/harmony-tx-tests/accounts"
	"github.com/SebastianJ/harmony-tx-tests/balances"
	"github.com/SebastianJ/harmony-tx-tests/testing"
	"github.com/SebastianJ/harmony-tx-tests/transactions"
	"github.com/SebastianJ/harmony-tx-tests/utils"
)

// Common test parameters are defined here - e.g. the test case name, expected result of the test and the required parameters to run the test case
var sbs1TestCase testing.TestCase = testing.TestCase{
	Scenario: "Same Beacon Shard",
	Name:     "SBS1",
	Goal:     "Single account",
	Priority: 0,
	Expected: true,
	Verbose:  true,
	Parameters: testing.TestCaseParameters{
		FromShardID:          0,
		ToShardID:            0,
		Data:                 "",
		Amount:               1.0,
		GasPrice:             0,
		Count:                1,
		ConfirmationWaitTime: 16,
	},
}

// Sbs1TestCase - Same Beacon Shard single account transfer A -> B, Shard 0 -> 0, Amount 1, Tx Data nil, expects: successful token transferred from A to B within 2 blocks time 16s
func Sbs1TestCase(accs map[string]string, passphrase string, node string) testing.TestCase {
	keyName, fromAddress := utils.RandomItemFromMap(accs)

	testing.Title(sbs1TestCase.Name, "header", sbs1TestCase.Verbose)
	testing.Log(sbs1TestCase.Name, fmt.Sprintf("Using source/sender key: %s and address: %s", keyName, fromAddress), sbs1TestCase.Verbose)

	sinkAccountName := fmt.Sprintf("%s_sink", keyName)
	testing.Log(sbs1TestCase.Name, fmt.Sprintf("Generating a new receiver/sink account: %s", sinkAccountName), sbs1TestCase.Verbose)
	toAddress, err := accounts.GenerateAccountAndReturnAddress(sinkAccountName, passphrase)

	senderStartingBalance, _ := balances.GetShardBalance(fromAddress, sbs1TestCase.Parameters.FromShardID, node)
	receiverStartingBalance, _ := balances.GetShardBalance(toAddress, sbs1TestCase.Parameters.ToShardID, node)

	testing.Log(sbs1TestCase.Name, fmt.Sprintf("Generated a new receiver/sink account: %s, address: %s", sinkAccountName, toAddress), sbs1TestCase.Verbose)
	testing.Log(sbs1TestCase.Name, fmt.Sprintf("Using source account %s (address: %s) and sink account %s (address : %s)", keyName, fromAddress, sinkAccountName, toAddress), sbs1TestCase.Verbose)
	testing.Log(sbs1TestCase.Name, fmt.Sprintf("Source account %s (address: %s) has a starting balance of %f in shard %d before the test", keyName, fromAddress, senderStartingBalance, sbs1TestCase.Parameters.FromShardID), sbs1TestCase.Verbose)
	testing.Log(sbs1TestCase.Name, fmt.Sprintf("Sink account %s (address: %s) has a starting balance of %f in shard %d before the test", sinkAccountName, toAddress, receiverStartingBalance, sbs1TestCase.Parameters.ToShardID), sbs1TestCase.Verbose)
	testing.Log(sbs1TestCase.Name, fmt.Sprintf("Will let the transaction wait up to %d seconds to try to get finalized within 2 blocks", sbs1TestCase.Parameters.ConfirmationWaitTime), sbs1TestCase.Verbose)
	testing.Log(sbs1TestCase.Name, "Sending transaction...", sbs1TestCase.Verbose)

	rawTx, err := transactions.SendTransaction(fromAddress, sbs1TestCase.Parameters.FromShardID, toAddress, sbs1TestCase.Parameters.ToShardID, sbs1TestCase.Parameters.Amount, sbs1TestCase.Parameters.GasPrice, sbs1TestCase.Parameters.Data, passphrase, node, sbs1TestCase.Parameters.ConfirmationWaitTime)
	testCaseTx := testing.ConvertToTestCaseTransaction(fromAddress, toAddress, rawTx, sbs1TestCase.Parameters, err)
	sbs1TestCase.Transactions = append(sbs1TestCase.Transactions, testCaseTx)

	testing.Log(sbs1TestCase.Name, fmt.Sprintf("Sent %f token(s) from %s to %s - transaction hash: %s, tx successful: %t", sbs1TestCase.Parameters.Amount, fromAddress, toAddress, testCaseTx.TransactionHash, testCaseTx.Success), sbs1TestCase.Verbose)

	senderEndingBalance, _ := balances.GetShardBalance(fromAddress, sbs1TestCase.Parameters.FromShardID, node)
	receiverEndingBalance, _ := balances.GetShardBalance(toAddress, sbs1TestCase.Parameters.ToShardID, node)

	testing.Log(sbs1TestCase.Name, fmt.Sprintf("Source account %s (address: %s) has an ending balance of %f in shard %d after the test", keyName, fromAddress, senderEndingBalance, sbs1TestCase.Parameters.FromShardID), sbs1TestCase.Verbose)
	testing.Log(sbs1TestCase.Name, fmt.Sprintf("Sink account %s (address: %s) has an ending balance of %f in shard %d after the test", sinkAccountName, toAddress, receiverEndingBalance, sbs1TestCase.Parameters.ToShardID), sbs1TestCase.Verbose)
	testing.Log(sbs1TestCase.Name, "Performing test teardown (returning funds and removing sink account)", sbs1TestCase.Verbose)
	testing.Title(sbs1TestCase.Name, "footer", sbs1TestCase.Verbose)

	testing.Teardown(sinkAccountName, toAddress, sbs1TestCase.Parameters.FromShardID, fromAddress, sbs1TestCase.Parameters.ToShardID, sbs1TestCase.Parameters.Amount, sbs1TestCase.Parameters.GasPrice, passphrase, node, 0)

	sbs1TestCase.Result = (testCaseTx.Success && ((receiverStartingBalance)+sbs1TestCase.Parameters.Amount == receiverEndingBalance))

	return sbs1TestCase
}
