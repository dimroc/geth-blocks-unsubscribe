package main

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"os"
	"os/signal"
	"syscall"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rpc"
)

const (
	EthUrl = "ws://localhost:18546"
)

type EthSubscription interface {
	Err() <-chan error
	Unsubscribe()
}

func main() {
	eth := connectToEth()

	logs := make(chan types.Log)
	done := make(chan struct{})

	go listen(logs, done)
	sub := subscribeToLogs(eth, logs, filterQueryFor(nil))

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs

	sub.Unsubscribe() // When called simultaneously with eth node shutdown, blocks indefinitely.
	close(done)       // In successful repro, this never gets called.
}

func subscribeToLogs(eth *rpc.Client, channel chan<- types.Log, q ethereum.FilterQuery) EthSubscription {
	ctx := context.Background()
	sub, err := eth.EthSubscribe(ctx, channel, "logs", toFilterArg(q))
	if err != nil {
		log.Panic(err)
	}
	return sub
}

func listen(logs chan types.Log, done chan struct{}) {
	for {
		select {
		case <-done:
			fmt.Println("Listening done")
			return
		case l := <-logs:
			fmt.Println("Received log:", l)
		}
	}
}

func connectToEth() *rpc.Client {
	dialed, err := rpc.Dial(EthUrl)
	if err != nil {
		log.Panic("Unable to connect to ", EthUrl)
	}
	return dialed
}

func filterQueryFor(fromBlock *big.Int) ethereum.FilterQuery {
	return ethereum.FilterQuery{
		FromBlock: fromBlock,
	}
}

// toFilterArg filters logs with the given FilterQuery
// https://github.com/ethereum/go-ethereum/blob/762f3a48a00da02fe58063cb6ce8dc2d08821f15/ethclient/ethclient.go#L363
func toFilterArg(q ethereum.FilterQuery) interface{} {
	arg := map[string]interface{}{
		"fromBlock": toBlockNumArg(q.FromBlock),
		"toBlock":   toBlockNumArg(q.ToBlock),
		"address":   []common.Address{},
		"topics":    q.Topics,
	}
	if q.FromBlock == nil {
		arg["fromBlock"] = "0x0"
	}
	return arg
}

func toBlockNumArg(number *big.Int) string {
	if number == nil {
		return "latest"
	}
	return hexutil.EncodeBig(number)
}
