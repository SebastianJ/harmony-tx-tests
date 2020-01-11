package testing

import "github.com/SebastianJ/harmony-tx-tests/transactions"

// ConvertToTestCaseTransaction - converts a raw tx response map to a typed TestCaseTransaction type
func ConvertToTestCaseTransaction(fromAddress string, toAddress string, rawTx map[string]interface{}, params TestCaseParameters, err error) TestCaseTransaction {
	if err != nil {
		return TestCaseTransaction{Error: err}
	}

	txHash := rawTx["transactionHash"].(string)
	success := transactions.IsTransactionSuccessful(rawTx)

	testCaseTransaction := TestCaseTransaction{
		FromAddress:     fromAddress,
		FromShardID:     params.FromShardID,
		ToAddress:       toAddress,
		ToShardID:       params.ToShardID,
		TransactionHash: txHash,
		Success:         success,
		Response:        rawTx,
	}

	return testCaseTransaction
}
