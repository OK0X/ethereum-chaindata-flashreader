package main

import (
	"fmt"
	"os/user"

	"github.com/OK0X/ethereum-chaindata-flashreader/src/ledger"
	"github.com/ethereum/go-ethereum/common/prque"
	"github.com/ethereum/go-ethereum/core/rawdb"
)

func main() {

	usr, _ := user.Current()
	benchDataDir := usr.HomeDir + "/geth/data/geth/chaindata"
	db, err := rawdb.NewLevelDBDatabase(benchDataDir, 128, 1024, "", false)

	if err != nil {
		fmt.Errorf("error opening database at %v: %v", benchDataDir, err)
	}
	from := uint64(0)
	to := uint64(8)
	lastNum := to
	blocks := 0
	txsCh := ledger.ReadTransactions(db, from, to, false, nil)

	queue := prque.New(nil)

	for chanDelivery := range txsCh {
		queue.Push(chanDelivery, int64(chanDelivery.Number))
		for !queue.Empty() {
			// If the next available item is gapped, return
			if _, priority := queue.Peek(); priority != int64(lastNum-1) {
				break
			}

			// Next block available, pop it off and index it
			delivery := queue.PopItem().(*ledger.BlockTxs)
			lastNum = delivery.Number
			blocks++
			fmt.Println(delivery.Number)

		}

	}

	db.Close()
}
