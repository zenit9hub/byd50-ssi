package main

import (
	"byd50-ssi/pkg/did/c-shared/foo"
	"byd50-ssi/pkg/did/core"
	"byd50-ssi/pkg/did/core/dids"
	"byd50-ssi/pkg/did/core/driver/scdid"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
	"time"
)

func geth() {
	ks := keystore.NewKeyStore("/path/to/keystore", keystore.StandardScryptN, keystore.StandardScryptP)
	am := accounts.NewManager(&accounts.Config{InsecureUnlockAllowed: false}, ks)
	log.Printf("am=%v", am)

	// Create a new account with the specified encryption passphrase.
	newAcc, _ := ks.NewAccount("Creation password")
	fmt.Println(newAcc)

	// Export the newly created account with a different passphrase. The returned
	// data from this method invocation is a JSON encoded, encrypted key-file.
	jsonAcc, _ := ks.Export(newAcc, "Creation password", "Export password")

	// Update the passphrase on the account created above inside the local keystore.
	_ = ks.Update(newAcc, "Creation password", "Update password")

	// Delete the account updated above from the local keystore.
	_ = ks.Delete(newAcc, "Update password")

	// Import back the account we've exported (and then deleted) above with yet
	// again a fresh passphrase.
	impAcc, _ := ks.Import(jsonAcc, "Export password", "Import password")
	log.Printf("impAcc=%v", impAcc)

	//Signing from Go
	// Create a new account to sign transactions with
	signer, _ := ks.NewAccount("Signer password")
	txHash := common.HexToHash("0x0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")

	// Sign a transaction with a single authorization
	signature, _ := ks.SignHashWithPassphrase(signer, "Signer password", txHash.Bytes())
	log.Printf("signature=%v", signature)

	// Sign a transaction with multiple manually cancelled authorizations
	_ = ks.Unlock(signer, "Signer password")
	signature, _ = ks.SignHash(signer, txHash.Bytes())
	_ = ks.Lock(signer.Address)
	log.Printf("signature=%v", signature)

	// Sign a transaction with multiple automatically cancelled authorizations
	_ = ks.TimedUnlock(signer, "Signer password", time.Second)
	signature, _ = ks.SignHash(signer, txHash.Bytes())

	log.Printf("signature=%v", signature)
}

