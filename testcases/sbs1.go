package testcases

import (
	"fmt"

	"github.com/SebastianJ/harmony-tx-tests/accounts"
	"github.com/SebastianJ/harmony-tx-tests/balances"
	"github.com/SebastianJ/harmony-tx-tests/transactions"
	"github.com/SebastianJ/harmony-tx-tests/utils"
)

var testCase TestCase = TestCase{
	Name:     "SBS1",
	Result:   false,
	Expected: true,
	Parameters: map[string]interface{}{
		"fromShardID":          0,
		"toShardID":            0,
		"txData":               "",
		"amount":               1.0,
		"gasPrice":             int64(1),
		"confirmationWaitTime": 16,
	},
}

// Sbs1TestCase - Same Beacon Shard single account transfer A -> B, Shard 0 -> 0, Amount 1, Tx Data nil, expects: successful token transferred from A to B within 2 blocks time 16s
func Sbs1TestCase(accs map[string]string, passphrase string, node string) TestCase {
	keyName, fromAddress := utils.RandomItemFromMap(accs)
	testCase.Parameters["fromAddress"] = fromAddress

	shardID := testCase.Parameters["fromShardID"].(int)
	amount := testCase.Parameters["amount"].(float64)
	gasPrice := testCase.Parameters["gasPrice"].(int64)
	txData := testCase.Parameters["txData"].(string)
	confirmationWaitTime := testCase.Parameters["confirmationWaitTime"].(int)

	TestTitle(testCase.Name, "header")
	TestLog(testCase.Name, fmt.Sprintf("Using source/sender key: %s and address: %s", keyName, fromAddress))

	sinkAccountName := fmt.Sprintf("%s_sink", keyName)
	TestLog(testCase.Name, fmt.Sprintf("Generating a new receiver/sink account: %s", sinkAccountName))

	err := accounts.GenerateAccount(sinkAccountName, passphrase)
	toAddress := accounts.FindAccountAddressByName(sinkAccountName)
	testCase.Parameters["toAddress"] = toAddress

	senderStartingBalance, _ := balances.GetShardBalance(fromAddress, shardID, node)
	receiverStartingBalance, _ := balances.GetShardBalance(toAddress, shardID, node)

	TestLog(testCase.Name, fmt.Sprintf("Generated a new receiver/sink account: %s, address: %s", sinkAccountName, toAddress))
	TestLog(testCase.Name, fmt.Sprintf("Using source account %s (address: %s) and sink account %s (address : %s)", keyName, fromAddress, sinkAccountName, toAddress))
	TestLog(testCase.Name, fmt.Sprintf("Source account %s (address: %s) has a starting balance of %f in shard %d before the test", keyName, fromAddress, senderStartingBalance, shardID))
	TestLog(testCase.Name, fmt.Sprintf("Sink account %s (address: %s) has a starting balance of %f in shard %d before the test", sinkAccountName, toAddress, receiverStartingBalance, shardID))
	TestLog(testCase.Name, fmt.Sprintf("Will let the transaction wait up to %d seconds to try to get finalized within 2 blocks", confirmationWaitTime))
	TestLog(testCase.Name, "Sending transaction...")

	testCase.TxData, err = transactions.SendSameShardTransaction(fromAddress, toAddress, uint32(shardID), amount, gasPrice, txData, passphrase, node, confirmationWaitTime)

	if err != nil {
		testCase.Error = err
		return testCase
	}

	txHash := testCase.TxData["transactionHash"].(string)
	success := transactions.IsTransactionSuccessful(testCase.TxData)

	TestLog(testCase.Name, fmt.Sprintf("Sent %f token(s) from %s to %s - transaction hash: %s, tx successful: %t", amount, fromAddress, toAddress, txHash, success))

	senderEndingBalance, _ := balances.GetShardBalance(fromAddress, shardID, node)
	receiverEndingBalance, _ := balances.GetShardBalance(toAddress, shardID, node)

	TestLog(testCase.Name, fmt.Sprintf("Source account %s (address: %s) has an ending balance of %f in shard %d after the test", keyName, fromAddress, senderEndingBalance, shardID))
	TestLog(testCase.Name, fmt.Sprintf("Sink account %s (address: %s) has an ending balance of %f in shard %d after the test", sinkAccountName, toAddress, receiverEndingBalance, shardID))
	TestLog(testCase.Name, "Performing test teardown (returning funds and removing sink account) ...")

	Teardown(sinkAccountName, toAddress, uint32(shardID), fromAddress, uint32(shardID), amount, gasPrice, passphrase, node, 0)
	TestTitle(testCase.Name, "footer")

	testCase.Result = (success && ((receiverStartingBalance)+amount == receiverEndingBalance))

	return testCase
}
