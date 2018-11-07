package estimator

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type MinerData struct {
	Miner       *common.Address
	NumBlocks   *big.Int
	MinGasPrice *big.Int
}

type Probability struct {
	GasPrice *big.Int
	Prob     float64
}
