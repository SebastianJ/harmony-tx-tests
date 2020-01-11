package testcases

import (
	"fmt"

	"github.com/SebastianJ/harmony-tx-tests/balances"
	"github.com/SebastianJ/harmony-tx-tests/testing"
	"github.com/SebastianJ/harmony-tx-tests/transactions"
	"github.com/SebastianJ/harmony-tx-tests/utils"
)

// Common test parameters are defined here - e.g. the test case name, expected result of the test and the required parameters to run the test case
var sbs2TestCase testing.TestCase = testing.TestCase{
	Scenario: "Same Beacon Shard",
	Name:     "SBS2",
	Goal:     "Same account",
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

// Sbs2TestCase - Same Beacon Shard same account transfer A -> A, Shard 0 -> 0, Amount 1, Tx Data nil, expects: successful token transferred from A to A within 2 blocks time 16s
func Sbs2TestCase(accs map[string]string, passphrase string, node string) testing.TestCase {
	keyName, fromAddress := utils.RandomItemFromMap(accs)

	testing.Title(sbs2TestCase.Name, "header", sbs2TestCase.Verbose)

	senderStartingBalance, _ := balances.GetShardBalance(fromAddress, sbs2TestCase.Parameters.FromShardID, node)

	testing.Log(sbs2TestCase.Name, fmt.Sprintf("Using account %s (address: %s) for a self transfer", keyName, fromAddress), sbs2TestCase.Verbose)
	testing.Log(sbs2TestCase.Name, fmt.Sprintf("Source account %s (address: %s) has a starting balance of %f in shard %d before the test", keyName, fromAddress, senderStartingBalance, sbs2TestCase.Parameters.FromShardID), sbs2TestCase.Verbose)
	testing.Log(sbs2TestCase.Name, fmt.Sprintf("Will let the transaction wait up to %d seconds to try to get finalized within 2 blocks", sbs2TestCase.Parameters.ConfirmationWaitTime), sbs2TestCase.Verbose)
	testing.Log(sbs2TestCase.Name, "Sending transaction...", sbs2TestCase.Verbose)

	rawTx, err := transactions.SendTransaction(fromAddress, sbs2TestCase.Parameters.FromShardID, fromAddress, sbs2TestCase.Parameters.ToShardID, sbs2TestCase.Parameters.Amount, sbs2TestCase.Parameters.GasPrice, sbs2TestCase.Parameters.Data, passphrase, node, sbs2TestCase.Parameters.ConfirmationWaitTime)
	testCaseTx := testing.ConvertToTestCaseTransaction(fromAddress, fromAddress, rawTx, sbs2TestCase.Parameters, err)
	sbs2TestCase.Transactions = append(sbs2TestCase.Transactions, testCaseTx)

	testing.Log(sbs2TestCase.Name, fmt.Sprintf("Sent %f token(s) from %s to %s - transaction hash: %s, tx successful: %t", sbs2TestCase.Parameters.Amount, fromAddress, fromAddress, testCaseTx.TransactionHash, testCaseTx.Success), sbs2TestCase.Verbose)

	senderEndingBalance, _ := balances.GetShardBalance(fromAddress, sbs2TestCase.Parameters.FromShardID, node)

	testing.Log(sbs2TestCase.Name, fmt.Sprintf("Account %s (address: %s) has an ending balance of %f in shard %d after the test", keyName, fromAddress, senderEndingBalance, sbs2TestCase.Parameters.FromShardID), sbs2TestCase.Verbose)
	testing.Title(sbs2TestCase.Name, "footer", sbs2TestCase.Verbose)

	// We should end up with a lesser amount compared to the initial amount since we pay a gas fee
	sbs2TestCase.Result = (testCaseTx.Success && (senderStartingBalance == senderEndingBalance))

	return sbs2TestCase
}
