package config

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
