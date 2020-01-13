package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/SebastianJ/harmony-tx-tests/config"
	"github.com/SebastianJ/harmony-tx-tests/testcases"
	"github.com/SebastianJ/harmony-tx-tests/funding"
	"github.com/SebastianJ/harmony-tx-tests/accounts"

	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "Harmony tx tests"
	app.Version = fmt.Sprintf("%s/%s-%s", runtime.Version(), runtime.GOOS, runtime.GOARCH)
	app.Usage = "Runs a set of Harmony tx tests"

	app.Authors = []cli.Author{
		{
			Name:  "Sebastian Johnsson",
			Email: "",
		},
	}

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "network",
			Usage: "Which network to use (valid options: localnet, devnet, testnet, mainnet)",
			Value: "",
		},

		cli.StringFlag{
			Name:  "config",
			Usage: "The path to the config containing the test suite settings",
			Value: "./config.yml",
		},

		cli.StringFlag{
			Name:  "funding-address",
			Usage: "Which address to use to fund test accounts (tokens will also be returned to this address",
			Value: "",
		},

		cli.Float64Flag{
			Name:  "minimum-funds",
			Usage: "The minimum funds a source wallet needs to have to be included in the funding process",
			Value: 10.0,
		},

		cli.StringFlag{
			Name:  "passphrase",
			Usage: "Passphrase to use for unlocking the keystores",
			Value: "",
		},

		cli.StringFlag{
			Name:  "keys",
			Usage: "Where the wallet keys are located",
			Value: "",
		},
	}

	app.Action = func(context *cli.Context) error {
		return startTests(context)
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Println("error: " + err.Error())
		os.Exit(1)
	}
}

func startTests(context *cli.Context) error {
	err := config.Configure(context)
	accs, err := accounts.LoadSourceAccounts()

	if err != nil {
		return err
	}

	if err = funding.SetupFundingAccount(accs); err != nil {
		return err
	}

	executeTests(accs)
	testResults()

	return nil
}

func executeTests(accs []string) {
	for _, testCase := range testcases.Registry {
		testCase.Function.(func())()
	}
}

func testResults() {
	successfulCount := 0
	failedCount := 0
	
	for _, testCase := range testcases.Results {
		if testCase.Result == testCase.Expected {
			successfulCount++
		} else {
			failedCount++
		}
	}

	fmt.Println("------------------------------------------------------------")
	fmt.Println(fmt.Sprintf("Test suite status - executed a total of %d test case(s):", len(testcases.Results)))
	fmt.Println(fmt.Sprintf("Successful: %d", successfulCount))
	fmt.Println(fmt.Sprintf("Failed: %d", failedCount))
	fmt.Println("------------------------------------------------------------\n")

	fmt.Println("Executed test cases:")
	fmt.Println("------------------------------------------------------------")
	for _, testCase := range testcases.Results {
		if testCase.Result == testCase.Expected {
			fmt.Println(fmt.Sprintf("Testcase %s: success", testCase.Name))
		} else {
			fmt.Println(fmt.Sprintf("Testcase %s: failed", testCase.Name))
		}
	}
	fmt.Println("------------------------------------------------------------\n")
}
