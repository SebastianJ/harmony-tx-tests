package testcases

// TestCase - represents a test case as a type
type TestCase struct {
	Name       	string
	Scenario   	string
	Goal 	   	string
	Priority   	int
	Result     	bool
	Expected   	bool
	Error      	error
	Parameters 	map[string]interface{}
	Transaction map[string]interface{}
}
