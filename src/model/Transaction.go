package model

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type Transaction struct {
	Type       uint8            `json:"type"`
	ChainId    *big.Int         `json:"chain_id"`
	Data       []byte           `json:"data"`        // Data returns the input data of the transaction.
	AccessList types.AccessList `json:"access_list"` // AccessList returns the access list of the transaction.
	Gas        uint64           `json:"gas"`
	GasPrice   *big.Int         `json:"gas_price"`
	GasTipCap  *big.Int         `json:"gas_tip_cap"`
	GasFeeCap  *big.Int         `json:"gas_fee_cap"`
	Value      *big.Int         `json:"value"`
	Nonce      uint64           `json:"nonce"`
	To         *common.Address  `json:"to"`
	Cost       *big.Int         `json:"cost"` // Cost returns gas * gasPrice + value.

	// RawSignatureValues returns the V, R, S signature values of the transaction.
	V *big.Int `json:"v"`
	R *big.Int `json:"r"`
	S *big.Int `json:"s"`

	Hash common.Hash        `json:"hash"`
	Size common.StorageSize `json:"size"`
}

type Transactions []*Transaction

type Block struct {
	Number       uint64 `json:"number"`
	Transactions `json:"transactions"`
}

type Blocks []*Block
