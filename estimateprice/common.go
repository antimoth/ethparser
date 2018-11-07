package estimateprice

import (
	"context"
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

func ensureContext(ctx context.Context) context.Context {
	if ctx == nil {
		return context.TODO()
	}
	return ctx
}
