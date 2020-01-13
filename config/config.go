package config

import (
	"errors"
	"fmt"

	"github.com/SebastianJ/harmony-tx-tests/utils"
	"github.com/harmony-one/go-sdk/pkg/sharding"
	"github.com/urfave/cli"
	"gopkg.in/yaml.v2"
)

// Configuration - the central configuration for the test suite tool
var Configuration Config

// Config - represents the config
type Config struct {
	Network Network `yaml:"network"`
	Account Account `yaml:"account"`
	Funding Funding `yaml:"funding"`
}

// Network - represents the network settings group
type Network struct {
	Name    string  `yaml:"name"`
	Node    string  `yaml:"node"`
	GasCost float64 `yaml:"gas_cost"`
	Shards  int
}

// Account - represents the account settings group
type Account struct {
	KeysPath   string `yaml:"keys_path"`
	Passphrase string `yaml:"passphrase"`
}

// Funding - represents the funding settings group
type Funding struct {
	Account              FundingAccount `yaml:"account"`
	MinimumFunds         float64        `yaml:"minimum_funds"`
	ConfirmationWaitTime int            `yaml:"confirmation_wait_time"`
	Attempts             int            `yaml:"attempts"`
	GasPrice             int64          `yaml:"gas_price"`
}

// FundingAccount - represents a funding account
type FundingAccount struct {
	Name    string `yaml:"name"`
	Address string `yaml:"address"`
}

// Configure - configures the test suite tool using a combination of the YAML config file as well as command arguments
func Configure(context *cli.Context) (err error) {
	configPath := context.GlobalString("config")

	if err = loadYamlConfig(configPath); err != nil {
		return err
	}

	if network := context.GlobalString("network"); network != "" && network != Configuration.Network.Name {
		Configuration.Network.Name = network
	}

	Configuration.Network.Name = NormalizedNetworkName(Configuration.Network.Name)

	if Configuration.Network.Name == "" {
		return errors.New("you need to specify a valid network name to use! Valid options: localnet, devnet, testnet, mainnet")
	}

	if err = setNetworkConfig(); err != nil {
		return err
	}

	if passphrase := context.GlobalString("passphrase"); passphrase != "" && passphrase != Configuration.Account.Passphrase {
		Configuration.Account.Passphrase = passphrase
	}

	if keysPath := context.GlobalString("keys"); keysPath != "" && keysPath != Configuration.Account.KeysPath {
		Configuration.Account.KeysPath = keysPath
	}

	if fundingAddress := context.GlobalString("funding-address"); fundingAddress != "" && fundingAddress != Configuration.Funding.Account.Address {
		Configuration.Funding.Account.Address = fundingAddress
	}

	if minimumFunds := context.GlobalFloat64("minimum-funds"); minimumFunds != 0.0 && minimumFunds != Configuration.Funding.MinimumFunds {
		Configuration.Funding.MinimumFunds = minimumFunds
	}

	Configuration.Account.KeysPath = fmt.Sprintf("%s/%s", Configuration.Account.KeysPath, Configuration.Network.Name)

	return nil
}

func loadYamlConfig(path string) error {
	Configuration = Config{}
	yamlData, err := utils.ReadFileToString(path)

	if err != nil {
		return err
	}

	err = yaml.Unmarshal([]byte(yamlData), &Configuration)

	if err != nil {
		return err
	}

	return nil
}

func setNetworkConfig() error {
	Configuration.Network.Node = GenerateNodeAddress(Configuration.Network.Name, 0)

	shardingStructure, err := sharding.Structure(Configuration.Network.Node)
	if err != nil {
		return err
	}

	Configuration.Network.Shards = len(shardingStructure)

	return nil
}

// NormalizedNetworkName - return a normalized network name
func NormalizedNetworkName(network string) string {
	switch network {
	case "local", "localnet":
		return "localnet"
	case "dev", "devnet", "pga":
		return "devnet"
	case "testnet", "pangaea", "p":
		return "testnet"
	case "mainnet", "main", "t":
		return "mainnet"
	default:
		return ""
	}
}

// GenerateNodeAddress - generates a node address given a network and a shardID
func GenerateNodeAddress(network string, shardID uint32) string {
	var node string

	switch network {
	case "local", "localnet":
		node = fmt.Sprintf("http://localhost:950%d", shardID)
	case "dev", "devnet", "pga":
		node = fmt.Sprintf("https://api.s%d.pga.hmny.io", shardID)
	case "testnet", "pangaea", "p":
		node = fmt.Sprintf("https://api.s%d.p.hmny.io", shardID)
	case "mainnet", "main", "t":
		node = fmt.Sprintf("https://api.s%d.t.hmny.io", shardID)
	default:
		node = fmt.Sprintf("http://localhost:950%d", shardID)
	}

	return node
}
