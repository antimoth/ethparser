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
package watcher

import (
	"context"
	"fmt"
	"math/big"

	"ethparser/log"

	"github.com/antimoth/ethparser/client"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/spf13/viper"
)

func init() {
	viper.SetDefault("loglevel", "info")
}

var (
	eLogger = log.NewLogger(viper.GetString("loglevel"))
)

type Watcher struct {
	c             *client.Client
	confirmHeight *big.Int
	currentHeight *big.Int
}

func NewWatcher(rawurl string, confirmHeight uint64) (*Watcher, error) {

	eLogger = log.NewLogger(viper.GetString("loglevel"))

	c, err := client.NewClient(rawurl)
	if err != nil {
		return nil, err
	}

	if confirmHeight < uint64(1) {
		confirmHeight = uint64(1)
	}

	return &Watcher{
		c:             c,
		confirmHeight: new(big.Int).SetUint64(confirmHeight),
		currentHeight: new(big.Int),
	}, nil
}

func (ew *Watcher) Close() {
	ew.c.Close()
}

func (ew *Watcher) ReviewBlock(start big.Int, ch chan<- *big.Int) {
	cur, err := ew.c.BlockNumber()
	if err != nil {
		panic("get current height error!")
	}
	ew.currentHeight = cur

	curConfirm := new(big.Int).Sub(cur, ew.confirmHeight)
	increaser := big.NewInt(1)
	go func() {
		for i := &start; i.Cmp(curConfirm) <= 0; i = new(big.Int).Add(i, increaser) {
			ch <- i
		}
		close(ch)
	}()
}

func (ew *Watcher) StartWatchBlock(start big.Int, heightCh chan<- *big.Int) {
	wCh := make(chan *client.RpcHeader, 1000)
	rCh := make(chan *big.Int, 1000)
	sub, err := ew.SubscribeNewHead(client.EnsureContext(nil), wCh)
	if err != nil {
		panic(fmt.Sprintf("subscribe new block error! e is %v!", err.Error()))
	}
	ew.ReviewBlock(start, rCh)
	bigConfirmH := ew.confirmHeight
	go func() {
		defer sub.Unsubscribe()
		for {
			select {
			case blockHeight := <-rCh:
				heightCh <- blockHeight
			default:
				select {
				case blockHeight := <-rCh:
					heightCh <- blockHeight
				case blockHeader := <-wCh:
					bigIntNumber := (*big.Int)(blockHeader.Number)
					if bigIntNumber.Cmp(ew.currentHeight) > 0 {
						heightCh <- new(big.Int).Sub(bigIntNumber, bigConfirmH)
						ew.currentHeight = bigIntNumber
					}
				case err := <-sub.Err():
					eLogger.Error("watch block error", "error", err.Error())
					return
				}
			}
		}
	}()
}

func (ew *Watcher) WatchPendingTx(ch chan<- *common.Hash) {
	txCh := make(chan *common.Hash, 1000)

	ctx := client.EnsureContext(nil)
	sub, err := ew.SubscribePendingTx(ctx, txCh)
	if err != nil {
		panic(err)
	}

	go func() {
		defer sub.Unsubscribe()
		for {
			select {
			case txHash := <-txCh:
				ch <- txHash
			case err := <-sub.Err():
				eLogger.Error("sub pending tranx error!", "error", err.Error())
				return
			}
		}
	}()
}

// SubscribeNewHead subscribes to notifications about the current blockchain head
// on the given channel.
func (ew *Watcher) SubscribeNewHead(ctx context.Context, ch chan<- *client.RpcHeader) (ethereum.Subscription, error) {
	return ew.c.EthSubscribe(ctx, ch, "newHeads")
}

func (ew *Watcher) SubscribePendingTx(ctx context.Context, ch chan<- *common.Hash) (ethereum.Subscription, error) {
	return ew.c.EthSubscribe(ctx, ch, "newPendingTransactions")
}
