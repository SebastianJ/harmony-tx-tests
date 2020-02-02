package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/SebastianJ/harmony-tf/config"
	"github.com/SebastianJ/harmony-tf/testing"

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
			Name:  "path",
			Usage: "The root path for the config file + the testcases",
			Value: "./",
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
	basePath, err := filepath.Abs(context.GlobalString("path"))
	if err != nil {
		return err
	}

	if err := config.Configure(basePath, context); err != nil {
		return err
	}

	if err := testing.ExecuteTestCases(); err != nil {
		return err
	}

	return nil
}
