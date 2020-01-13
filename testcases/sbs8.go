package testcases

import (
	"fmt"

	"github.com/SebastianJ/harmony-tx-tests/accounts"
	"github.com/SebastianJ/harmony-tx-tests/balances"
	"github.com/SebastianJ/harmony-tx-tests/config"
	"github.com/SebastianJ/harmony-tx-tests/testing"
	"github.com/SebastianJ/harmony-tx-tests/transactions"
)

// Sbs8TestCase - defines the common properties for the SBS8 test case
var Sbs8TestCase testing.TestCase = testing.TestCase{
	Scenario: "Same Beacon Shard",
	Name:     "SBS8",
	Goal:     "One byte tx data",
	Priority: 0,
	Expected: true,
	Verbose:  true,
	Parameters: testing.TestCaseParameters{
		FromShardID:          0,
		ToShardID:            0,
		Data:                 "a",
		Amount:               1.0,
		Nonce:                -1, //negative nonce value = fetch the latest nonce from the network
		GasPrice:             1,
		Count:                1,
		ConfirmationWaitTime: 16,
	},
}

// RunSbs8TestCase - Same Beacon Shard one byte tx data transfer A -> B, Shard 0 -> 0, Amount 1, Tx Data 1 byte, expects: successful token transferred from A to B within 2 blocks time 16s
func RunSbs8TestCase() {
	testing.Title(Sbs8TestCase.Name, "header", Sbs8TestCase.Verbose)
	testing.Log(Sbs8TestCase.Name, fmt.Sprintf("Using sender address: %s", config.Configuration.Funding.Account.Address), Sbs8TestCase.Verbose)

	sinkAccountName := fmt.Sprintf("%s_receiver", config.Configuration.Funding.Account.Address)
	testing.Log(Sbs8TestCase.Name, fmt.Sprintf("Generating a new receiver account: %s", sinkAccountName), Sbs8TestCase.Verbose)
	toAddress, err := accounts.GenerateAccountAndReturnAddress(sinkAccountName)

	senderStartingBalance, _ := balances.GetShardBalance(config.Configuration.Funding.Account.Address, Sbs8TestCase.Parameters.FromShardID)
	receiverStartingBalance, _ := balances.GetShardBalance(toAddress, Sbs8TestCase.Parameters.ToShardID)

	testing.Log(Sbs8TestCase.Name, fmt.Sprintf("Generated a new receiver account: %s, address: %s", sinkAccountName, toAddress), Sbs8TestCase.Verbose)
	testing.Log(Sbs8TestCase.Name, fmt.Sprintf("Using sender address: %s and receiver address : %s", config.Configuration.Funding.Account.Address, toAddress), Sbs8TestCase.Verbose)
	testing.Log(Sbs8TestCase.Name, fmt.Sprintf("Sender address: %s has a starting balance of %f in shard %d before the test", config.Configuration.Funding.Account.Address, senderStartingBalance, Sbs8TestCase.Parameters.FromShardID), Sbs8TestCase.Verbose)
	testing.Log(Sbs8TestCase.Name, fmt.Sprintf("Receiver account %s (address: %s) has a starting balance of %f in shard %d before the test", sinkAccountName, toAddress, receiverStartingBalance, Sbs8TestCase.Parameters.ToShardID), Sbs8TestCase.Verbose)
	testing.Log(Sbs8TestCase.Name, fmt.Sprintf("Will let the transaction wait up to %d seconds to try to get finalized within 2 blocks", Sbs8TestCase.Parameters.ConfirmationWaitTime), Sbs8TestCase.Verbose)
	testing.Log(Sbs8TestCase.Name, fmt.Sprintf("Sending transaction using tx data %s: %d byte(s)", Sbs8TestCase.Parameters.Data, len([]byte(Sbs8TestCase.Parameters.Data))), Sbs8TestCase.Verbose)

	rawTx, err := transactions.SendTransaction(config.Configuration.Funding.Account.Address, Sbs8TestCase.Parameters.FromShardID, toAddress, Sbs8TestCase.Parameters.ToShardID, Sbs8TestCase.Parameters.Amount, Sbs8TestCase.Parameters.Nonce, Sbs8TestCase.Parameters.GasPrice, Sbs8TestCase.Parameters.Data, Sbs8TestCase.Parameters.ConfirmationWaitTime)
	testCaseTx := testing.ConvertToTestCaseTransaction(config.Configuration.Funding.Account.Address, toAddress, rawTx, Sbs8TestCase.Parameters, err)
	Sbs8TestCase.Transactions = append(Sbs8TestCase.Transactions, testCaseTx)

	testing.Log(Sbs8TestCase.Name, fmt.Sprintf("Sent %f token(s) from %s to %s - transaction hash: %s, tx successful: %t", Sbs8TestCase.Parameters.Amount, config.Configuration.Funding.Account.Address, toAddress, testCaseTx.TransactionHash, testCaseTx.Success), Sbs8TestCase.Verbose)

	senderEndingBalance, _ := balances.GetShardBalance(config.Configuration.Funding.Account.Address, Sbs8TestCase.Parameters.FromShardID)
	receiverEndingBalance, _ := balances.GetShardBalance(toAddress, Sbs8TestCase.Parameters.ToShardID)

	testing.Log(Sbs8TestCase.Name, fmt.Sprintf("Sender address: %s has an ending balance of %f in shard %d after the test", config.Configuration.Funding.Account.Address, senderEndingBalance, Sbs8TestCase.Parameters.FromShardID), Sbs8TestCase.Verbose)
	testing.Log(Sbs8TestCase.Name, fmt.Sprintf("Receiver address: %s has an ending balance of %f in shard %d after the test", toAddress, receiverEndingBalance, Sbs8TestCase.Parameters.ToShardID), Sbs8TestCase.Verbose)
	testing.Log(Sbs8TestCase.Name, "Performing test teardown (returning funds and removing sink account)", Sbs8TestCase.Verbose)
	testing.Title(Sbs8TestCase.Name, "footer", Sbs8TestCase.Verbose)

	testing.Teardown(sinkAccountName, toAddress, Sbs8TestCase.Parameters.FromShardID, config.Configuration.Funding.Account.Address, Sbs8TestCase.Parameters.ToShardID, Sbs8TestCase.Parameters.Amount, Sbs8TestCase.Parameters.GasPrice, 0)

	Sbs8TestCase.Result = (testCaseTx.Success && ((receiverStartingBalance)+Sbs8TestCase.Parameters.Amount == receiverEndingBalance))

	Results = append(Results, Sbs8TestCase)
}
