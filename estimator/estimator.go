// Copyright 2016 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

// Package ethclient provides a client for the Ethereum RPC API.
package estimator

import (
	"context"
	"fmt"
	"math/big"

	"github.com/antimoth/ethparser/log"
	"github.com/antimoth/ethparser/watcher"
	"github.com/antimoth/lvldb"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault("loglevel", "info")
}

var (
	epLogger = log.NewLogger(viper.GetString("loglevel"))
)

type Estimator struct {
	ew  *watcher.Watcher
	ldb *lvldb.LDBDatabase
}

func NewEstimator(wsUrl string, confirms uint64, ldbPath string) (*Estimator, error) {
	epLogger = log.NewLogger(viper.GetString("loglevel"))
	ew, err := watcher.NewWatcher(wsUrl, confirms)
	if err != nil {
		return nil, err
	}

	ldb, err := lvldb.NewLDBDatabase(viper.GetString("LEVELDB.estimate_store"), 16, 16)
	if err != nil {
		return nil, err
	}

	return &Estimator{
		ew:  ew,
		ldb: ldb,
	}, nil
}

func (et *Estimator) Close() {
	et.ew.Close()
}

func (et *Estimator) Run(startH *big.Int) {
	heightCh := make(chan *big.Int, 1000)
	et.ew.WatchBlock(startH, heightCh)

	txCh := make(chan *common.Hash, 1000)
	et.ew.WatchPendingTx(txCh)

	for {
		select {
		case msg := <-txCh:
			logger.Info("receive tranx", "txhash", msg.Hex())
		case msg := <-heightCh:
			logger.Info("receive block", "height", msg.Uint64())
		}
	}
}
