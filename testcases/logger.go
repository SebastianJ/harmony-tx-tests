package testcases

import (
	"fmt"
	"time"
)

var (
	timeFormat = "2006-01-02 15:04:05"
)

// TestLegend - header/footer for test cases
func TestLegend(name string) {
	fmt.Println(fmt.Sprintf("\n-----Test case: %s---------------------------------------------------------------------------", name))
}

// TestLog - time stamped logging messages for test cases
func TestLog(name string, message string) {
	fmt.Println(fmt.Sprintf("%s - [Test Case - %s]: %s", time.Now().Format(timeFormat), name, message))
}
