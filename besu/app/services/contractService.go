package services

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"os"
	"strings"
	"time"

	"github.com/Xandyhoss/goledger-challenge-besu/app/db"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func loadABI(filePath string) (*abi.ABI, error) {
	type Artifact struct {
		ABI json.RawMessage `json:"abi"`
	}

	abiContent, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read ABI file: %w", err)
	}

	var artifact Artifact
	if err := json.Unmarshal(abiContent, &artifact); err != nil {
		return nil, fmt.Errorf("failed to unmarshal artifact: %w", err)
	}

	parsedABI, err := abi.JSON(strings.NewReader(string(artifact.ABI)))
	if err != nil {
		return nil, fmt.Errorf("failed to parse ABI: %w", err)
	}
	return &parsedABI, nil
}

func createClient() (*ethclient.Client, *big.Int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := ethclient.DialContext(ctx, os.Getenv("BESU_NODE_URL"))
	if err != nil {
		return nil, nil, fmt.Errorf("error dialing node: %w", err)
	}

	chainID, err := client.ChainID(ctx)
	if err != nil {
		client.Close()
		return nil, nil, fmt.Errorf("error querying chain ID: %w", err)
	}

	return client, chainID, nil
}

func createTransactor(chainID *big.Int) (*bind.TransactOpts, error) {
	privateKey, err := crypto.HexToECDSA(os.Getenv("PRIVATE_KEY"))
	if err != nil {
		return nil, fmt.Errorf("error loading private key: %w", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return nil, fmt.Errorf("error creating transactor: %w", err)
	}

	return auth, nil
}

func ExecContract(value uint) (*uint, error) {
	convertedValue := big.NewInt(int64(value))
	abi, err := loadABI("artifacts/contracts/SimpleStorage.sol/SimpleStorage.json")
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("error loading ABI: %w", err)
	}

	client, chainID, err := createClient()
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("error creating client: %w", err)
	}
	defer client.Close()

	auth, err := createTransactor(chainID)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("error creating transactor: %w", err)
	}

	contractAddress := common.HexToAddress(os.Getenv("CONTRACT_ADDRESS"))
	boundContract := bind.NewBoundContract(contractAddress, *abi, client, client, client)

	tx, err := boundContract.Transact(auth, "set", convertedValue)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("error sending transaction: %w", err)
	}

	fmt.Println("Transaction sent:", tx.Hash().Hex())

	receipt, err := bind.WaitMined(context.Background(), client, tx)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("error waiting for transaction: %w", err)
	}

	fmt.Println("Transaction mined:", receipt)

	returnValue := value
	return &returnValue, nil
}

func CallContract() (*big.Int, error) {
	var result *big.Int

	abi, err := loadABI("artifacts/contracts/SimpleStorage.sol/SimpleStorage.json")
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("error loading ABI: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, _, err := createClient()
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("error creating client: %v", err)
	}
	defer client.Close()

	contractAddress := common.HexToAddress(os.Getenv("CONTRACT_ADDRESS"))
	boundContract := bind.NewBoundContract(contractAddress, *abi, client, client, client)

	caller := bind.CallOpts{
		Pending: false,
		Context: ctx,
	}

	outputArr := []interface{}{}
	err = boundContract.Call(&caller, &outputArr, "get")
	if err != nil {
		log.Println("error calling contract: %w", err)
		return nil, fmt.Errorf("error calling contract: %v", err)
	}
	result = outputArr[0].(*big.Int)

	fmt.Println("Successfully called contract!", result)
	return result, nil
}

func CheckContract() bool {
	database, err := db.ConnectDB()
	if err != nil {
		log.Println("error connecting to database: %w", err)
		return false
	}

	var contract db.Contract

	result := database.First(&contract, 1)
	if result.Error != nil {
		log.Println("error querying contract from db: %w", result.Error)
		return false
	}

	calledContractValue, err := CallContract()
	if err != nil {
		log.Println("error calling smart contract: %w", err)
		return false
	}

	return calledContractValue.Uint64() == uint64(contract.ContractNumber)
}

func SyncContract() bool {
	database, err := db.ConnectDB()
	if err != nil {
		log.Println("error connecting to database: %w", err)
		return false
	}

	calledContractValue, err := CallContract()
	if err != nil {
		log.Println("error calling smart contract: %w", err)
		return false
	}

	var contract db.Contract
	result := database.First(&contract, 1)
	if result.Error != nil {
		log.Println("error querying contract from db: %w", result.Error)
		return true
	}

	database.Model(&contract).Update("contract_number", calledContractValue.Uint64())
	return true
}
