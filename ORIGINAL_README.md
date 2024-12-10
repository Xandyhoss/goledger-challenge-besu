# GoLedger Challenge - Besu Edition

On this challenge, you will interact with a Besu node. The goal is to create a simple application that will interact with a Besu node to transact in a smart contract, check the value of a smart contract variable and sync that value to an external database.

To accomplish that, we recommend you use a UNIX-like machine (Linux/macOS). Besides that, we will need to install NPM/NPX, Hardhat and Docker.

## Install the prerequisites

- Install NPM and NPX (https://www.npmjs.com/get-npm)
- Install Hardhat (https://hardhat.org/getting-started/)
- Install Docker and Docker Compose (https://www.docker.com/)
- Install Besu (https://besu.hyperledger.org/private-networks/get-started/install/binary-distribution)
- Install Go (https://golang.org/dl/)
- Fork the repository https://github.com/goledgerdev/goledger-challenge-besu 
    - Fork it, do **NOT** clone it, since you will need to send us your forked repository
	- If you cannot fork it, create a private repository and give access to `samuelvenzi` and `Tubar2`

### Hardhat installation details

Hardhat is a development environment to compile, deploy, test, and debug your Ethereum software. It helps developers manage and automate the recurring tasks that are inherent to the process of building smart contracts and dApps.

To install Hardhat, you need to have Node.js installed. If you don't have it, you can download it [here](https://nodejs.org/).

After installing Node.js, you can install Hardhat by running the following command:

```bash
npm install --save-dev hardhat
```

Note: Your system might require a slightly different command to install Hardhat. Check the [Hardhat installation guide](https://hardhat.org/getting-started/) for more information.

## Set up the environment

To set up the environment, you need to fork this repository. Make sure you have installed the requirements. To set up the environment, you need to run the following commands:

```bash
cd besu
./startDev.sh
```

This will bring up a local Besu netwwork with 4 nodes. You can check the logs of each node by running the following command:

```bash
docker logs -f besu_node-0
```

This will also deploy a smart contract to the network. The contract is a simple storage contract that has a variable that can be set and get. Note that it will log the contracts address, which will be important later. If you want to check the contract's source code, you can find it in the `contracts` folder. The contract's ABI can be found in the `/besu/artifacts/contracts/SimpleStorage.sol/SimpleStorage.json` file.

# The challenge

Your task is to create a simple application that interacts with a Besu blockchain network and an SQL database. The application should be implemented in Go and expose its functionality as either a REST API or a gRPC service.

## Requirements

1. **Programming Language:**
   - The application must be written in Go.

2. **API Type:**
   - Choose either REST or gRPC for the service interface.
   - If implementing gRPC, enable reflection so we can test it using tools like Postman.

3. **Database Integration:**
   - Use an SQL database (e.g., PostgreSQL or MySQL).
   - Store the value of the smart contract variable in the database.

4. **Endpoints:**
   - The application should provide the following functionality via appropriately named endpoints or methods:

     1. **SET:**
        - Set a new value for the smart contract variable.
        - The application should send this value to the deployed smart contract on the Besu network.

     2. **GET:**
        - Retrieve the current value of the smart contract variable from the blockchain.

     3. **SYNC:**
        - Synchronize the value of the smart contract variable from the blockchain to the SQL database.

     4. **CHECK:**
        - Compare the value stored in the database with the current value of the smart contract variable.
        - Return `true` if they are the same, otherwise return `false`.

   - **Endpoint Naming:**
     - You may name the endpoints/methods as you see fit, provided their functionality meets the requirements outlined above.

   - **General Notes:**
     - The Besu network will have a smart contract deployed that includes a single variable to store a value (similar to a SimpleStorage contract).
     - Ensure the application handles blockchain interactions (reads/writes) correctly.
     - Add appropriate error handling for all interactions (blockchain, database, and API).

## Deliverables

1. **Source Code:**
   - The source code of the application should be hosted on a public GitHub repository forked from this one.
   - Include a README file with instructions on how to run the application.
2. **Documentation:**
   - Provide a brief explanation of the application's architecture and how it interacts with the Besu network and the SQL database.
   - Include any additional information you think is relevant.
   - This can be done in the README file or as a separate Markdown file.

Remember to commit your changes to your forked repository. Commits will be used during the evaluation process.

## Interaction with the Besu network

To interact with the Besu network, you can use the Go Ethereum client. Below we provide two functions that interact with the Besu network, one for writing data (`ExecContract`) and one for reading data (`CallContract`). Feel free to include and change this function in your application.

```go
package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func ExecContract() {
	abi, err := abi.JSON(strings.NewReader("REPLACE: abi JSON as string goes here")) // found under besu/artifacts/contracts/SimpleStorage.sol/SimpleStorage.json
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := ethclient.DialContext(ctx, "REPLACE: network URL") // e.g., http://localhost:8545
	if err != nil {
		log.Fatalf("error dialing node: %v", err)
	}

	slog.Info("querying chain id")

	chainId, err := client.ChainID(ctx)
	if err != nil {
		log.Fatalf("error querying chain id: %v", err)
	}
	defer client.Close()

	contractAddress := common.HexToAddress("REPLACE: contract address") // will be returned during startDev.sh execution

	boundContract := bind.NewBoundContract(
		contractAddress,
		abi,
		client,
		client,
		client,
	)

	priv, err := crypto.HexToECDSA("REPLACE: private key") // this can be found in the genesis.json file
	if err != nil {
		log.Fatalf("error loading private key: %v", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(priv, chainId)
	if err != nil {
		log.Fatalf("error creating transactor: %v", err)
	}

	tx, err := boundContract.Transact(auth, "REPLACE: method name")
	if err != nil {
		log.Fatalf("error transacting: %v", err)
	}

	fmt.Println("waiting until transaction is mined",
		"tx", tx.Hash().Hex(),
	)

	receipt, err := bind.WaitMined(
		context.Background(),
		client,
		tx,
	)
	if err != nil {
		log.Fatalf("error waiting for transaction to be mined: %v", err)
	}

	fmt.Printf("transaction mined: %v\n", receipt)
}
```

You can also use the following code to call `view` functions on the contract.

```go
package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func CallContract()  {
	var result interface{}

	abi, err := abi.JSON(strings.NewReader("REPLACE: abi JSON as string goes here")) // found under besu/artifacts/contracts/SimpleStorage.sol/SimpleStorage.json
	if err != nil {
		log.Fatalf("error parsing abi: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := ethclient.DialContext(ctx, "REPLACE: network URL") // e.g., http://localhost:8545
	if err != nil {
		log.Fatalf("error connecting to eth client: %v", err)
	}
	defer client.Close()

	contractAddress := common.HexToAddress("REPLACE: contract address") // will be returned during startDev.sh execution
	caller := bind.CallOpts{
		Pending: false,
		Context: ctx,
	}

	boundContract := bind.NewBoundContract(
		contractAddress,
		abi,
		client,
		client,
		client,
	)

	var output []interface{}
	err = boundContract.Call(&caller, &output, "REPLACE: method name")
	if err != nil {
		log.Fatalf("error calling contract: %v", err)
	}
	result = output

	fmt.Println("Successfully called contract!", result)
}
```

To complete the challenge, you must send us the link to your repository with the alterations you made.
