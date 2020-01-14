package transactions

import (
	"encoding/base64"

	"github.com/SebastianJ/harmony-tx-tests/config"
	sdkNonces "github.com/SebastianJ/harmony-sdk/nonces"
	sdkShards "github.com/SebastianJ/harmony-sdk/shards"
	sdkTxs "github.com/SebastianJ/harmony-sdk/transactions"
	"github.com/harmony-one/go-sdk/pkg/common"
	"github.com/harmony-one/go-sdk/pkg/rpc"
	"github.com/harmony-one/go-sdk/pkg/store"
)

// IsTransactionSuccessful - checks if a transaction is successful given a transaction response
func IsTransactionSuccessful(txResponse map[string]interface{}) (success bool) {
	txStatus, ok := txResponse["status"].(string)

	if txStatus != "" && ok {
		success = (txStatus == "0x1")
	}

	return success
}

// SendSameShardTransaction - send a transaction using the same shard for both the receiver and the sender
func SendSameShardTransaction(fromAddress string, toAddress string, shardID uint32, amount float64, nonce int, gasPrice int64, txData string, confirmationWaitTime int) (map[string]interface{}, error) {
	return SendTransaction(fromAddress, shardID, toAddress, shardID, amount, nonce, gasPrice, txData, confirmationWaitTime)
}

// NetworkHandler - resolve the RPC/HTTP Messenger to use for remote commands
func NetworkHandler(shardID uint32) (*rpc.HTTPMessenger, error) {
	node := config.GenerateNodeAddress(config.Configuration.Network.Name, shardID)
	networkHandler, err := sdkShards.HandlerForShard(shardID, node)
	if err != nil {
		return nil, err
	}

	return networkHandler, nil
}

// CurrentNonce - fetch the current nonce for a given address and RPC interface
func CurrentNonce(address string, networkHandler *rpc.HTTPMessenger) (uint64, error) {
	currentNonce, err := sdkNonces.GetNonceFromInput(address, "", networkHandler)

	if err != nil {
		return 0, err
	}

	return currentNonce, nil
}

// SendTransaction - send transactions
func SendTransaction(fromAddress string, fromShardID uint32, toAddress string, toShardID uint32, amount float64, nonce int, gasPrice int64, txData string, confirmationWaitTime int) (map[string]interface{}, error) {
	node := config.GenerateNodeAddress(config.Configuration.Network.Name, fromShardID)

	networkHandler, err := NetworkHandler(fromShardID)
	if err != nil {
		return nil, err
	}

	chain := &common.Chain.TestNet

	if config.Configuration.Network.Name == "localnet" {
		if confirmationWaitTime > 0 {
			confirmationWaitTime = confirmationWaitTime * 2
		}
	} else {
		chain, err = common.StringToChainID(config.Configuration.Network.Name)
	}

	var currentNonce uint64

	if nonce < 0 {
		currentNonce, err = CurrentNonce(fromAddress, networkHandler)
		if err != nil {
			return nil, err
		}
	} else {
		currentNonce = uint64(nonce)
	}

	if txData != "" {
		txData = base64.StdEncoding.EncodeToString([]byte(txData))
	}

	keystore, account, err := store.UnlockedKeystore(fromAddress, config.Configuration.Account.Passphrase)
	if err != nil {
		return nil, err
	}

	txResult, err := sdkTxs.SendTransaction(keystore, account, networkHandler, chain, fromAddress, fromShardID, toAddress, toShardID, amount, gasPrice, currentNonce, txData, config.Configuration.Account.Passphrase, node, confirmationWaitTime)

	if err != nil {
		return nil, err
	}

	return txResult, nil
}
