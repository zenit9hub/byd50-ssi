# On-chain Registry Deployment Guide (EVM Testnet)

This guide describes a minimal path to deploy a DID Registry contract and connect it to the `geth_client` sample.

## 1) Choose a Testnet
Recommended for fast confirmations:
- Polygon Amoy (chainId: 80002)
- Avalanche Fuji (chainId: 43113)
- Base Sepolia (chainId: 84532)
- Ethereum Sepolia (chainId: 11155111)

## 2) Contract Source (Minimal Example)
Use a minimal registry contract if you do not already have one.

```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

contract DIDRegistry {
    mapping(string => string) private docs;

    function CreateDid(string memory _did, string memory _document) public {
        docs[_did] = _document;
    }

    function ResolveDid(string memory _did) public view returns (string memory) {
        return docs[_did];
    }
}
```

## 3) Deploy (Remix or Hardhat)
### Remix (fastest for training)
1. Open Remix IDE and paste the contract.
2. Compile with Solidity 0.8.x.
3. Deploy to your selected testnet.
4. Record the deployed contract address.

### Hardhat (repeatable)
1. Initialize Hardhat project.
2. Add network RPC and private key in `.env`.
3. Deploy and record contract address.

## 4) Wire to the Project
Set environment variables for `geth_client`:

```
ETH_RPC_URL=https://YOUR_TESTNET_RPC
ETH_CHAIN_ID=YOUR_CHAIN_ID
ETH_REGISTRY_ADDRESS=0xYourDeployedContract
ETH_PRIVATE_KEY_HEX=your_hex_private_key
ETH_RESOLVE_DID=did:byd50:12345555
```

## 5) Run the Sample
```
go run ./apps/geth_client/main.go
```

Expected:
- ResolveDid prints existing document (if any)
- CreateDid sends a transaction and prints tx hash

## 6) Notes
- This repo contains ABI bindings in `pkg/did/core/driver/scdid`. If you use a different contract,
  regenerate bindings and update the client.
- For training, keep this as an optional lab to avoid network issues blocking core SSI exercises.
