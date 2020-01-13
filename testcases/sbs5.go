package testcases

import (
	"fmt"

	"github.com/SebastianJ/harmony-tx-tests/accounts"
	"github.com/SebastianJ/harmony-tx-tests/balances"
	"github.com/SebastianJ/harmony-tx-tests/config"
	"github.com/SebastianJ/harmony-tx-tests/testing"
	"github.com/SebastianJ/harmony-tx-tests/transactions"
)

// Sbs5TestCase - defines the common properties for the SBS5 test case
var Sbs5TestCase testing.TestCase = testing.TestCase{
	Scenario: "Same Beacon Shard",
	Name:     "SBS5",
	Goal:     "Insufficient amount",
	Priority: 0,
	Expected: true,
	Verbose:  true,
	Parameters: testing.TestCaseParameters{
		FromShardID:          0,
		ToShardID:            0,
		Data:                 "",
		Amount:               1.00E+20,
		Nonce:                -1, //negative nonce value = fetch the latest nonce from the network
		GasPrice:             1,
		Count:                1,
		ConfirmationWaitTime: 16,
	},
}

// RunSbs5TestCase - Same Beacon Shard single account transfer A -> B, Shard 0 -> 0, Amount 1.00E+20, Tx Data nil, expects: unsuccessful token transfer from A to B within 2 blocks time 16s
func RunSbs5TestCase() {
	testing.Title(Sbs5TestCase.Name, "header", Sbs5TestCase.Verbose)
	testing.Log(Sbs5TestCase.Name, fmt.Sprintf("Using source/sender key: %s and address: %s", config.Configuration.Funding.Account.Name, config.Configuration.Funding.Account.Address), Sbs5TestCase.Verbose)

	sinkAccountName := fmt.Sprintf("%s_sink", config.Configuration.Funding.Account.Name)
	testing.Log(Sbs5TestCase.Name, fmt.Sprintf("Generating a new receiver/sink account: %s", sinkAccountName), Sbs5TestCase.Verbose)
	toAddress, err := accounts.GenerateAccountAndReturnAddress(sinkAccountName)

	senderStartingBalance, _ := balances.GetShardBalance(config.Configuration.Funding.Account.Address, Sbs5TestCase.Parameters.FromShardID)
	receiverStartingBalance, _ := balances.GetShardBalance(toAddress, Sbs5TestCase.Parameters.ToShardID)

	testing.Log(Sbs5TestCase.Name, fmt.Sprintf("Generated a new receiver/sink account: %s, address: %s", sinkAccountName, toAddress), Sbs5TestCase.Verbose)
	testing.Log(Sbs5TestCase.Name, fmt.Sprintf("Using source account %s (address: %s) and sink account %s (address : %s)", config.Configuration.Funding.Account.Name, config.Configuration.Funding.Account.Address, sinkAccountName, toAddress), Sbs5TestCase.Verbose)
	testing.Log(Sbs5TestCase.Name, fmt.Sprintf("Source account %s (address: %s) has a starting balance of %f in shard %d before the test", config.Configuration.Funding.Account.Name, config.Configuration.Funding.Account.Address, senderStartingBalance, Sbs5TestCase.Parameters.FromShardID), Sbs5TestCase.Verbose)
	testing.Log(Sbs5TestCase.Name, fmt.Sprintf("Sink account %s (address: %s) has a starting balance of %f in shard %d before the test", sinkAccountName, toAddress, receiverStartingBalance, Sbs5TestCase.Parameters.ToShardID), Sbs5TestCase.Verbose)
	testing.Log(Sbs5TestCase.Name, fmt.Sprintf("Will let the transaction wait up to %d seconds to try to get finalized within 2 blocks", Sbs5TestCase.Parameters.ConfirmationWaitTime), Sbs5TestCase.Verbose)
	testing.Log(Sbs5TestCase.Name, "Sending transaction...", Sbs5TestCase.Verbose)

	rawTx, err := transactions.SendTransaction(config.Configuration.Funding.Account.Address, Sbs5TestCase.Parameters.FromShardID, toAddress, Sbs5TestCase.Parameters.ToShardID, Sbs5TestCase.Parameters.Amount, Sbs5TestCase.Parameters.Nonce, Sbs5TestCase.Parameters.GasPrice, Sbs5TestCase.Parameters.Data, Sbs5TestCase.Parameters.ConfirmationWaitTime)
	testCaseTx := testing.ConvertToTestCaseTransaction(config.Configuration.Funding.Account.Address, toAddress, rawTx, Sbs5TestCase.Parameters, err)
	Sbs5TestCase.Transactions = append(Sbs5TestCase.Transactions, testCaseTx)

	testing.Log(Sbs5TestCase.Name, fmt.Sprintf("Sent %f token(s) from %s to %s - transaction hash: %s, tx successful: %t", Sbs5TestCase.Parameters.Amount, config.Configuration.Funding.Account.Address, toAddress, testCaseTx.TransactionHash, testCaseTx.Success), Sbs5TestCase.Verbose)

	senderEndingBalance, _ := balances.GetShardBalance(config.Configuration.Funding.Account.Address, Sbs5TestCase.Parameters.FromShardID)
	receiverEndingBalance, _ := balances.GetShardBalance(toAddress, Sbs5TestCase.Parameters.ToShardID)

	testing.Log(Sbs5TestCase.Name, fmt.Sprintf("Source account %s (address: %s) has an ending balance of %f in shard %d after the test", config.Configuration.Funding.Account.Name, config.Configuration.Funding.Account.Address, senderEndingBalance, Sbs5TestCase.Parameters.FromShardID), Sbs5TestCase.Verbose)
	testing.Log(Sbs5TestCase.Name, fmt.Sprintf("Sink account %s (address: %s) has an ending balance of %f in shard %d after the test", sinkAccountName, toAddress, receiverEndingBalance, Sbs5TestCase.Parameters.ToShardID), Sbs5TestCase.Verbose)
	testing.Log(Sbs5TestCase.Name, "Performing test teardown (returning funds and removing sink account)", Sbs5TestCase.Verbose)
	testing.Title(Sbs5TestCase.Name, "footer", Sbs5TestCase.Verbose)

	testing.Teardown(sinkAccountName, toAddress, Sbs5TestCase.Parameters.FromShardID, config.Configuration.Funding.Account.Address, Sbs5TestCase.Parameters.ToShardID, Sbs5TestCase.Parameters.Amount, Sbs5TestCase.Parameters.GasPrice, 0)

	Sbs5TestCase.Result = (!testCaseTx.Success && receiverEndingBalance == 0.0)

	Results = append(Results, Sbs5TestCase)
}