func main() {
	// Create an IPC based RPC connection to a remote node
	client, err := ethclient.Dial("http://3.37.125.54:8545")
	if err != nil {
		log.Fatalf("Failed to connect to the Ethereum client: %v", err)
	}

	vpJwtSample := "eyJhbGciOiJFUzI1NiIsImtpZCI6ImRpZDpyZWFwOjViZTE1YmIwOTVmNDU3NThiZjY4MTEwODEzNjA4NTBjNmZiY2NjZWUiLCJ0eXAiOiJKV1QifQ.eyJub25jZSI6IkJ2MGpEbEdMV1FwTiIsInZjIjp7IkBjb250ZXh0IjpbImh0dHBzOi8vd3d3LnczLm9yZy8yMDE4L2NyZWRlbnRpYWxzL3YxIiwiaHR0cHM6Ly93d3cudzMub3JnLzIwMTgvY3JlZGVudGlhbHMvZXhhbXBsZXMvdjEiXSwiY3JlZGVudGlhbFN1YmplY3QiOnsiYmlydGgiOiIyMDAwLTExLTA4IiwiY291bnRyeSI6IlMuS29yZWEiLCJuYW1lIjoiSG9uZyBHaWwtRG9uZyJ9LCJ0eXBlIjpbIlZlcmlmaWFibGVDcmVkZW50aWFsIiwiZUlkQ2FyZENyZWRlbnRpYWwiXX0sImV4cCI6MTY0MTI3NzYyMywianRpIjoiMDg5YTQxMWYtMGQ4OC00NTBmLThjYzAtMWEzYWNmZWJlY2QzIiwiaWF0IjoxNjQxMjc3NDQzLCJpc3MiOiJodHRwOi8vd3d3Lmdvdi5rci9yZXNpZGVudHJlZ2lzdHJhdGlvbiIsIm5iZiI6MTY0MTI3NzQ0Mywic3ViIjoiZGlkOnJlYXA6ZDM4MmExNTFiZjQ2NzlmMzQ4YjA3M2NkNWJlODg2MTRiYzIyOTJmZCJ9.nZVvh9wv7XYNTNl2VrQlD3qaUvZ2ljGzPl6iKypPU_WrRn6bH2T6y_WyHbjzzZTahwrQqO3DU0buq02l5NF_fA"
	exp := foo.ClaimsGetInt64(vpJwtSample, "exp")
	iat := foo.ClaimsGetInt64(vpJwtSample, "iat")
	log.Printf("~~~~~~~~~~>    exp : %v  iat : %v", exp, iat)

	// Instantiate the contract and display its name
	address := common.HexToAddress("0x928B45C153fF4df3220C79D30a6b490D7a1e4878")
	instance, err := scdid.NewScdid(address, client)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("contract is loaded")
	_ = instance
	resolveDid := "did:byd50:12345555"
	docs, err := instance.ResolveDid(&bind.CallOpts{Pending: true}, resolveDid)
	fmt.Println(docs)

	// Generate a new random account and a funded simulator
	newKey, _ := crypto.GenerateKey()
	bytes := crypto.FromECDSA(newKey)
	newKeyPbKey := newKey.Public()
	newKeyPbKeyECDSA, _ := newKeyPbKey.(*ecdsa.PublicKey)
	newKeyAddr := crypto.PubkeyToAddress(*newKeyPbKeyECDSA).Hex()
	fmt.Printf("bytes: %v\n newKeyAddr: %v\n", bytes, newKeyAddr)

	myDKMS, err := core.InitDKMS(core.KeyTypeECDSA)
	myPvKey := myDKMS.PvKey().(*ecdsa.PrivateKey)
	fmt.Printf("pvkey: %v\n", myPvKey)
	myBytes := crypto.FromECDSA(myPvKey)
	fmt.Printf("myBytes: %v\n", myBytes)

	pvKeyBytes4import := []byte{144, 228, 44, 70, 165, 36, 200, 74, 21, 141, 93, 62, 93, 191, 5, 52, 115, 225, 36, 88, 248, 189, 54, 105, 44, 106, 45, 246, 239, 192, 156, 50}
	importedPvKey, err := crypto.ToECDSA(pvKeyBytes4import)
	importedPvKeyBytes := crypto.FromECDSA(importedPvKey)
	pbKey := importedPvKey.Public()
	pbKeyECDSA, ok := pbKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("cannot assert type: publicKey is not of type *ecdsa.PublicKey")
	}
	importedAddr := crypto.PubkeyToAddress(*pbKeyECDSA).Hex()
	fmt.Printf("importedpvKeyBytes: %v\n importedAddr: %v\n", importedPvKeyBytes, importedAddr)

	auth, err := bind.NewKeyedTransactorWithChainID(importedPvKey, big.NewInt(3333))

	fmt.Printf("auth.From : %v\n", auth.From)
	createdDid, createdDoc := dids.CreateDID("eth", myDKMS.PbKeyBase58())
	fmt.Printf("createdDid: %v\n createdDoc: %v\n", createdDid, string(createdDoc))

	transact, err2 := instance.CreateDid(&bind.TransactOpts{
		From:     auth.From,
		Signer:   auth.Signer,
		GasPrice: big.NewInt(0),
		GasLimit: 3141592,
	}, createdDid, string(createdDoc))
	if err2 != nil {
		fmt.Printf("error() : %v", err2.Error())
	} else {
		fmt.Printf("transact.Hash(): %v\n", transact.Hash())
	}

}
