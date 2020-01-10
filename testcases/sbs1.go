package testcases

import (
	"fmt"
	"time"

	"github.com/SebastianJ/harmony-tx-tests/accounts"
	"github.com/SebastianJ/harmony-tx-tests/balances"
	"github.com/SebastianJ/harmony-tx-tests/transactions"
	"github.com/SebastianJ/harmony-tx-tests/utils"
)

// Sbs1TestCase - Same Beacon Shard single account transfer A -> B, Shard 0 -> 0, Amount 1, Tx Data nil, expects: successful token transferred from A to B within 2 blocks time 16s
func Sbs1TestCase(accs map[string]string, passphrase string, node string) (bool, error) {
	keyName, fromAddress := utils.RandomItemFromMap(accs)
	var success bool
	timeFormat := "2006-01-02 15:04:05"

	fmt.Println("\n-----Test case: SBS1---------------------------------------------------------------------------")
	fmt.Println(fmt.Sprintf("%s - [Test Case - SBS1]: Using source/sender key: %s and address: %s", time.Now().Format(timeFormat), keyName, fromAddress))

	shardID := uint32(0)
	txData := ``
	amount := 1.0
	gasPrice := int64(1)
	confirmationWaitTime := 16

	sinkAccountName := fmt.Sprintf("%s_sink", keyName)
	err := accounts.GenerateAccount(sinkAccountName, passphrase)
	toAddress := accounts.FindAccountAddressByName(sinkAccountName)

	senderStartingBalance, _ := balances.GetTotalBalance(fromAddress, node)
	receiverStartingBalance, _ := balances.GetTotalBalance(toAddress, node)

	fmt.Println(fmt.Sprintf("%s - [Test Case - SBS1]: Generated a new receiver/sink account: %s, address: %s", time.Now().Format(timeFormat), sinkAccountName, toAddress))
	fmt.Println(fmt.Sprintf("%s - [Test Case - SBS1]: Using source account %s (address: %s) and sink account %s (address : %s)", time.Now().Format(timeFormat), keyName, fromAddress, sinkAccountName, toAddress))
	fmt.Println(fmt.Sprintf("%s - [Test Case - SBS1]: Source account %s (address: %s) has a starting balance of %f across all shards before the test", time.Now().Format(timeFormat), keyName, fromAddress, senderStartingBalance))
	fmt.Println(fmt.Sprintf("%s - [Test Case - SBS1]: Sink account %s (address: %s) has a starting balance of %f across all shards before the test", time.Now().Format(timeFormat), sinkAccountName, toAddress, receiverStartingBalance))

	fmt.Println(fmt.Sprintf("%s - [Test Case - SBS1]: Will let the transaction wait up to %d seconds to try to get finalized within 2 blocks", time.Now().Format(timeFormat), confirmationWaitTime))
	fmt.Println(fmt.Sprintf("%s - [Test Case - SBS1]: Sending transaction...", time.Now().Format(timeFormat)))

	txResponse, err := transactions.SendSameShardTransaction(fromAddress, toAddress, shardID, amount, gasPrice, txData, passphrase, node, confirmationWaitTime)

	if err != nil {
		fmt.Println(fmt.Sprintf(`Error occurred: %s`, err))
		return false, err
	}

	txHash := txResponse["transactionHash"].(string)
	txStatus, ok := txResponse["status"].(string)

	if txStatus != "" && ok {
		success = (txStatus == "0x1")
	}

	fmt.Println(fmt.Sprintf(`%s - [Test Case - SBS1]: Sent %f token(s) from %s to %s - transaction hash: %s, tx status: %s`, time.Now().Format(timeFormat), amount, fromAddress, toAddress, txHash, txStatus))

	senderEndingBalance, _ := balances.GetTotalBalance(fromAddress, node)
	receiverEndingBalance, _ := balances.GetTotalBalance(toAddress, node)

	fmt.Println(fmt.Sprintf("%s - [Test Case - SBS1]: Source account %s (address: %s) has an ending balance of %f across all shards after the test", time.Now().Format(timeFormat), keyName, fromAddress, senderEndingBalance))
	fmt.Println(fmt.Sprintf("%s - [Test Case - SBS1]: Sink account %s (address: %s) has an ending balance of %f across all shards after the test", time.Now().Format(timeFormat), sinkAccountName, toAddress, receiverEndingBalance))

	fmt.Println("-----Test case: SBS1---------------------------------------------------------------------------\n")
	success = (success && ((receiverStartingBalance)+amount == receiverEndingBalance))

	return success, nil
}
