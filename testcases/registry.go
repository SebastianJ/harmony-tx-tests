package testcases

import (
	"github.com/SebastianJ/harmony-tx-tests/testing"
)

// RegistryItem - represents a given TestCase and its function to execute in order to perform a test
type RegistryItem struct {
	TestCase testing.TestCase
	Function interface{}
}

var (
	// Registry - contains all registered test cases
	Registry []RegistryItem

	// Results - contains all executed test case results
	Results []testing.TestCase
)

func init() {
	Registry = append(Registry, RegistryItem{
		TestCase: Sbs1TestCase,
		Function: RunSbs1TestCase,
	})

	Registry = append(Registry, RegistryItem{
		TestCase: Sbs2TestCase,
		Function: RunSbs2TestCase,
	})

	Registry = append(Registry, RegistryItem{
		TestCase: Sbs3TestCase,
		Function: RunSbs3TestCase,
	})

	Registry = append(Registry, RegistryItem{
		TestCase: Sbs4TestCase,
		Function: RunSbs4TestCase,
	})

	Registry = append(Registry, RegistryItem{
		TestCase: Sbs5TestCase,
		Function: RunSbs5TestCase,
	})

	Registry = append(Registry, RegistryItem{
		TestCase: Sbs6TestCase,
		Function: RunSbs6TestCase,
	})
}
