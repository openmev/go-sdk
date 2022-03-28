// SPDX-Liceense-Identifier: Apache-2.0
package openmev

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"strconv"
	s "strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
)

// @title OpenMeV
// @version 0.1.0
// @description OpenMeV is a Go SDK for the EVM

// FloatToBigInt() can be used to convert float to big int
func FloatToBigInt(val float64) *big.Int {
	bigval := new(big.Float)
	bigval.SetFloat64(val)

	// Set precision if required.
	// bigval.SetPrec(64)

	coin := new(big.Float)
	coin.SetInt(big.NewInt(1000000000000000000))

	bigval.Mul(bigval, coin)

	result := new(big.Int)

	// store converted number in result
	bigval.Int(result)

	return result
}

// CreateRawTransaction creates a raw transaction to be signed
func createRawTransaction(
	chainID *big.Int,
	nonce uint64,
	value float64,
	to string,
	gasLimit string,
	data string,
	txType int,
	gasPrice string,
	maxPriorityFeePerGas string,
	maxFeePerGas string,
	targetTx map[string]interface{},
) *types.Transaction {
	_chainID := chainID
	_nonce := nonce
	_value := FloatToBigInt(value)
	_to := common.HexToAddress(to)
	_gasLimit, _ := strconv.ParseInt(gasLimit, 10, 64)
	_data_tmp := data
	if _data_tmp != "0x" {
		_data_tmp = s.Replace(_data_tmp, "0x", "", 1)
	}
	_data, _ := hex.DecodeString(_data_tmp)

	var _strTemp string
	var _int64Temp int64
	var tx *types.Transaction

	if txType == 0 {
		var _gasPrice *big.Int

		if gasPrice == "auto" {
			_strTemp = targetTx["gasPrice"].(string)
			_strTemp = s.Replace(_strTemp, "0x", "", 1)
			_int64Temp, _ = strconv.ParseInt(_strTemp, 16, 64)
			_gasPrice = big.NewInt(_int64Temp)
		} else {
			_int64Temp, _ = strconv.ParseInt(gasPrice, 10, 64)
			_gasPrice = big.NewInt(_int64Temp)
			_gasPrice = _gasPrice.Mul(_gasPrice, big.NewInt(1000000000))
		}

		tx = types.NewTx(&types.LegacyTx{
			Nonce:    _nonce,
			GasPrice: _gasPrice,
			Gas:      uint64(_gasLimit),
			To:       &_to,
			Value:    _value,
			Data:     _data,
		})
	} else if txType == 2 {
		var _maxPriorityFeePerGas, _maxFeePerGas *big.Int

		if maxPriorityFeePerGas == "auto" {
			_strTemp = targetTx["maxPriorityFeePerGas"].(string)
			_strTemp = s.Replace(_strTemp, "0x", "", 1)
			_int64Temp, _ = strconv.ParseInt(_strTemp, 16, 64)
			_maxPriorityFeePerGas = big.NewInt(_int64Temp)
		} else {
			_int64Temp, _ = strconv.ParseInt(maxPriorityFeePerGas, 10, 64)
			_maxPriorityFeePerGas = big.NewInt(_int64Temp)
			_maxPriorityFeePerGas = _maxPriorityFeePerGas.Mul(_maxPriorityFeePerGas, big.NewInt(1000000000))
		}

		if maxFeePerGas == "auto" {
			_strTemp = targetTx["maxFeePerGas"].(string)
			_strTemp = s.Replace(_strTemp, "0x", "", 1)
			_int64Temp, _ = strconv.ParseInt(_strTemp, 16, 64)
			_maxFeePerGas = big.NewInt(_int64Temp)
		} else {
			_int64Temp, _ = strconv.ParseInt(maxFeePerGas, 10, 64)
			_maxFeePerGas = big.NewInt(_int64Temp)
			_maxFeePerGas = _maxFeePerGas.Mul(_maxFeePerGas, big.NewInt(1000000000))
		}

		tx = types.NewTx(&types.DynamicFeeTx{
			ChainID:   _chainID,
			Nonce:     _nonce,
			GasTipCap: _maxPriorityFeePerGas,
			GasFeeCap: _maxFeePerGas,
			Gas:       uint64(_gasLimit),
			To:        &_to,
			Value:     _value,
			Data:      _data,
		})
	}

	return tx
}

// signRawTransaction() can be used to sign the transaction
func signRawTransaction(_pvKey string, txType int, tx *types.Transaction) *types.Transaction {
	_pvKey = s.Replace(_pvKey, "0x", "", 1)
	pvKey, err := crypto.HexToECDSA(_pvKey)
	if err != nil {
		fmt.Println("pvKey_Err:", err)
	}

	var signedTx *types.Transaction
	if txType == 0 {
		signedTx, err = types.SignTx(tx, types.LatestSignerForChainID(nil), pvKey)
	} else if txType == 2 {
		signedTx, err = types.SignTx(tx, types.LatestSignerForChainID(tx.ChainId()), pvKey)
	}

	if err != nil {
		fmt.Println("SignNewTx_Err:", err.Error())
	}

	return signedTx
}
