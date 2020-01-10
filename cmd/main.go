package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/SebastianJ/harmony-tx-tests/accounts"
	"github.com/SebastianJ/harmony-tx-tests/testcases"
	"github.com/SebastianJ/harmony-tx-tests/utils"

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
			Name:  "node",
			Usage: "Which node endpoint to use for API commands",
			Value: "https://api.s0.pga.hmny.io",
		},

		cli.StringFlag{
			Name:  "passphrase",
			Usage: "Passphrase to use for unlocking the keystore",
			Value: "",
		},

		cli.StringFlag{
			Name:  "keys",
			Usage: "Where the wallet keys are located",
			Value: "./keys",
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
	node := context.GlobalString("node")
	passphrase := context.GlobalString("passphrase")

	if node == "" {
		return errors.New("you need to specify a node to use for the API calls")
	}

	keysPath, _ := filepath.Abs(context.GlobalString("keys"))
	keyFiles, err := utils.IdentifyKeyFiles(keysPath)

	if err != nil {
		return err
	}

	fmt.Println(fmt.Sprintf("Found a total of %d keys", len(keyFiles)))

	accs := make(map[string]string)

	for _, keyFile := range keyFiles {
		//fmt.Println(fmt.Sprintf("Keyfile path found: %s", keyFile))

		keyDetails := utils.ParseKeyDetailsFromKeyFile(keyFile)
		keyName := fmt.Sprintf("_tx_gen_NEW_FUNDED_ACC_%s", keyDetails["id"])

		err := accounts.ImportAccount(keyFile, keyName, passphrase, keyDetails)

		if err != nil {
			return err
		}

		accs[keyName] = keyDetails["address"]
	}

	//err = balances.OutputBalanceStatusForAddresses(accs, node)

	testcaseStatuses := make(map[string]testcases.TestCase)
	testcaseStatuses["SBS1"] = testcases.Sbs1TestCase(accs, passphrase, node)

	testResults(testcaseStatuses)

	return nil
}

func testResults(testcaseStatuses map[string]testcases.TestCase) {
	successfulCount := 0
	failedCount := 0

	for _, testCase := range testcaseStatuses {
		if testCase.Result == testCase.Expected {
			successfulCount++
		} else {
			failedCount++
		}
	}

	fmt.Println("\n------------------------------------------------------------")
	fmt.Println(fmt.Sprintf("Test suite status - executed a total of %d test case(s):", len(testcaseStatuses)))
	fmt.Println(fmt.Sprintf("Successful: %d", successfulCount))
	fmt.Println(fmt.Sprintf("Failed: %d", failedCount))
	fmt.Println("------------------------------------------------------------\n")

	fmt.Println("Executed test cases:")
	fmt.Println("------------------------------------------------------------")
	for testCaseName, testCase := range testcaseStatuses {
		if testCase.Result == testCase.Expected {
			fmt.Println(fmt.Sprintf("Testcase %s: success", testCaseName))
		} else {
			fmt.Println(fmt.Sprintf("Testcase %s: failed", testCaseName))
		}
	}
	fmt.Println("------------------------------------------------------------\n")
}
