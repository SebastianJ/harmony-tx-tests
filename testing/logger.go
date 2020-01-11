package testing

import (
	"fmt"
	"time"
)

var (
	timeFormat = "2006-01-02 15:04:05"
)

// TestTitle - header/footer for test cases
func Title(name string, titleType string, verbose bool) {
	if verbose {
		if titleType == "header" {
			fmt.Println("\n")
		}

		fmt.Println(fmt.Sprintf("-----Test case: %s---------------------------------------------------------------------------", name))

		if titleType == "footer" {
			fmt.Println("\n")
		}
	}
}

// TestLog - time stamped logging messages for test cases
func Log(name string, message string, verbose bool) {
	if verbose {
		fmt.Println(fmt.Sprintf("%s - [Test Case - %s]: %s", time.Now().Format(timeFormat), name, message))
	}
}
