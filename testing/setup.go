package testing

import (
	"fmt"
	"github.com/SebastianJ/harmony-tx-tests/utils"
	"strings"
)

var (
	// TestCases - contains all test cases that will get executed
	TestCases []TestCase

	// Results - contains all executed test case results
	Results []TestCase
)

func init() {
	testCases := loadTestCases()
	TestCases = testCases
}

func loadTestCases() (testCases []TestCase) {
	testCaseFiles, err := utils.GlobFiles("../testcases/*.yml")

	if err != nil {
		return nil
	}

	fmt.Println(fmt.Sprintf("Found a total of %d test case files", len(testCaseFiles)))

	for _, testCaseFile := range testCaseFiles {
		testCase := TestCase{}
		err := utils.ParseYaml(testCaseFile, &testCase)

		if err == nil {
			if testCase.TestType != "" {
				testCase.TestType = strings.ToLower(testCase.TestType)
			}

			testCases = append(testCases, testCase)
		}
	}

	return testCases
}
