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
	deploy2()
	// 	client, err := ethclient.Dial("https://sepolia.infura.io/v3/72ee2f483626429ab30c674d52862ef7")
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	deployContract(client)

	// // 检查节点同步状态
	// syncProgress, err := client.SyncProgress(context.Background())
	// if err != nil {
	//     log.Fatal("获取同步状态失败: ", err)
	// }
	// if syncProgress != nil {
	//     log.Printf("警告：节点正在同步，当前区块：%d / 最新区块：%d\n", syncProgress.CurrentBlock, syncProgress.HighestBlock)
	//     // 如果不同步，需要等待或更换节点
	// }

	// contractAddr := common.HexToAddress("0x26c61e4a28EC6ec78C32cbe1B3Dea457c001DA9E")
	// contractAddr := common.HexToAddress("0x6e8aa187e78CCF1f5B9e7c2e5C117fB5459b813b")
	// contractAddr := common.HexToAddress("0xD9A7a44E20f7eb4788Ac0CE89B1AbE73f7BE08ad")
	//   contractAddr := common.HexToAddress("0xD8cC5681a32C8E92b12C0D3fbA95633bf5E0c646")
	//   fmt.Printf("调用合约地址: %s\n", contractAddr.Hex())

	//   code, err := client.CodeAt(context.Background(), contractAddr, nil)
	//   if err != nil {
	//       log.Fatal("检查合约代码失败: ", err)
	//   }
	//   if len(code) == 0 {
	//       log.Fatal("❌ 该地址没有合约代码！请确认：\n" +
	//           "1. 地址是否正确\n" +
	//           "2. 是否已部署成功\n" +
	//           "3. 是否连接了正确的网络（Sepolia）")
	//   }
	//   fmt.Printf("✅ 合约代码长度: %d 字节\n", len(code))

	// 	// instance, err := count.NewCount(receipt.ContractAddress, client)
	// 	instance, err := count.NewCount(contractAddr, client)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	instance.Increment(&bind.TransactOpts{Context: context.Background()})
	// 	value, err := instance.GetCount(&bind.CallOpts{Context: context.Background()})
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	fmt.Printf("initial count value: %d\n", value.Uint64())

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
	if receipt.Status != 1 {
		log.Fatal("Contract deployment transaction failed! Receipt Status: ", receipt.Status)
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
func increment() {
	client, err := ethclient.Dial("https://sepolia.infura.io/v3/72ee2f483626429ab30c674d52862ef7")
	if err != nil {
		log.Fatal(err)
	}

	// privateKey, err := crypto.GenerateKey()
	// privateKeyBytes := crypto.FromECDSA(privateKey)
	// privateKeyHex := hex.EncodeToString(privateKeyBytes)
	// fmt.Println("Private Key:", privateKeyHex)
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

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	chainId, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainId)
	if err != nil {
		log.Fatal(err)
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)     // in wei
	auth.GasLimit = uint64(300000) // in units
	auth.GasPrice = gasPrice

	instance, err := count.NewCount(common.HexToAddress("0x6fDF97B4A3dcF8CCb5A9EDCDBB69137Cc0638512"), client)
	if err != nil {
		log.Fatal(err)
	}
	value, err := instance.GetCount(&bind.CallOpts{Context: context.Background()})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("initial count value: %d\n", value.Uint64())
}

func deploy2() {
	client, err := ethclient.Dial("https://sepolia.infura.io/v3/72ee2f483626429ab30c674d52862ef7")
	if err != nil {
		log.Fatal(err)
	}

	// privateKey, err := crypto.GenerateKey()
	// privateKeyBytes := crypto.FromECDSA(privateKey)
	// privateKeyHex := hex.EncodeToString(privateKeyBytes)
	// fmt.Println("Private Key:", privateKeyHex)
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

	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	chainId, err := client.NetworkID(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainId)
	if err != nil {
		log.Fatal(err)
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)     // in wei
	auth.GasLimit = uint64(300000) // in units
	auth.GasPrice = gasPrice

	// input := "1.0"
	address, tx, instance, err := count.DeployCount(auth, client)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(address.Hex())
	fmt.Println(tx.Hash().Hex())

	// _ = instance
	// instance, err := count.NewCount(receipt.ContractAddress, client)
	// instance, err := count.NewCount(address, client)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	 ctx := context.Background()
	receipt, err := bind.WaitMined(ctx, client, tx)
    if err != nil {
        log.Fatal("等待交易确认失败:", err)
    }
	if receipt.Status != 1 {
        log.Fatal("❌ 交易失败，合约未部署")
    }
	fmt.Printf("✅ 交易已确认！区块: %d\n", receipt.BlockNumber)
    fmt.Println("⏳ 给节点一些时间同步...")
    time.Sleep(3 * time.Second) // 额外等待确保节点同步


	value, err := instance.GetCount(&bind.CallOpts{Context: context.Background()})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("initial count value: %d\n", value.Uint64())
}
