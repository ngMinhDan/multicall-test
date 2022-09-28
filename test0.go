package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"log"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

const bytecode = "0x608060405234801561001057600080fd5b50604051610915380380610915833981810160405281019061003291906105a0565b6000825190508082511461007b576040517f08c379a00000000000000000000000000000000000000000000000000000000081526004016100729061069b565b60405180910390fd5b60008167ffffffffffffffff8111156100975761009661024e565b5b6040519080825280602002602001820160405280156100ca57816020015b60608152602001906001900390816100b55790505b50905060005b828110156101f45760008582815181106100ed576100ec6106bb565b5b60200260200101519050600085838151811061010c5761010b6106bb565b5b602002602001015190506000808373ffffffffffffffffffffffffffffffffffffffff168360405161013e9190610731565b6000604051808303816000865af19150503d806000811461017b576040519150601f19603f3d011682016040523d82523d6000602084013e610180565b606091505b5091509150816101bd57604051806020016040528060008152508686815181106101ad576101ac6106bb565b5b60200260200101819052506101dd565b808686815181106101d1576101d06106bb565b5b60200260200101819052505b5050505080806101ec90610781565b9150506100d0565b506000438260405160200161020a9291906108e4565b604051602081830303815290604052905080516101008201f35b6000604051905090565b600080fd5b600080fd5b600080fd5b6000601f19601f8301169050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052604160045260246000fd5b6102868261023d565b810181811067ffffffffffffffff821117156102a5576102a461024e565b5b80604052505050565b60006102b8610224565b90506102c4828261027d565b919050565b600067ffffffffffffffff8211156102e4576102e361024e565b5b602082029050602081019050919050565b600080fd5b600073ffffffffffffffffffffffffffffffffffffffff82169050919050565b6000610325826102fa565b9050919050565b6103358161031a565b811461034057600080fd5b50565b6000815190506103528161032c565b92915050565b600061036b610366846102c9565b6102ae565b9050808382526020820190506020840283018581111561038e5761038d6102f5565b5b835b818110156103b757806103a38882610343565b845260208401935050602081019050610390565b5050509392505050565b600082601f8301126103d6576103d5610238565b5b81516103e6848260208601610358565b91505092915050565b600067ffffffffffffffff82111561040a5761040961024e565b5b602082029050602081019050919050565b600080fd5b600067ffffffffffffffff82111561043b5761043a61024e565b5b6104448261023d565b9050602081019050919050565b60005b8381101561046f578082015181840152602081019050610454565b60008484015250505050565b600061048e61048984610420565b6102ae565b9050828152602081018484840111156104aa576104a961041b565b5b6104b5848285610451565b509392505050565b600082601f8301126104d2576104d1610238565b5b81516104e284826020860161047b565b91505092915050565b60006104fe6104f9846103ef565b6102ae565b90508083825260208201905060208402830185811115610521576105206102f5565b5b835b8181101561056857805167ffffffffffffffff81111561054657610545610238565b5b80860161055389826104bd565b85526020850194505050602081019050610523565b5050509392505050565b600082601f83011261058757610586610238565b5b81516105978482602086016104eb565b91505092915050565b600080604083850312156105b7576105b661022e565b5b600083015167ffffffffffffffff8111156105d5576105d4610233565b5b6105e1858286016103c1565b925050602083015167ffffffffffffffff81111561060257610601610233565b5b61060e85828601610572565b9150509250929050565b600082825260208201905092915050565b7f4572726f723a204172726179206c656e6774687320646f206e6f74206d61746360008201527f682e000000000000000000000000000000000000000000000000000000000000602082015250565b6000610685602283610618565b915061069082610629565b604082019050919050565b600060208201905081810360008301526106b481610678565b9050919050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052603260045260246000fd5b600081519050919050565b600081905092915050565b600061070b826106ea565b61071581856106f5565b9350610725818560208601610451565b80840191505092915050565b600061073d8284610700565b915081905092915050565b7f4e487b7100000000000000000000000000000000000000000000000000000000600052601160045260246000fd5b6000819050919050565b600061078c82610777565b91507fffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff82036107be576107bd610748565b5b600182019050919050565b6107d281610777565b82525050565b600081519050919050565b600082825260208201905092915050565b6000819050602082019050919050565b600082825260208201905092915050565b6000610820826106ea565b61082a8185610804565b935061083a818560208601610451565b6108438161023d565b840191505092915050565b600061085a8383610815565b905092915050565b6000602082019050919050565b600061087a826107d8565b61088481856107e3565b935083602082028501610896856107f4565b8060005b858110156108d257848403895281516108b3858261084e565b94506108be83610862565b925060208a0199505060018101905061089a565b50829750879550505050505092915050565b60006040820190506108f960008301856107c9565b818103602083015261090b818461086f565b9050939250505056fe"
const MulticallABI = "[{\"inputs\":[{\"internalType\":\"address[]\",\"name\":\"targets\",\"type\":\"address[]\"},{\"internalType\":\"bytes[]\",\"name\":\"datas\",\"type\":\"bytes[]\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"}]"
const _aaveProtocolDataProviderAddress = "0x057835ad21a177dbdd3090bb1cae03eacf78fc6d"
const aaveProtocolDataProviderABI = "[{\"inputs\":[{\"internalType\":\"contractILendingPoolAddressesProvider\",\"name\":\"addressesProvider\",\"type\":\"address\"}],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[],\"name\":\"ADDRESSES_PROVIDER\",\"outputs\":[{\"internalType\":\"contractILendingPoolAddressesProvider\",\"name\":\"\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllATokens\",\"outputs\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"symbol\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"}],\"internalType\":\"structAaveProtocolDataProvider.TokenData[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[],\"name\":\"getAllReservesTokens\",\"outputs\":[{\"components\":[{\"internalType\":\"string\",\"name\":\"symbol\",\"type\":\"string\"},{\"internalType\":\"address\",\"name\":\"tokenAddress\",\"type\":\"address\"}],\"internalType\":\"structAaveProtocolDataProvider.TokenData[]\",\"name\":\"\",\"type\":\"tuple[]\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"}],\"name\":\"getReserveConfigurationData\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"decimals\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"ltv\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"liquidationThreshold\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"liquidationBonus\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"reserveFactor\",\"type\":\"uint256\"},{\"internalType\":\"bool\",\"name\":\"usageAsCollateralEnabled\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"borrowingEnabled\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"stableBorrowRateEnabled\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"isActive\",\"type\":\"bool\"},{\"internalType\":\"bool\",\"name\":\"isFrozen\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"}],\"name\":\"getReserveData\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"availableLiquidity\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"totalStableDebt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"totalVariableDebt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"liquidityRate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"variableBorrowRate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"stableBorrowRate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"averageStableBorrowRate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"liquidityIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"variableBorrowIndex\",\"type\":\"uint256\"},{\"internalType\":\"uint40\",\"name\":\"lastUpdateTimestamp\",\"type\":\"uint40\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"}],\"name\":\"getReserveTokensAddresses\",\"outputs\":[{\"internalType\":\"address\",\"name\":\"aTokenAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"stableDebtTokenAddress\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"variableDebtTokenAddress\",\"type\":\"address\"}],\"stateMutability\":\"view\",\"type\":\"function\"},{\"inputs\":[{\"internalType\":\"address\",\"name\":\"asset\",\"type\":\"address\"},{\"internalType\":\"address\",\"name\":\"user\",\"type\":\"address\"}],\"name\":\"getUserReserveData\",\"outputs\":[{\"internalType\":\"uint256\",\"name\":\"currentATokenBalance\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"currentStableDebt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"currentVariableDebt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"principalStableDebt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"scaledVariableDebt\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"stableBorrowRate\",\"type\":\"uint256\"},{\"internalType\":\"uint256\",\"name\":\"liquidityRate\",\"type\":\"uint256\"},{\"internalType\":\"uint40\",\"name\":\"stableRateLastUpdated\",\"type\":\"uint40\"},{\"internalType\":\"bool\",\"name\":\"usageAsCollateralEnabled\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"

