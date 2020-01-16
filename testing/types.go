package testing

import (
	"github.com/SebastianJ/harmony-tx-tests/accounts"
)

// TestCase - represents a test case
type TestCase struct {
	Name         string `yaml:"name"`
	Scenario     string `yaml:"scenario"`
	Goal         string `yaml:"goal"`
	Priority     int    `yaml:"priority"`
	Execute      bool   `yaml:"execute"`
	Result       bool   `yaml:"result"`
	Expected     bool   `yaml:"expected"`
	Verbose      bool   `yaml:"verbose"`
	TestType     string `yaml:"test_type"`
	Error        error
	Parameters   TestCaseParameters `yaml:"parameters"`
	Transactions []TestCaseTransaction
	Function     interface{}
}

// TestCaseParameters - represents the test case parameters
type TestCaseParameters struct {
	SenderCount          int `yaml:"sender_count"`
	Senders              []accounts.Account
	ReceiverCount        int `yaml:"receiver_count"`
	Receivers            []accounts.Account
	FromShardID          uint32  `yaml:"from_shard_id"`
	ToShardID            uint32  `yaml:"to_shard_id"`
	Data                 string  `yaml:"data"`
	DataSize             int     `yaml:"data_size,omitempty"`
	Amount               float64 `yaml:"amount"`
	GasPrice             int64   `yaml:"gas_price"`
	Nonce                int     `yaml:"nonce"`
	Count                int     `yaml:"count"`
	ConfirmationWaitTime int     `yaml:"confirmation_wait_time"`
}

// TestCaseTransaction - represents an executed test case transaction
type TestCaseTransaction struct {
	FromAddress          string
	FromShardID          uint32
	ToAddress            string
	ToShardID            uint32
	Data                 string
	Amount               float64
	GasPrice             int64
	ConfirmationWaitTime int
	TransactionHash      string
	Success              bool
	Response             map[string]interface{}
	Error                error
}
