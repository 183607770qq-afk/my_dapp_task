package main

import (
    "context"
    "crypto/ecdsa"
    "fmt"
    "log"
    "math/big"

    "github.com/ethereum/go-ethereum/common"
    "github.com/ethereum/go-ethereum/core/types"
    "github.com/ethereum/go-ethereum/crypto"
    "github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	client, err := ethclient.Dial("https://sepolia.infura.io/v3/72ee2f483626429ab30c674d52862ef7")
	if err != nil {
		log.Fatal(err)
	}
	// sendTransaction(*client)
	selectTransaction(*client)
  
}

func selectTransaction(client ethclient.Client) {
	// txHash := common.HexToHash("0x6d3afc0ff3cca16c5733e5bc1eea4c804eab7c01d344a8e795c61bb9f5cb4c1e")
	// tx,isPending, err := client.TransactionByHash(context.Background(), txHash)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// if isPending {
	// 	fmt.Print("tx is pending\n")
	// } else {
	// 	fmt.Print("tx is not pending\n")
	// }
	// fmt.Printf("tx hash: %s\n", tx.Hash().Hex())
	// fmt.Printf("tx value: %s\n", tx.Value().String())
	// fmt.Printf("tx gas: %d\n", tx.Gas())	
	// fmt.Printf("tx gas price: %s\n", tx.GasPrice().String())	
	// fmt.Printf("tx data: %x\n", tx.Data())
	// fmt.Printf("tx nonce: %d\n", tx.Nonce())
	// fmt.Printf("tx to: %s\n", tx.To().Hex())

    blockNumber ,err := client.BlockNumber(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("latest block number: %d\n", blockNumber)
	block, err :=client.BlockByNumber(context.Background(),nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("block hash: %s\n", block.Hash().Hex())
	fmt.Printf("block number: %d\n", block.Number().Uint64())
	fmt.Printf("block time: %d\n", block.Time())
	fmt.Printf("block txs: %d\n", block.Transactions().Len())
	fmt.Printf("block difficulty: %s\n", block.Difficulty().String())
}
func sendTransaction(client ethclient.Client) {
	privateKey, err := crypto.HexToECDSA("97c2242873584e7a8a5e20456e74dca8a2ca4d9252b8916a9dda8615b607fcd6")
	if err != nil {
		log.Fatal(err)
	}
	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	    	crypto.PubkeyToAddress(*publicKeyECDSA)
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce,err := client.PendingNonceAt(context.Background(),fromAddress)
	if err != nil {
		log.Fatal(err)
	}
	value := big.NewInt(10000000000000000) // in wei (0 eth)
	gasLimit := uint64(21000)              // in units
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	toAddress := common.HexToAddress("0xef798e10e69952d110ccda6fc8df9e178c77cedc")
	var data []byte
	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, data)
	chainID, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	signedTx,err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
	if err != nil {
		log.Fatal(err)
	}
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("sended tx: %s", signedTx.Hash().Hex()) //0x6d3afc0ff3cca16c5733e5bc1eea4c804eab7c01d344a8e795c61bb9f5cb4c1e
	fmt.Println()
}


