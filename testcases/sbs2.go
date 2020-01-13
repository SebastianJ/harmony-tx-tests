package testcases

import (
	"fmt"

	"github.com/SebastianJ/harmony-tx-tests/balances"
	"github.com/SebastianJ/harmony-tx-tests/config"
	"github.com/SebastianJ/harmony-tx-tests/testing"
	"github.com/SebastianJ/harmony-tx-tests/transactions"
)

// Sbs2TestCase - defines the common properties for the SBS2 test case
var Sbs2TestCase testing.TestCase = testing.TestCase{
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
		Nonce:                -1, //negative nonce value = fetch the latest nonce from the network
		GasPrice:             1,
		Count:                1,
		ConfirmationWaitTime: 16,
	},
}

// RunSbs2TestCase - Same Beacon Shard same account transfer A -> A, Shard 0 -> 0, Amount 1, Tx Data nil, expects: successful token transferred from A to A within 2 blocks time 16s
func RunSbs2TestCase() {
	testing.Title(Sbs2TestCase.Name, "header", Sbs2TestCase.Verbose)

	senderStartingBalance, _ := balances.GetShardBalance(config.Configuration.Funding.Account.Address, Sbs2TestCase.Parameters.FromShardID)

	testing.Log(Sbs2TestCase.Name, fmt.Sprintf("Using account %s (address: %s) for a self transfer", config.Configuration.Funding.Account.Name, config.Configuration.Funding.Account.Address), Sbs2TestCase.Verbose)
	testing.Log(Sbs2TestCase.Name, fmt.Sprintf("Source account %s (address: %s) has a starting balance of %f in shard %d before the test", config.Configuration.Funding.Account.Name, config.Configuration.Funding.Account.Address, senderStartingBalance, Sbs2TestCase.Parameters.FromShardID), Sbs2TestCase.Verbose)
	testing.Log(Sbs2TestCase.Name, fmt.Sprintf("Will let the transaction wait up to %d seconds to try to get finalized within 2 blocks", Sbs2TestCase.Parameters.ConfirmationWaitTime), Sbs2TestCase.Verbose)
	testing.Log(Sbs2TestCase.Name, "Sending transaction...", Sbs2TestCase.Verbose)

	rawTx, err := transactions.SendTransaction(config.Configuration.Funding.Account.Address, Sbs2TestCase.Parameters.FromShardID, config.Configuration.Funding.Account.Address, Sbs2TestCase.Parameters.ToShardID, Sbs2TestCase.Parameters.Amount, Sbs2TestCase.Parameters.Nonce, Sbs2TestCase.Parameters.GasPrice, Sbs2TestCase.Parameters.Data, Sbs2TestCase.Parameters.ConfirmationWaitTime)
	testCaseTx := testing.ConvertToTestCaseTransaction(config.Configuration.Funding.Account.Address, config.Configuration.Funding.Account.Address, rawTx, Sbs2TestCase.Parameters, err)
	Sbs2TestCase.Transactions = append(Sbs2TestCase.Transactions, testCaseTx)

	testing.Log(Sbs2TestCase.Name, fmt.Sprintf("Sent %f token(s) from %s to %s - transaction hash: %s, tx successful: %t", Sbs2TestCase.Parameters.Amount, config.Configuration.Funding.Account.Address, config.Configuration.Funding.Account.Address, testCaseTx.TransactionHash, testCaseTx.Success), Sbs2TestCase.Verbose)

	senderEndingBalance, _ := balances.GetShardBalance(config.Configuration.Funding.Account.Address, Sbs2TestCase.Parameters.FromShardID)

	testing.Log(Sbs2TestCase.Name, fmt.Sprintf("Account %s (address: %s) has an ending balance of %f in shard %d after the test", config.Configuration.Funding.Account.Name, config.Configuration.Funding.Account.Address, senderEndingBalance, Sbs2TestCase.Parameters.FromShardID), Sbs2TestCase.Verbose)
	testing.Title(Sbs2TestCase.Name, "footer", Sbs2TestCase.Verbose)

	// We should end up with a lesser amount compared to the initial amount since we pay a gas fee
	Sbs2TestCase.Result = (testCaseTx.Success && (senderStartingBalance <= senderEndingBalance))

	Results = append(Results, Sbs2TestCase)
}
