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
var sbs4TestCase testing.TestCase = testing.TestCase{
	Scenario: "Same Beacon Shard",
	Name:     "SBS4",
	Goal:     "Nano transfer",
	Priority: 0,
	Expected: true,
	Verbose:  true,
	Parameters: testing.TestCaseParameters{
		FromShardID:          0,
		ToShardID:            0,
		Data:                 "",
		Amount:               1.00E-09,
		GasPrice:             0,
		Count:                1,
		ConfirmationWaitTime: 16,
	},
}

// Sbs4TestCase - Same Beacon Shard single account transfer A -> B, Shard 0 -> 0, Amount 1.00E-09, Tx Data nil, expects: successful token transferred from A to B within 2 blocks time 16s
func Sbs4TestCase(accs map[string]string, passphrase string, node string) testing.TestCase {
	keyName, fromAddress := utils.RandomItemFromMap(accs)

	testing.Title(sbs4TestCase.Name, "header", sbs4TestCase.Verbose)
	testing.Log(sbs4TestCase.Name, fmt.Sprintf("Using source/sender key: %s and address: %s", keyName, fromAddress), sbs4TestCase.Verbose)

	sinkAccountName := fmt.Sprintf("%s_sink", keyName)
	testing.Log(sbs4TestCase.Name, fmt.Sprintf("Generating a new receiver/sink account: %s", sinkAccountName), sbs4TestCase.Verbose)
	toAddress, err := accounts.GenerateAccountAndReturnAddress(sinkAccountName, passphrase)

	senderStartingBalance, _ := balances.GetShardBalance(fromAddress, sbs4TestCase.Parameters.FromShardID, node)
	receiverStartingBalance, _ := balances.GetShardBalance(toAddress, sbs4TestCase.Parameters.ToShardID, node)

	testing.Log(sbs4TestCase.Name, fmt.Sprintf("Generated a new receiver/sink account: %s, address: %s", sinkAccountName, toAddress), sbs4TestCase.Verbose)
	testing.Log(sbs4TestCase.Name, fmt.Sprintf("Using source account %s (address: %s) and sink account %s (address : %s)", keyName, fromAddress, sinkAccountName, toAddress), sbs4TestCase.Verbose)
	testing.Log(sbs4TestCase.Name, fmt.Sprintf("Source account %s (address: %s) has a starting balance of %f in shard %d before the test", keyName, fromAddress, senderStartingBalance, sbs4TestCase.Parameters.FromShardID), sbs4TestCase.Verbose)
	testing.Log(sbs4TestCase.Name, fmt.Sprintf("Sink account %s (address: %s) has a starting balance of %f in shard %d before the test", sinkAccountName, toAddress, receiverStartingBalance, sbs4TestCase.Parameters.ToShardID), sbs4TestCase.Verbose)
	testing.Log(sbs4TestCase.Name, fmt.Sprintf("Will let the transaction wait up to %d seconds to try to get finalized within 2 blocks", sbs4TestCase.Parameters.ConfirmationWaitTime), sbs4TestCase.Verbose)
	testing.Log(sbs4TestCase.Name, "Sending transaction...", sbs4TestCase.Verbose)

	rawTx, err := transactions.SendTransaction(fromAddress, sbs4TestCase.Parameters.FromShardID, toAddress, sbs4TestCase.Parameters.ToShardID, sbs4TestCase.Parameters.Amount, sbs4TestCase.Parameters.GasPrice, sbs4TestCase.Parameters.Data, passphrase, node, sbs4TestCase.Parameters.ConfirmationWaitTime)
	testCaseTx := testing.ConvertToTestCaseTransaction(fromAddress, toAddress, rawTx, sbs4TestCase.Parameters, err)
	sbs4TestCase.Transactions = append(sbs4TestCase.Transactions, testCaseTx)

	testing.Log(sbs4TestCase.Name, fmt.Sprintf("Sent %f token(s) from %s to %s - transaction hash: %s, tx successful: %t", sbs4TestCase.Parameters.Amount, fromAddress, toAddress, testCaseTx.TransactionHash, testCaseTx.Success), sbs4TestCase.Verbose)

	senderEndingBalance, _ := balances.GetShardBalance(fromAddress, sbs4TestCase.Parameters.FromShardID, node)
	receiverEndingBalance, _ := balances.GetShardBalance(toAddress, sbs4TestCase.Parameters.ToShardID, node)

	testing.Log(sbs4TestCase.Name, fmt.Sprintf("Source account %s (address: %s) has an ending balance of %f in shard %d after the test", keyName, fromAddress, senderEndingBalance, sbs4TestCase.Parameters.FromShardID), sbs4TestCase.Verbose)
	testing.Log(sbs4TestCase.Name, fmt.Sprintf("Sink account %s (address: %s) has an ending balance of %f in shard %d after the test", sinkAccountName, toAddress, receiverEndingBalance, sbs4TestCase.Parameters.ToShardID), sbs4TestCase.Verbose)
	testing.Log(sbs4TestCase.Name, "Performing test teardown (returning funds and removing sink account)", sbs4TestCase.Verbose)
	testing.Title(sbs4TestCase.Name, "footer", sbs4TestCase.Verbose)

	testing.Teardown(sinkAccountName, toAddress, sbs4TestCase.Parameters.FromShardID, fromAddress, sbs4TestCase.Parameters.ToShardID, sbs4TestCase.Parameters.Amount, sbs4TestCase.Parameters.GasPrice, passphrase, node, 0)

	//TODO: tests with super small denominations can't rely on the balances supplied from balances.GetShardBalance since that function only returns 6 decimals
	sbs4TestCase.Result = testCaseTx.Success

	return sbs4TestCase
}
