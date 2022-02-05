// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package ERC20Abi

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

// ERC20AbiMetaData contains all meta data concerning the ERC20Abi contract.
var ERC20AbiMetaData = &bind.MetaData{
	ABI: "[{\"constant\":true,\"inputs\":[],\"name\":\"decimals\",\"outputs\":[{\"name\":\"\",\"type\":\"uint8\"}],\"payable\":false,\"type\":\"function\"},{\"constant\":true,\"inputs\":[{\"name\":\"_owner\",\"type\":\"address\"}],\"name\":\"balanceOf\",\"outputs\":[{\"name\":\"balance\",\"type\":\"uint256\"}],\"payable\":false,\"type\":\"function\"}]",
}

// ERC20AbiABI is the input ABI used to generate the binding from.
// Deprecated: Use ERC20AbiMetaData.ABI instead.
var ERC20AbiABI = ERC20AbiMetaData.ABI

// ERC20Abi is an auto generated Go binding around an Ethereum contract.
type ERC20Abi struct {
	ERC20AbiCaller     // Read-only binding to the contract
	ERC20AbiTransactor // Write-only binding to the contract
	ERC20AbiFilterer   // Log filterer for contract events
}

// ERC20AbiCaller is an auto generated read-only Go binding around an Ethereum contract.
type ERC20AbiCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC20AbiTransactor is an auto generated write-only Go binding around an Ethereum contract.
type ERC20AbiTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC20AbiFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type ERC20AbiFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// ERC20AbiSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type ERC20AbiSession struct {
	Contract     *ERC20Abi         // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// ERC20AbiCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type ERC20AbiCallerSession struct {
	Contract *ERC20AbiCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts   // Call options to use throughout this session
}

// ERC20AbiTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type ERC20AbiTransactorSession struct {
	Contract     *ERC20AbiTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts   // Transaction auth options to use throughout this session
}

// ERC20AbiRaw is an auto generated low-level Go binding around an Ethereum contract.
type ERC20AbiRaw struct {
	Contract *ERC20Abi // Generic contract binding to access the raw methods on
}

// ERC20AbiCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type ERC20AbiCallerRaw struct {
	Contract *ERC20AbiCaller // Generic read-only contract binding to access the raw methods on
}

// ERC20AbiTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type ERC20AbiTransactorRaw struct {
	Contract *ERC20AbiTransactor // Generic write-only contract binding to access the raw methods on
}

// NewERC20Abi creates a new instance of ERC20Abi, bound to a specific deployed contract.
func NewERC20Abi(address common.Address, backend bind.ContractBackend) (*ERC20Abi, error) {
	contract, err := bindERC20Abi(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &ERC20Abi{ERC20AbiCaller: ERC20AbiCaller{contract: contract}, ERC20AbiTransactor: ERC20AbiTransactor{contract: contract}, ERC20AbiFilterer: ERC20AbiFilterer{contract: contract}}, nil
}

// NewERC20AbiCaller creates a new read-only instance of ERC20Abi, bound to a specific deployed contract.
func NewERC20AbiCaller(address common.Address, caller bind.ContractCaller) (*ERC20AbiCaller, error) {
	contract, err := bindERC20Abi(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &ERC20AbiCaller{contract: contract}, nil
}

// NewERC20AbiTransactor creates a new write-only instance of ERC20Abi, bound to a specific deployed contract.
func NewERC20AbiTransactor(address common.Address, transactor bind.ContractTransactor) (*ERC20AbiTransactor, error) {
	contract, err := bindERC20Abi(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &ERC20AbiTransactor{contract: contract}, nil
}

// NewERC20AbiFilterer creates a new log filterer instance of ERC20Abi, bound to a specific deployed contract.
func NewERC20AbiFilterer(address common.Address, filterer bind.ContractFilterer) (*ERC20AbiFilterer, error) {
	contract, err := bindERC20Abi(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &ERC20AbiFilterer{contract: contract}, nil
}

// bindERC20Abi binds a generic wrapper to an already deployed contract.
func bindERC20Abi(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := abi.JSON(strings.NewReader(ERC20AbiABI))
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ERC20Abi *ERC20AbiRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ERC20Abi.Contract.ERC20AbiCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ERC20Abi *ERC20AbiRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ERC20Abi.Contract.ERC20AbiTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ERC20Abi *ERC20AbiRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ERC20Abi.Contract.ERC20AbiTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_ERC20Abi *ERC20AbiCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _ERC20Abi.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_ERC20Abi *ERC20AbiTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _ERC20Abi.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_ERC20Abi *ERC20AbiTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _ERC20Abi.Contract.contract.Transact(opts, method, params...)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address _owner) returns(uint256 balance)
func (_ERC20Abi *ERC20AbiCaller) BalanceOf(opts *bind.CallOpts, _owner common.Address) (*big.Int, error) {
	var out []interface{}
	err := _ERC20Abi.contract.Call(opts, &out, "balanceOf", _owner)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address _owner) returns(uint256 balance)
func (_ERC20Abi *ERC20AbiSession) BalanceOf(_owner common.Address) (*big.Int, error) {
	return _ERC20Abi.Contract.BalanceOf(&_ERC20Abi.CallOpts, _owner)
}

// BalanceOf is a free data retrieval call binding the contract method 0x70a08231.
//
// Solidity: function balanceOf(address _owner) returns(uint256 balance)
func (_ERC20Abi *ERC20AbiCallerSession) BalanceOf(_owner common.Address) (*big.Int, error) {
	return _ERC20Abi.Contract.BalanceOf(&_ERC20Abi.CallOpts, _owner)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() returns(uint8)
func (_ERC20Abi *ERC20AbiCaller) Decimals(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _ERC20Abi.contract.Call(opts, &out, "decimals")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() returns(uint8)
func (_ERC20Abi *ERC20AbiSession) Decimals() (uint8, error) {
	return _ERC20Abi.Contract.Decimals(&_ERC20Abi.CallOpts)
}

// Decimals is a free data retrieval call binding the contract method 0x313ce567.
//
// Solidity: function decimals() returns(uint8)
func (_ERC20Abi *ERC20AbiCallerSession) Decimals() (uint8, error) {
	return _ERC20Abi.Contract.Decimals(&_ERC20Abi.CallOpts)
}
