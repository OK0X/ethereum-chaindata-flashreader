package ledger

import (
	"fmt"
	"runtime"
	"sync/atomic"

	"github.com/OK0X/ethereum-chaindata-flashreader/src/model"
	"github.com/ethereum/go-ethereum/common/prque"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/rlp"
)

type FlashRead struct {
	DB ethdb.Database
}

type BlockTxs struct {
	Number uint64
	Txs    types.Transactions
}

func (fr *FlashRead) Initialize(db ethdb.Database) {
	fr.DB = db
}

func (fr *FlashRead) ReadTransactions(from uint64, to uint64, reverse bool, interrupt chan struct{}) (model.Blocks, error) {
	if from < 0 {
		return nil, fmt.Errorf("block number from should bigger or equal 0")
	}

	// ToDo
	// //headBlockHash := rawdb.ReadHeadBlockHash(db)
	// //headBlockNumber := rawdb.ReadHeaderNumber(db, headBlockHash)
	// if to > lastBlockNumber {
	// 	return nil, fmt.Errorf("block number to should less or equal lastest")
	// }

	if from >= to {
		return nil, fmt.Errorf("block number from should less than to")
	}

	blocks := make(model.Blocks, 0, to-from)

	lastNum := to
	txsCh := fr.readTransactions(from, to, false, nil)

	queue := prque.New(nil)

	for chanDelivery := range txsCh {
		queue.Push(chanDelivery, int64(chanDelivery.Number))
		for !queue.Empty() {
			// If the next available item is gapped, return
			if _, priority := queue.Peek(); priority != int64(lastNum-1) {
				break
			}

			// Next block available, pop it off and index it
			delivery := queue.PopItem().(*BlockTxs)
			lastNum = delivery.Number
			mtxs := make([]*model.Transaction, 0, delivery.Txs.Len())
			for i := 0; i < delivery.Txs.Len(); i++ {
				tx := delivery.Txs[i]
				v, r, s := tx.RawSignatureValues()
				mtx := &model.Transaction{
					Type:       tx.Type(),
					ChainId:    tx.ChainId(),
					Data:       tx.Data(),
					AccessList: tx.AccessList(),
					Gas:        tx.Gas(),
					GasPrice:   tx.GasPrice(),
					GasTipCap:  tx.GasTipCap(),
					GasFeeCap:  tx.GasFeeCap(),
					Value:      tx.Value(),
					Nonce:      tx.Nonce(),
					To:         tx.To(),
					Cost:       tx.Cost(),
					V:          v,
					R:          r,
					S:          s,
					Hash:       tx.Hash(),
					Size:       tx.Size(),
				}

				mtxs = append(mtxs, mtx)

			}

			block := &model.Block{
				Number:       lastNum,
				Transactions: mtxs,
			}

			blocks = append(blocks, block)

		}

	}

	return blocks, nil
}

func (fr *FlashRead) readTransactions(from uint64, to uint64, reverse bool, interrupt chan struct{}) chan *BlockTxs {
	// One thread sequentially reads data from db
	type numberRlp struct {
		number uint64
		rlp    rlp.RawValue
	}
	if to == from {
		return nil
	}
	threads := to - from
	if cpus := runtime.NumCPU(); threads > uint64(cpus) {
		threads = uint64(cpus)
	}
	var (
		rlpCh = make(chan *numberRlp, threads*2) // we send raw rlp over this channel
		txsCh = make(chan *BlockTxs, threads*2)  // send hashes over hashesCh
	)
	// lookup runs in one instance
	lookup := func() {
		n, end := from, to
		if reverse {
			n, end = to-1, from-1
		}
		defer close(rlpCh)
		for n != end {
			data := rawdb.ReadCanonicalBodyRLP(fr.DB, n)
			// Feed the block to the aggregator, or abort on interrupt
			select {
			case rlpCh <- &numberRlp{n, data}:
			case <-interrupt:
				return
			}
			if reverse {
				n--
			} else {
				n++
			}
		}
	}
	// process runs in parallel
	nThreadsAlive := int32(threads)
	process := func() {
		defer func() {
			// Last processor closes the result channel
			if atomic.AddInt32(&nThreadsAlive, -1) == 0 {
				close(txsCh)
			}
		}()
		for data := range rlpCh {
			var body types.Body
			if err := rlp.DecodeBytes(data.rlp, &body); err != nil {
				log.Warn("Failed to decode block body", "block", data.number, "error", err)
				return
			}
			result := &BlockTxs{
				Number: data.number,
				Txs:    body.Transactions,
			}
			// Feed the block to the aggregator, or abort on interrupt
			select {
			case txsCh <- result:
			case <-interrupt:
				return
			}
		}
	}
	go lookup() // start the sequential db accessor
	for i := 0; i < int(threads); i++ {
		go process()
	}
	return txsCh
}
