//go:build eth
// +build eth

package driver

import (
	"byd50-ssi/pkg/did/configs"
	"byd50-ssi/pkg/did/core/dids"
	"byd50-ssi/pkg/did/core/driver/scdid"
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"errors"
	"fmt"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

var (
	// ErrResolveDidEth - Sadly this is missing from method/eth registry
	ErrResolveDidEth = errors.New("method/eth: resolve did error")
)

// DidMethodETH - Implements the ETH did methods
type DidMethodETH struct {
	Name string
}

// Specific instances for EC256 and company
var (
	didMethodETH *DidMethodETH
)

func init() {
	// Register eth driver for Ethereum Smart Contract
	didMethodETH = &DidMethodETH{"eth"}
	RegisterDidMethod(didMethodETH.Method(), func() DidMethod {
		return didMethodETH
	})
}

func (m *DidMethodETH) Method() string {
	return m.Name
}

func toECDSAFromHex(hexString string) (*ecdsa.PrivateKey, error) {
	pk := new(ecdsa.PrivateKey)
	pk.D, _ = new(big.Int).SetString(hexString, 16)
	pk.PublicKey.Curve = elliptic.P256()
	pk.PublicKey.X, pk.PublicKey.Y = pk.PublicKey.Curve.ScalarBaseMult(pk.D.Bytes())
	return pk, nil
}

func loadEthPrivateKeyHex() (string, error) {
	privateKeyHex := os.Getenv("ETH_PRIVATE_KEY_HEX")
	if privateKeyHex == "" {
		return "", fmt.Errorf("missing ETH_PRIVATE_KEY_HEX env for eth driver")
	}
	return privateKeyHex, nil
}

// CreateDid - Implements the RegisterDid method from DidMethod
// For this register did method, pbKeyBase58 must be an base58 encoded string
func (m *DidMethodETH) CreateDid(pbKeyBase58 string) (string, error) {
	// Create an IPC based RPC connection to a remote node
	client, err := ethclient.Dial(configs.UseConfig.EthClientUrl)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	// Instantiate the contract and display its name
	address := common.HexToAddress(configs.UseConfig.EthClientScAddress)
	instance, err := scdid.NewScdid(address, client)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("contract is loaded")

	//myDKMS := kms.GetKMS()
	//pvKeyECDSA, ok := myDKMS.PvKey().(*ecdsa.PrivateKey)
	//pvKeyECDSA, err := toECDSAFromHex("e757f43e4d62c271e0e4713fe8e33f8451fc379026ed01dd02a933b9ae750c9d")
	privateKeyHex, err := loadEthPrivateKeyHex()
	if err != nil {
		return "", err
	}
	pvKeyECDSA, err := crypto.HexToECDSA(privateKeyHex)

	if err != nil {
		log.Fatal("cannot assert type: PrivateKey is not of type *ecdsa.PrivateKey")
	}

	gasLimit := uint64(50000)
	//gasPrice := big.Int()
	gasPrice, err := client.SuggestGasPrice(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("gasLimit(%v)  gasPrice(%v)", gasLimit, gasPrice)

	//Todo... Upgrade Smart Contract. 2022-01-14
	//Smart Contract is just beginning to development.
	//therefore currently using the core API to create a 'DID'. and is have to changed.
	//Requirement: To be consistent, a 'DID' should be created in the smart contract.

	createdDid, createdDoc := dids.CreateDID("eth", pbKeyBase58)
	fmt.Printf("createdDid: %v\n createdDoc: %v\n", createdDid, string(createdDoc))

	auth, err := bind.NewKeyedTransactorWithChainID(pvKeyECDSA, big.NewInt(97))

	transact, err2 := instance.CreateDid(&bind.TransactOpts{
		From:     auth.From,
		Signer:   auth.Signer,
		GasPrice: gasPrice,
		GasLimit: gasLimit,
	}, createdDid, string(createdDoc))
	if err2 != nil {
		fmt.Printf("error() : %v", err2.Error())
	} else {
		fmt.Printf("transact.Hash(): %v\n", transact.Hash())
	}

	return createdDid, nil
}

// ResolveDid - Implements the ResolveDid method from DidMethod
// For this resolve did method, did must be an string
func (m *DidMethodETH) ResolveDid(did string) (string, string, string, error) {
	// Create an IPC based RPC connection to a remote node
	client, err := ethclient.Dial(configs.UseConfig.EthClientUrl)
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	// Instantiate the contract and display its name
	address := common.HexToAddress(configs.UseConfig.EthClientScAddress)
	instance, err := scdid.NewScdid(address, client)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("contract is loaded")
	resolveDid := did
	docs, err := instance.ResolveDid(&bind.CallOpts{Pending: true}, resolveDid)
	fmt.Println(docs)

	didDocument := docs
	DidDocumentMetadata := ""
	ResolutionError := ""

	return didDocument, DidDocumentMetadata, ResolutionError, err
}
