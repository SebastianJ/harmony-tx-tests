package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

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

	if node == "" {
		return errors.New("you need to specify a node to use for the API calls")
	}

	keysPath, _ := filepath.Abs(context.GlobalString("keys"))
	//keyFiles, err := utils.IdentifyKeyFiles(keysPath)

	return nil
}