func main() {
	client, err := ethclient.Dial("https://mainnet.infura.io/v3/ce46e74b5e0c491ab6c217987f0220c6")
	if err != nil {
		panic(err)
	} else {
		fmt.Println("Connecting to node ...")
	}

	msg := ethereum.CallMsg{}
	bytecodeBytes, err := hex.DecodeString(strings.ReplaceAll(bytecode, "0x", ""))
	if err != nil {
		panic(err)
	}

	aaveProtocolDataProviderABI, err := abi.JSON(strings.NewReader(string(aaveProtocolDataProviderABI)))
	if err != nil {
		log.Fatal("abi.JSON(strings.NewReader(string(aaveProtocolDataProviderABIByte))): ", err)
	}

	abiMultiCall, err := abi.JSON(strings.NewReader(string(MulticallABI)))
	if err != nil {
		log.Fatal("abi.JSON(strings.NewReader(string(MulticallABI))): ", err)
	}

	method := aaveProtocolDataProviderABI.Methods["getUserReserveData"]

	id := method.ID
	usdcAddress := common.HexToAddress("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48")
	userAddress := common.HexToAddress("0xC74B707b122fd65c98957b5F30392208ECf74317")
	inputs, err := method.Inputs.Pack(usdcAddress, userAddress)

	if err != nil {
		panic(err)
	}
	inputs = append(id[:], inputs[:]...)

	targets := []common.Address{
		common.HexToAddress(_aaveProtocolDataProviderAddress),
	}

	datas := [][]byte{
		inputs,
	}
	// tạo input data bằng khởi tạo : Inputs.Pack(targets, datas)
	inputData, err := abiMultiCall.Constructor.Inputs.Pack(targets, datas)

	if err != nil {
		panic(err)
	}

	msg.Data = append(bytecodeBytes[:], inputData[:]...)
	result, err := client.PendingCallContract(context.TODO(), msg)
	if err != nil {
		panic(err)
	}

	// count := 0
	// for {
	res, err := method.Outputs.Unpack(result[:])
	fmt.Println(method)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(res)

	// KET QUA MONG MUON
	// ===================================
	// currentATokenBalance   uint256 :  0
	// currentStableDebt   uint256 :  0
	// currentVariableDebt   uint256 :  900414
	// principalStableDebt   uint256 :  0
	// scaledVariableDebt   uint256 :  809287
	// stableBorrowRate   uint256 :  0
	// liquidityRate   uint256 :  4948714111130000762018747
	// stableRateLastUpdated   uint40 :  0
	// usageAsCollateralEnabled   bool :  false
	// ===================================

	// KET QUA DAT DUOC
	// 2022/09/28 23:21:16 [900425 0 809287 0 4938747366840871168343614 0 0 0 false]

}
