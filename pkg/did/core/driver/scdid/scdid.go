// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package scdid

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
)

// ScdidMetaData contains all meta data concerning the Scdid contract.
var ScdidMetaData = &bind.MetaData{
	ABI: "[{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_did\",\"type\":\"string\"},{\"internalType\":\"string\",\"name\":\"_document\",\"type\":\"string\"}],\"name\":\"CreateDid\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"GetHelloWorld\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"_did\",\"type\":\"string\"}],\"name\":\"ResolveDid\",\"outputs\":[{\"internalType\":\"string\",\"name\":\"\",\"type\":\"string\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"string\",\"name\":\"str\",\"type\":\"string\"}],\"name\":\"SetHelloWorld\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]",
}

// ScdidABI is the input ABI used to generate the binding from.
// Deprecated: Use ScdidMetaData.ABI instead.
var ScdidABI = ScdidMetaData.ABI

// Scdid is an auto generated Go binding around an Ethereum contract.
type Scdid struct {
	ScdidCaller     // Read-only binding to the contract
	ScdidTransactor // Write-only binding to the contract
	ScdidFilterer   // Log filterer for contract events
}

// ScdidCaller is an auto generated read-only Go binding around an Ethereum contract.
type ScdidCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ScdidTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ScdidTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ScdidFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ScdidFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ScdidSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ScdidSession struct {
	Contract     *Scdid            // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ScdidCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ScdidCallerSession struct {
	Contract *ScdidCaller  // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts // Call options to use throughout this session
}

// ScdidTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ScdidTransactorSession struct {
	Contract     *ScdidTransactor  // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ScdidRaw is an auto generated low-level Go binding around an Ethereum contract.
type ScdidRaw struct {
	Contract *Scdid // Generic contract binding to access the raw methods on
}

// ScdidCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ScdidCallerRaw struct {
	Contract *ScdidCaller // Generic read-only contract binding to access the raw methods on
}

// ScdidTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ScdidTransactorRaw struct {
	Contract *ScdidTransactor // Generic write-only contract binding to access the raw methods on
}

