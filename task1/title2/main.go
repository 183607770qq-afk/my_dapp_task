package main

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"time"

	"title2/count"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

const (
	contractBytecode = "6080604052348015600e575f5ffd5b505f5f81905550610132806100225f395ff3fe6080604052348015600e575f5ffd5b50600436106030575f3560e01c8063a87d942c146034578063d09de08a14604e575b5f5ffd5b603a6056565b60405160459190608c565b60405180910390f35b6054605e565b005b5f5f54905090565b60015f5f828254606d919060d0565b92505081905550565b5f819050919050565b6086816076565b82525050565b5f602082019050609d5f830184607f565b92915050565b7f4e487b71000000000000000000000000000000000000000000000000000000005f52601160045260245ffd5b5f60d8826076565b915060e1836076565b925082820190508082111560f65760f560a3565b5b9291505056fea2646970667358221220bb12840bf20785318a6683e66fa752325e613b7e9666645ecab62962d34c2d6464736f6c634300081f0033"
)

func main() {
	client, err := ethclient.Dial("https://sepolia.infura.io/v3/72ee2f483626429ab30c674d52862ef7")
	if err != nil {
		log.Fatal(err)
	}
	deployContract(client)

	// instance, err := count.NewCount(receipt.ContractAddress, client)
	instance, err := count.NewCount(common.HexToAddress("0x6e8aa187e78CCF1f5B9e7c2e5C117fB5459b813b"), client)
	if err != nil {
		log.Fatal(err)
	}
	value, err := instance.GetCount(&bind.CallOpts{Context: context.Background()})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("initial count value: %d\n", value.Uint64())

}

func deployContract(client *ethclient.Client) {
	privateKey, err := crypto.HexToECDSA("97c2242873584e7a8a5e20456e74dca8a2ca4d9252b8916a9dda8615b607fcd6")
	if err != nil {
		log.Fatal(err)
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}
	gasprice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	data, err := hex.DecodeString(contractBytecode)
	if err != nil {
		log.Fatal(err)
	}
	tx := types.NewContractCreation(nonce, big.NewInt(0), 58078, gasprice, data)
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	typedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatal(err)
	}
	err = client.SendTransaction(context.Background(), typedTx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("tx sent: %s\n", typedTx.Hash().Hex())
	receipt, err := waitForReceipt(client, typedTx.Hash())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("contract deployed at address: %s\n", receipt.ContractAddress.Hex())
}
func waitForReceipt(client *ethclient.Client, txHash common.Hash) (*types.Receipt, error) {
	for {
		receipt, err := client.TransactionReceipt(context.Background(), txHash)
		if err == nil {
			return receipt, nil
		}
		if err != ethereum.NotFound {
			return nil, err
		}
		// 等待一段时间后再次查询
		time.Sleep(1 * time.Second)
	}
}
