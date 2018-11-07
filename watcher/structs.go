package watcher

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
)

type Watcher struct {
	c             *rpc.Client
	confirmHeight *big.Int
	currentHeight *big.Int
}

type RpcHeader struct {
	Number    *hexutil.Big  `json:"number"`
	Timestamp *hexutil.Uint `json:"timestamp"`
}