// NewScdid creates a new instance of Scdid, bound to a specific deployed contract.
func NewScdid(address common.Address, backend bind.ContractBackend) (*Scdid, error) {
	contract, err := bindScdid(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Scdid{ScdidCaller: ScdidCaller{contract: contract}, ScdidTransactor: ScdidTransactor{contract: contract}, ScdidFilterer: ScdidFilterer{contract: contract}}, nil
}

// NewScdidCaller creates a new read-only instance of Scdid, bound to a specific deployed contract.
func NewScdidCaller(address common.Address, caller bind.ContractCaller) (*ScdidCaller, error) {
	contract, err := bindScdid(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ScdidCaller{contract: contract}, nil
}

// NewScdidTransactor creates a new write-only instance of Scdid, bound to a specific deployed contract.
func NewScdidTransactor(address common.Address, transactor bind.ContractTransactor) (*ScdidTransactor, error) {
	contract, err := bindScdid(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ScdidTransactor{contract: contract}, nil
}

// NewScdidFilterer creates a new log filterer instance of Scdid, bound to a specific deployed contract.
func NewScdidFilterer(address common.Address, filterer bind.ContractFilterer) (*ScdidFilterer, error) {
	contract, err := bindScdid(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ScdidFilterer{contract: contract}, nil
}

// bindScdid binds a generic wrapper to an already deployed contract.
func bindScdid(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ScdidABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Scdid *ScdidRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Scdid.Contract.ScdidCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Scdid *ScdidRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Scdid.Contract.ScdidTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Scdid *ScdidRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Scdid.Contract.ScdidTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Scdid *ScdidCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Scdid.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Scdid *ScdidTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Scdid.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Scdid *ScdidTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Scdid.Contract.contract.Transact(opts, method, params...)
}

// GetHelloWorld is a free data retrieval call binding the contract method 0xa1aa8ab1.
//
// Solidity: function GetHelloWorld() view returns(string)
func (_Scdid *ScdidCaller) GetHelloWorld(opts *bind.CallOpts) (string, error) {
	var out []interface{}
	err := _Scdid.contract.Call(opts, &out, "GetHelloWorld")

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// GetHelloWorld is a free data retrieval call binding the contract method 0xa1aa8ab1.
//
// Solidity: function GetHelloWorld() view returns(string)
func (_Scdid *ScdidSession) GetHelloWorld() (string, error) {
	return _Scdid.Contract.GetHelloWorld(&_Scdid.CallOpts)
}

// GetHelloWorld is a free data retrieval call binding the contract method 0xa1aa8ab1.
//
// Solidity: function GetHelloWorld() view returns(string)
func (_Scdid *ScdidCallerSession) GetHelloWorld() (string, error) {
	return _Scdid.Contract.GetHelloWorld(&_Scdid.CallOpts)
}

// ResolveDid is a free data retrieval call binding the contract method 0x15f0f802.
//
// Solidity: function ResolveDid(string _did) view returns(string)
func (_Scdid *ScdidCaller) ResolveDid(opts *bind.CallOpts, _did string) (string, error) {
	var out []interface{}
	err := _Scdid.contract.Call(opts, &out, "ResolveDid", _did)

	if err != nil {
		return *new(string), err
	}

	out0 := *abi.ConvertType(out[0], new(string)).(*string)

	return out0, err

}

// ResolveDid is a free data retrieval call binding the contract method 0x15f0f802.
//
// Solidity: function ResolveDid(string _did) view returns(string)
func (_Scdid *ScdidSession) ResolveDid(_did string) (string, error) {
	return _Scdid.Contract.ResolveDid(&_Scdid.CallOpts, _did)
}

// ResolveDid is a free data retrieval call binding the contract method 0x15f0f802.
//
// Solidity: function ResolveDid(string _did) view returns(string)
func (_Scdid *ScdidCallerSession) ResolveDid(_did string) (string, error) {
	return _Scdid.Contract.ResolveDid(&_Scdid.CallOpts, _did)
}

// CreateDid is a paid mutator transaction binding the contract method 0x77aeac44.
//
// Solidity: function CreateDid(string _did, string _document) returns()
func (_Scdid *ScdidTransactor) CreateDid(opts *bind.TransactOpts, _did string, _document string) (*types.Transaction, error) {
	return _Scdid.contract.Transact(opts, "CreateDid", _did, _document)
}

// CreateDid is a paid mutator transaction binding the contract method 0x77aeac44.
//
// Solidity: function CreateDid(string _did, string _document) returns()
func (_Scdid *ScdidSession) CreateDid(_did string, _document string) (*types.Transaction, error) {
	return _Scdid.Contract.CreateDid(&_Scdid.TransactOpts, _did, _document)
}

// CreateDid is a paid mutator transaction binding the contract method 0x77aeac44.
//
// Solidity: function CreateDid(string _did, string _document) returns()
func (_Scdid *ScdidTransactorSession) CreateDid(_did string, _document string) (*types.Transaction, error) {
	return _Scdid.Contract.CreateDid(&_Scdid.TransactOpts, _did, _document)
}

// SetHelloWorld is a paid mutator transaction binding the contract method 0x91c094ac.
//
// Solidity: function SetHelloWorld(string str) returns()
func (_Scdid *ScdidTransactor) SetHelloWorld(opts *bind.TransactOpts, str string) (*types.Transaction, error) {
	return _Scdid.contract.Transact(opts, "SetHelloWorld", str)
}

// SetHelloWorld is a paid mutator transaction binding the contract method 0x91c094ac.
//
// Solidity: function SetHelloWorld(string str) returns()
func (_Scdid *ScdidSession) SetHelloWorld(str string) (*types.Transaction, error) {
	return _Scdid.Contract.SetHelloWorld(&_Scdid.TransactOpts, str)
}

// SetHelloWorld is a paid mutator transaction binding the contract method 0x91c094ac.
//
// Solidity: function SetHelloWorld(string str) returns()
func (_Scdid *ScdidTransactorSession) SetHelloWorld(str string) (*types.Transaction, error) {
	return _Scdid.Contract.SetHelloWorld(&_Scdid.TransactOpts, str)
}
