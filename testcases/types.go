package testcases

// TestCase - represents a test case as a type
type TestCase struct {
	Name       string
	Result     bool
	Expected   bool
	Error      error
	Parameters map[string]interface{}
	TxData     map[string]interface{}
}
