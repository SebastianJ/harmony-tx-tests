package testcases

import (
	"fmt"

	"github.com/SebastianJ/harmony-tx-tests/config"
	"github.com/SebastianJ/harmony-tx-tests/accounts"
	"github.com/SebastianJ/harmony-tx-tests/balances"
	"github.com/SebastianJ/harmony-tx-tests/testing"
	"github.com/SebastianJ/harmony-tx-tests/transactions"
)

// Sbs1TestCase - defines the common properties for the SBS1 test case
var Sbs1TestCase testing.TestCase = testing.TestCase{
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
		GasPrice:             1,
		Count:                1,
		ConfirmationWaitTime: 16,
	},
}

// RunSbs1TestCase - Same Beacon Shard single account transfer A -> B, Shard 0 -> 0, Amount 1, Tx Data nil, expects: successful token transferred from A to B within 2 blocks time 16s
func RunSbs1TestCase() {
	testing.Title(Sbs1TestCase.Name, "header", Sbs1TestCase.Verbose)
	testing.Log(Sbs1TestCase.Name, fmt.Sprintf("Using sender address: %s", config.Configuration.Funding.Account.Address), Sbs1TestCase.Verbose)

	sinkAccountName := fmt.Sprintf("%s_receiver", config.Configuration.Funding.Account.Address)
	testing.Log(Sbs1TestCase.Name, fmt.Sprintf("Generating a new receiver account: %s", sinkAccountName), Sbs1TestCase.Verbose)
	toAddress, err := accounts.GenerateAccountAndReturnAddress(sinkAccountName)

	senderStartingBalance, _ := balances.GetShardBalance(config.Configuration.Funding.Account.Address, Sbs1TestCase.Parameters.FromShardID)
	receiverStartingBalance, _ := balances.GetShardBalance(toAddress, Sbs1TestCase.Parameters.ToShardID)

	testing.Log(Sbs1TestCase.Name, fmt.Sprintf("Generated a new receiver account: %s, address: %s", sinkAccountName, toAddress), Sbs1TestCase.Verbose)
	testing.Log(Sbs1TestCase.Name, fmt.Sprintf("Using sender address: %s and receiver address : %s", config.Configuration.Funding.Account.Address, toAddress), Sbs1TestCase.Verbose)
	testing.Log(Sbs1TestCase.Name, fmt.Sprintf("Sender address: %s has a starting balance of %f in shard %d before the test", config.Configuration.Funding.Account.Address, senderStartingBalance, Sbs1TestCase.Parameters.FromShardID), Sbs1TestCase.Verbose)
	testing.Log(Sbs1TestCase.Name, fmt.Sprintf("Receiver account %s (address: %s) has a starting balance of %f in shard %d before the test", sinkAccountName, toAddress, receiverStartingBalance, Sbs1TestCase.Parameters.ToShardID), Sbs1TestCase.Verbose)
	testing.Log(Sbs1TestCase.Name, fmt.Sprintf("Will let the transaction wait up to %d seconds to try to get finalized within 2 blocks", Sbs1TestCase.Parameters.ConfirmationWaitTime), Sbs1TestCase.Verbose)
	testing.Log(Sbs1TestCase.Name, "Sending transaction...", Sbs1TestCase.Verbose)

	rawTx, err := transactions.SendTransaction(config.Configuration.Funding.Account.Address, Sbs1TestCase.Parameters.FromShardID, toAddress, Sbs1TestCase.Parameters.ToShardID, Sbs1TestCase.Parameters.Amount, Sbs1TestCase.Parameters.GasPrice, Sbs1TestCase.Parameters.Data, Sbs1TestCase.Parameters.ConfirmationWaitTime)
	testCaseTx := testing.ConvertToTestCaseTransaction(config.Configuration.Funding.Account.Address, toAddress, rawTx, Sbs1TestCase.Parameters, err)
	Sbs1TestCase.Transactions = append(Sbs1TestCase.Transactions, testCaseTx)

	testing.Log(Sbs1TestCase.Name, fmt.Sprintf("Sent %f token(s) from %s to %s - transaction hash: %s, tx successful: %t", Sbs1TestCase.Parameters.Amount, config.Configuration.Funding.Account.Address, toAddress, testCaseTx.TransactionHash, testCaseTx.Success), Sbs1TestCase.Verbose)

	senderEndingBalance, _ := balances.GetShardBalance(config.Configuration.Funding.Account.Address, Sbs1TestCase.Parameters.FromShardID)
	receiverEndingBalance, _ := balances.GetShardBalance(toAddress, Sbs1TestCase.Parameters.ToShardID)

	testing.Log(Sbs1TestCase.Name, fmt.Sprintf("Sender address: %s has an ending balance of %f in shard %d after the test", config.Configuration.Funding.Account.Address, senderEndingBalance, Sbs1TestCase.Parameters.FromShardID), Sbs1TestCase.Verbose)
	testing.Log(Sbs1TestCase.Name, fmt.Sprintf("Receiver address: %s has an ending balance of %f in shard %d after the test", toAddress, receiverEndingBalance, Sbs1TestCase.Parameters.ToShardID), Sbs1TestCase.Verbose)
	testing.Log(Sbs1TestCase.Name, "Performing test teardown (returning funds and removing sink account)", Sbs1TestCase.Verbose)
	testing.Title(Sbs1TestCase.Name, "footer", Sbs1TestCase.Verbose)

	testing.Teardown(sinkAccountName, toAddress, Sbs1TestCase.Parameters.FromShardID, config.Configuration.Funding.Account.Address, Sbs1TestCase.Parameters.ToShardID, Sbs1TestCase.Parameters.Amount, Sbs1TestCase.Parameters.GasPrice, 0)

	Sbs1TestCase.Result = (testCaseTx.Success && ((receiverStartingBalance)+Sbs1TestCase.Parameters.Amount == receiverEndingBalance))

	Results = append(Results, Sbs1TestCase)
}
