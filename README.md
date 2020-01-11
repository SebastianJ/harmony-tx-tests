# harmony-tx-tests
Harmony tx tests is a transaction test suite tool.

It has support for:

* Automatically importing keys placed in keys/ into the keystore
* Generating temporary receiver accounts
* Sending back any eventual test funds to the originator and subsequently removing the account from the keystore
* Defining test cases and evaluating if a given test case's result matches the expected test result.
* Sending transactions without relying on hmy or any CLI - i.e. directly communication with the API:s

