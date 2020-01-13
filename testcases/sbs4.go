package testcases

import (
	"fmt"

	"github.com/SebastianJ/harmony-tx-tests/accounts"
	"github.com/SebastianJ/harmony-tx-tests/balances"
	"github.com/SebastianJ/harmony-tx-tests/config"
	"github.com/SebastianJ/harmony-tx-tests/testing"
	"github.com/SebastianJ/harmony-tx-tests/transactions"
)

// Sbs4TestCase - defines the common properties for the SBS4 test case
var Sbs4TestCase testing.TestCase = testing.TestCase{
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
		Nonce:                -1, //negative nonce value = fetch the latest nonce from the network
		GasPrice:             1,
		Count:                1,
		ConfirmationWaitTime: 16,
	},
}

// RunSbs4TestCase - Same Beacon Shard single account transfer A -> B, Shard 0 -> 0, Amount 1.00E-09, Tx Data nil, expects: successful token transferred from A to B within 2 blocks time 16s
func RunSbs4TestCase() {
	testing.Title(Sbs4TestCase.Name, "header", Sbs4TestCase.Verbose)
	testing.Log(Sbs4TestCase.Name, fmt.Sprintf("Using source/sender key: %s and address: %s", config.Configuration.Funding.Account.Name, config.Configuration.Funding.Account.Address), Sbs4TestCase.Verbose)

	sinkAccountName := fmt.Sprintf("%s_sink", config.Configuration.Funding.Account.Name)
	testing.Log(Sbs4TestCase.Name, fmt.Sprintf("Generating a new receiver/sink account: %s", sinkAccountName), Sbs4TestCase.Verbose)
	toAddress, err := accounts.GenerateAccountAndReturnAddress(sinkAccountName)

	senderStartingBalance, _ := balances.GetShardBalance(config.Configuration.Funding.Account.Address, Sbs4TestCase.Parameters.FromShardID)
	receiverStartingBalance, _ := balances.GetShardBalance(toAddress, Sbs4TestCase.Parameters.ToShardID)

	testing.Log(Sbs4TestCase.Name, fmt.Sprintf("Generated a new receiver/sink account: %s, address: %s", sinkAccountName, toAddress), Sbs4TestCase.Verbose)
	testing.Log(Sbs4TestCase.Name, fmt.Sprintf("Using source account %s (address: %s) and sink account %s (address : %s)", config.Configuration.Funding.Account.Name, config.Configuration.Funding.Account.Address, sinkAccountName, toAddress), Sbs4TestCase.Verbose)
	testing.Log(Sbs4TestCase.Name, fmt.Sprintf("Source account %s (address: %s) has a starting balance of %f in shard %d before the test", config.Configuration.Funding.Account.Name, config.Configuration.Funding.Account.Address, senderStartingBalance, Sbs4TestCase.Parameters.FromShardID), Sbs4TestCase.Verbose)
	testing.Log(Sbs4TestCase.Name, fmt.Sprintf("Sink account %s (address: %s) has a starting balance of %f in shard %d before the test", sinkAccountName, toAddress, receiverStartingBalance, Sbs4TestCase.Parameters.ToShardID), Sbs4TestCase.Verbose)
	testing.Log(Sbs4TestCase.Name, fmt.Sprintf("Will let the transaction wait up to %d seconds to try to get finalized within 2 blocks", Sbs4TestCase.Parameters.ConfirmationWaitTime), Sbs4TestCase.Verbose)
	testing.Log(Sbs4TestCase.Name, "Sending transaction...", Sbs4TestCase.Verbose)

	rawTx, err := transactions.SendTransaction(config.Configuration.Funding.Account.Address, Sbs4TestCase.Parameters.FromShardID, toAddress, Sbs4TestCase.Parameters.ToShardID, Sbs4TestCase.Parameters.Amount, Sbs4TestCase.Parameters.Nonce, Sbs4TestCase.Parameters.GasPrice, Sbs4TestCase.Parameters.Data, Sbs4TestCase.Parameters.ConfirmationWaitTime)
	testCaseTx := testing.ConvertToTestCaseTransaction(config.Configuration.Funding.Account.Address, toAddress, rawTx, Sbs4TestCase.Parameters, err)
	Sbs4TestCase.Transactions = append(Sbs4TestCase.Transactions, testCaseTx)

	testing.Log(Sbs4TestCase.Name, fmt.Sprintf("Sent %f token(s) from %s to %s - transaction hash: %s, tx successful: %t", Sbs4TestCase.Parameters.Amount, config.Configuration.Funding.Account.Address, toAddress, testCaseTx.TransactionHash, testCaseTx.Success), Sbs4TestCase.Verbose)

	senderEndingBalance, _ := balances.GetShardBalance(config.Configuration.Funding.Account.Address, Sbs4TestCase.Parameters.FromShardID)
	receiverEndingBalance, _ := balances.GetShardBalance(toAddress, Sbs4TestCase.Parameters.ToShardID)

	testing.Log(Sbs4TestCase.Name, fmt.Sprintf("Source account %s (address: %s) has an ending balance of %f in shard %d after the test", config.Configuration.Funding.Account.Name, config.Configuration.Funding.Account.Address, senderEndingBalance, Sbs4TestCase.Parameters.FromShardID), Sbs4TestCase.Verbose)
	testing.Log(Sbs4TestCase.Name, fmt.Sprintf("Sink account %s (address: %s) has an ending balance of %f in shard %d after the test", sinkAccountName, toAddress, receiverEndingBalance, Sbs4TestCase.Parameters.ToShardID), Sbs4TestCase.Verbose)
	testing.Log(Sbs4TestCase.Name, "Performing test teardown (returning funds and removing sink account)", Sbs4TestCase.Verbose)
	testing.Title(Sbs4TestCase.Name, "footer", Sbs4TestCase.Verbose)

	testing.Teardown(sinkAccountName, toAddress, Sbs4TestCase.Parameters.FromShardID, config.Configuration.Funding.Account.Address, Sbs4TestCase.Parameters.ToShardID, Sbs4TestCase.Parameters.Amount, Sbs4TestCase.Parameters.GasPrice, 0)

	//TODO: tests with super small denominations can't rely on the balances supplied from balances.GetShardBalance since that function only returns 6 decimals
	Sbs4TestCase.Result = testCaseTx.Success

	Results = append(Results, Sbs4TestCase)
}
