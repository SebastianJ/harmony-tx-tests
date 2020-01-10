package transactions

import (
	"github.com/SebastianJ/harmony-tx-sender/nonces"
	"github.com/SebastianJ/harmony-tx-sender/shards"
	senderTxs "github.com/SebastianJ/harmony-tx-sender/transactions"
	"github.com/harmony-one/go-sdk/pkg/common"
	"github.com/harmony-one/go-sdk/pkg/store"
)

type interfaceWrapper []interface{}

// SendSameShardTransaction - send a transaction using the same shard for both the receiver and the sender
func SendSameShardTransaction(fromAddress string, toAddress string, shardID uint32, amount float64, gasPrice int64, txData string, passphrase string, node string, confirmationWaitTime int) (map[string]interface{}, error) {
	return SendTransaction(fromAddress, shardID, toAddress, shardID, amount, gasPrice, txData, passphrase, node, confirmationWaitTime)
}

// SendTransaction - send transactions
func SendTransaction(fromAddress string, fromShardID uint32, toAddress string, toShardID uint32, amount float64, gasPrice int64, txData string, passphrase string, node string, confirmationWaitTime int) (map[string]interface{}, error) {
	networkHandler, err := shards.HandlerForShard(fromShardID, node)
	if err != nil {
		return nil, err
	}

	chain := &common.Chain.DevNet

	currentNonce, err := nonces.GetNonceFromInput(fromAddress, "", networkHandler)

	if err != nil {
		return nil, err
	}

	keystore, account, err := store.UnlockedKeystore(fromAddress, passphrase)
	if err != nil {
		return nil, err
	}

	txResult, err := senderTxs.SendTransaction(keystore, account, networkHandler, chain, fromAddress, fromShardID, toAddress, toShardID, amount, gasPrice, currentNonce, txData, passphrase, node, confirmationWaitTime)

	if err != nil {
		return nil, err
	}

	return txResult, nil
}
