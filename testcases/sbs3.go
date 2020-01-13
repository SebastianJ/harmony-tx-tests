package testcases

import (
	"fmt"

	"github.com/SebastianJ/harmony-tx-tests/accounts"
	"github.com/SebastianJ/harmony-tx-tests/balances"
	"github.com/SebastianJ/harmony-tx-tests/config"
	"github.com/SebastianJ/harmony-tx-tests/testing"
	"github.com/SebastianJ/harmony-tx-tests/transactions"
)

// Sbs3TestCase - defines the common properties for the SBS3 test case
var Sbs3TestCase testing.TestCase = testing.TestCase{
	Scenario: "Same Beacon Shard",
	Name:     "SBS3",
	Goal:     "Atto transfer",
	Priority: 0,
	Expected: true,
	Verbose:  true,
	Parameters: testing.TestCaseParameters{
		FromShardID:          0,
		ToShardID:            0,
		Data:                 "",
		Amount:               1.00E-18,
		Nonce:                -1, //negative nonce value = fetch the latest nonce from the network
		GasPrice:             1,
		Count:                1,
		ConfirmationWaitTime: 16,
	},
}

// RunSbs3TestCase - Same Beacon Shard single account transfer A -> B, Shard 0 -> 0, Amount 1.00E-18, Tx Data nil, expects: successful token transferred from A to B within 2 blocks time 16s
func RunSbs3TestCase() {
	testing.Title(Sbs3TestCase.Name, "header", Sbs3TestCase.Verbose)
	testing.Log(Sbs3TestCase.Name, fmt.Sprintf("Using sender account %s (address: %s)", config.Configuration.Funding.Account.Name, config.Configuration.Funding.Account.Address), Sbs3TestCase.Verbose)

	sinkAccountName := fmt.Sprintf("%s_sink", config.Configuration.Funding.Account.Name)
	testing.Log(Sbs3TestCase.Name, fmt.Sprintf("Generating a new receiver/sink account: %s", sinkAccountName), Sbs3TestCase.Verbose)
	toAddress, err := accounts.GenerateAccountAndReturnAddress(sinkAccountName)

	senderStartingBalance, _ := balances.GetShardBalance(config.Configuration.Funding.Account.Address, Sbs3TestCase.Parameters.FromShardID)
	receiverStartingBalance, _ := balances.GetShardBalance(toAddress, Sbs3TestCase.Parameters.ToShardID)

	testing.Log(Sbs3TestCase.Name, fmt.Sprintf("Generated a new receiver account: %s, address: %s", sinkAccountName, toAddress), Sbs3TestCase.Verbose)
	testing.Log(Sbs3TestCase.Name, fmt.Sprintf("Using sender account %s (address: %s) and sink account %s (address : %s)", config.Configuration.Funding.Account.Name, config.Configuration.Funding.Account.Address, sinkAccountName, toAddress), Sbs3TestCase.Verbose)
	testing.Log(Sbs3TestCase.Name, fmt.Sprintf("Sender account %s (address: %s) has a starting balance of %f in shard %d before the test", config.Configuration.Funding.Account.Name, config.Configuration.Funding.Account.Address, senderStartingBalance, Sbs3TestCase.Parameters.FromShardID), Sbs3TestCase.Verbose)
	testing.Log(Sbs3TestCase.Name, fmt.Sprintf("Receiver account %s (address: %s) has a starting balance of %f in shard %d before the test", sinkAccountName, toAddress, receiverStartingBalance, Sbs3TestCase.Parameters.ToShardID), Sbs3TestCase.Verbose)
	testing.Log(Sbs3TestCase.Name, fmt.Sprintf("Will let the transaction wait up to %d seconds to try to get finalized within 2 blocks", Sbs3TestCase.Parameters.ConfirmationWaitTime), Sbs3TestCase.Verbose)
	testing.Log(Sbs3TestCase.Name, "Sending transaction...", Sbs3TestCase.Verbose)

	rawTx, err := transactions.SendTransaction(config.Configuration.Funding.Account.Address, Sbs3TestCase.Parameters.FromShardID, toAddress, Sbs3TestCase.Parameters.ToShardID, Sbs3TestCase.Parameters.Amount, Sbs3TestCase.Parameters.Nonce, Sbs3TestCase.Parameters.GasPrice, Sbs3TestCase.Parameters.Data, Sbs3TestCase.Parameters.ConfirmationWaitTime)
	testCaseTx := testing.ConvertToTestCaseTransaction(config.Configuration.Funding.Account.Address, toAddress, rawTx, Sbs3TestCase.Parameters, err)
	Sbs3TestCase.Transactions = append(Sbs3TestCase.Transactions, testCaseTx)

	testing.Log(Sbs3TestCase.Name, fmt.Sprintf("Sent %f token(s) from %s to %s - transaction hash: %s, tx successful: %t", Sbs3TestCase.Parameters.Amount, config.Configuration.Funding.Account.Address, toAddress, testCaseTx.TransactionHash, testCaseTx.Success), Sbs3TestCase.Verbose)

	senderEndingBalance, _ := balances.GetShardBalance(config.Configuration.Funding.Account.Address, Sbs3TestCase.Parameters.FromShardID)
	receiverEndingBalance, _ := balances.GetShardBalance(toAddress, Sbs3TestCase.Parameters.ToShardID)

	testing.Log(Sbs3TestCase.Name, fmt.Sprintf("Source account %s (address: %s) has an ending balance of %f in shard %d after the test", config.Configuration.Funding.Account.Name, config.Configuration.Funding.Account.Address, senderEndingBalance, Sbs3TestCase.Parameters.FromShardID), Sbs3TestCase.Verbose)
	testing.Log(Sbs3TestCase.Name, fmt.Sprintf("Sink account %s (address: %s) has an ending balance of %f in shard %d after the test", sinkAccountName, toAddress, receiverEndingBalance, Sbs3TestCase.Parameters.ToShardID), Sbs3TestCase.Verbose)
	testing.Log(Sbs3TestCase.Name, "Performing test teardown (returning funds and removing sink account)", Sbs3TestCase.Verbose)
	testing.Title(Sbs3TestCase.Name, "footer", Sbs3TestCase.Verbose)

	testing.Teardown(sinkAccountName, toAddress, Sbs3TestCase.Parameters.FromShardID, config.Configuration.Funding.Account.Address, Sbs3TestCase.Parameters.ToShardID, Sbs3TestCase.Parameters.Amount, Sbs3TestCase.Parameters.GasPrice, 0)

	//TODO: tests with super small denominations can't rely on the balances supplied from balances.GetShardBalance since that function only returns 6 decimals
	Sbs3TestCase.Result = testCaseTx.Success

	Results = append(Results, Sbs3TestCase)
}
