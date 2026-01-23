package main

import (
	"byd50-ssi/pkg/did/core/dids"
	"byd50-ssi/pkg/did/core/driver/scdid"
	"byd50-ssi/pkg/did/kms"
	"context"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
	"os"
	"strings"
)

func main() {
	rpcURL := os.Getenv("ETH_RPC_URL")
	if rpcURL == "" {
		log.Fatal("missing ETH_RPC_URL")
	}
	registryAddress := os.Getenv("ETH_REGISTRY_ADDRESS")
	if registryAddress == "" {
		log.Fatal("missing ETH_REGISTRY_ADDRESS")
	}

	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		log.Fatalf("failed to connect to the Ethereum client: %v", err)
	}

	instance, err := scdid.NewScdid(common.HexToAddress(registryAddress), client)
	if err != nil {
		log.Fatalf("failed to load registry contract: %v", err)
	}

	resolveDid := os.Getenv("ETH_RESOLVE_DID")
	if resolveDid == "" {
		resolveDid = "did:byd50:12345555"
	}
	docs, err := instance.ResolveDid(&bind.CallOpts{Pending: true}, resolveDid)
	if err != nil {
		log.Printf("ResolveDid failed: %v", err)
	} else {
		log.Printf("ResolveDid(%s) => %s", resolveDid, docs)
	}

	privateKeyHex := strings.TrimPrefix(os.Getenv("ETH_PRIVATE_KEY_HEX"), "0x")
	if privateKeyHex == "" {
		log.Printf("skipping CreateDid; set ETH_PRIVATE_KEY_HEX to write to chain")
		return
	}

	chainIDStr := os.Getenv("ETH_CHAIN_ID")
	if chainIDStr == "" {
		log.Fatal("missing ETH_CHAIN_ID")
	}
	chainID, ok := new(big.Int).SetString(chainIDStr, 10)
	if !ok {
		log.Fatalf("invalid ETH_CHAIN_ID: %s", chainIDStr)
	}

	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		log.Fatalf("invalid ETH_PRIVATE_KEY_HEX: %v", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		log.Fatalf("failed to create transactor: %v", err)
	}
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Printf("failed to suggest gas price: %v", err)
		gasPrice = big.NewInt(0)
	}

	localKMS, err := kms.InitKMS(kms.KeyTypeECDSA)
	if err != nil {
		log.Fatalf("failed to init KMS: %v", err)
	}
	createdDid, createdDoc := dids.CreateDID("eth", localKMS.PbKeyBase58())
	log.Printf("CreateDid input => %s", createdDid)

	tx, err := instance.CreateDid(&bind.TransactOpts{
		From:     auth.From,
		Signer:   auth.Signer,
		GasPrice: gasPrice,
		GasLimit: 3141592,
	}, createdDid, string(createdDoc))
	if err != nil {
		log.Printf("CreateDid failed: %v", err)
		return
	}

	log.Printf("CreateDid tx hash: %s", tx.Hash())
}
