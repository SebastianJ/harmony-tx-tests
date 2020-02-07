#!/usr/bin/env bash

echo "Installing Harmony Tx Tests"
curl -LOs http://tools.harmony.one.s3.amazonaws.com/release/linux-x86_64/harmony-tx-tests/harmony-tx-tests && chmod u+x harmony-tx-sender
curl -LOs http://tools.harmony.one.s3.amazonaws.com/release/linux-x86_64/harmony-tx-tests/testcases.tar.gz && tar -xzvf testcases.tar.gz && rm -rf testcases.tar.gz
curl -LOs https://raw.githubusercontent.com/SebastianJ/harmony-tx-tests/master/config.yml
mkdir -p keys/testnet
echo "Harmony Tx Tests have now been downloaded"
echo "Make sure to either add keystore files to keys/testnet or create the file keys/testnet/private_keys.txt and add testnet private keys to it"
echo "The tests can be run on multiple network - testnet is just used as an example. To run the framework on other networks, e.g. localnet, use --network localnet"
echo "When you've added keyfiles or private keys, invoke the tests using ./harmony-tx-tests"
echo "To see all available configuration options, run ./harmony-tx-tests --help or check out the configuration specified in config.yml"
