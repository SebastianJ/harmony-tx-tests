package testing

import (
	"fmt"

	"github.com/SebastianJ/harmony-tx-tests/accounts"
	"github.com/SebastianJ/harmony-tx-tests/funding"
)

// ExecuteTestCases - executes all registered/identified test cases
func ExecuteTestCases() error {
	accs, err := accounts.LoadSourceAccounts()

	if err != nil {
		return err
	}

	if err = funding.SetupFundingAccount(accs); err != nil {
		return err
	}

	execute()
	results()

	return nil
}

func execute() {
	for _, testCase := range TestCases {
		if testCase.Execute {
			executed := false

			switch testCase.TestType {
			case "standard":
				testCase = RunStandardTestCase(testCase)
				executed = true
			case "same_account":
				testCase = RunSameAccountTestCase(testCase)
				executed = true
			case "multiple_senders":
				testCase = RunMultipleSenderTestCase(testCase)
				executed = true
			case "multiple_receivers_invalid_nonce":
				testCase = RunMultipleReceiverInvalidNonceTestCase(testCase)
				executed = true
			default:
				fmt.Println(fmt.Sprintf("Please specify a valid test type for your test case %s", testCase.Name))
			}

			if executed {
				Results = append(Results, testCase)
			}

			//registryItem.Function.(func(TestCase))(registryItem.TestCase)
		} else {
			fmt.Println(fmt.Sprintf("\nTest case %s has the execute attribute set to false - make sure to set it to true if you want to execute this test case\n", testCase.Name))
		}
	}
}

func fund() {
	for _, testCase := range TestCases {
		if testCase.Execute {

		}
	}
}

func results() {
	successfulCount := 0
	failedCount := 0

	for _, testCase := range Results {
		if testCase.Result == testCase.Expected {
			successfulCount++
		} else {
			failedCount++
		}
	}

	fmt.Println("------------------------------------------------------------")
	fmt.Println(fmt.Sprintf("Test suite status - executed a total of %d test case(s):", len(Results)))
	fmt.Println(fmt.Sprintf("Successful: %d", successfulCount))
	fmt.Println(fmt.Sprintf("Failed: %d", failedCount))
	fmt.Println("------------------------------------------------------------\n")

	if len(Results) > 0 {
		fmt.Println("Executed test cases:")
		fmt.Println("------------------------------------------------------------")
		for _, testCase := range Results {
			if testCase.Result == testCase.Expected {
				fmt.Println(fmt.Sprintf("Testcase %s: success", testCase.Name))
			} else {
				fmt.Println(fmt.Sprintf("Testcase %s: failed", testCase.Name))
			}
		}
		fmt.Println("------------------------------------------------------------\n")
	}
}
