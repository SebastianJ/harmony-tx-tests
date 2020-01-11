package testing

// TestCase - represents a test case
type TestCase struct {
	Name         string
	Scenario     string
	Goal         string
	Priority     int
	Result       bool
	Expected     bool
	Error        error
	Verbose      bool
	Parameters   TestCaseParameters
	Transactions []TestCaseTransaction
}

// TestCaseParameters - represents the test case parameters
type TestCaseParameters struct {
	FromShardID          uint32
	ToShardID            uint32
	Data                 string
	Amount               float64
	GasPrice             int64
	Count                int
	ConfirmationWaitTime int
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
