package testcases

import (
	"fmt"

	"github.com/SebastianJ/harmony-tx-tests/accounts"
	"github.com/SebastianJ/harmony-tx-tests/balances"
	"github.com/SebastianJ/harmony-tx-tests/transactions"
	"github.com/SebastianJ/harmony-tx-tests/utils"
)

var (
	testCaseName         = "SBS1"
	shardID              = uint32(0)
	txData               = ``
	amount               = 1.0
	gasPrice             = int64(1)
	confirmationWaitTime = 16
)

// Sbs1TestCase - Same Beacon Shard single account transfer A -> B, Shard 0 -> 0, Amount 1, Tx Data nil, expects: successful token transferred from A to B within 2 blocks time 16s
func Sbs1TestCase(accs map[string]string, passphrase string, node string) (bool, error) {
	var success bool

	keyName, fromAddress := utils.RandomItemFromMap(accs)

	TestLegend(testCaseName)
	TestLog(testCaseName, fmt.Sprintf("Using source/sender key: %s and address: %s", keyName, fromAddress))

	sinkAccountName := fmt.Sprintf("%s_sink", keyName)
	err := accounts.GenerateAccount(sinkAccountName, passphrase)
	toAddress := accounts.FindAccountAddressByName(sinkAccountName)

	senderStartingBalance, _ := balances.GetShardBalance(fromAddress, int(shardID), node)
	receiverStartingBalance, _ := balances.GetShardBalance(toAddress, int(shardID), node)

	TestLog(testCaseName, fmt.Sprintf("Generated a new receiver/sink account: %s, address: %s", sinkAccountName, toAddress))
	TestLog(testCaseName, fmt.Sprintf("Using source account %s (address: %s) and sink account %s (address : %s)", keyName, fromAddress, sinkAccountName, toAddress))
	TestLog(testCaseName, fmt.Sprintf("Source account %s (address: %s) has a starting balance of %f in shard %d before the test", keyName, fromAddress, senderStartingBalance, shardID))
	TestLog(testCaseName, fmt.Sprintf("Sink account %s (address: %s) has a starting balance of %f in shard %d before the test", sinkAccountName, toAddress, receiverStartingBalance, shardID))
	TestLog(testCaseName, fmt.Sprintf("Will let the transaction wait up to %d seconds to try to get finalized within 2 blocks", confirmationWaitTime))
	TestLog(testCaseName, "Sending transaction...")

	txResponse, err := transactions.SendSameShardTransaction(fromAddress, toAddress, shardID, amount, gasPrice, txData, passphrase, node, confirmationWaitTime)

	if err != nil {
		fmt.Println(fmt.Sprintf(`Error occurred: %s`, err))
		return false, err
	}

	txHash := txResponse["transactionHash"].(string)
	success = transactions.IsTransactionSuccessful(txResponse)

	TestLog(testCaseName, fmt.Sprintf("Sent %f token(s) from %s to %s - transaction hash: %s, tx successful: %t", amount, fromAddress, toAddress, txHash, success))

	senderEndingBalance, _ := balances.GetShardBalance(fromAddress, int(shardID), node)
	receiverEndingBalance, _ := balances.GetShardBalance(toAddress, int(shardID), node)

	TestLog(testCaseName, fmt.Sprintf("Source account %s (address: %s) has an ending balance of %f in shard %d after the test", keyName, fromAddress, senderEndingBalance, shardID))
	TestLog(testCaseName, fmt.Sprintf("Sink account %s (address: %s) has an ending balance of %f in shard %d after the test", sinkAccountName, toAddress, receiverEndingBalance, shardID))
	TestLog(testCaseName, "Performing test teardown (returning funds and removing sink account) ...")

	Teardown(sinkAccountName, toAddress, shardID, fromAddress, shardID, amount, gasPrice, passphrase, node, 0)
	TestLegend(testCaseName)

	success = (success && ((receiverStartingBalance)+amount == receiverEndingBalance))

	return success, nil
}
