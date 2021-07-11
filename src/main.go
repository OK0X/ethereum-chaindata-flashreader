package main

import (
	"encoding/json"
	"fmt"
	"os/user"

	"github.com/ethereum/go-ethereum/core/rawdb"

	"github.com/OK0X/ethereum-chaindata-flashreader/src/ledger"
)

func main() {

	usr, _ := user.Current()
	dataDir := usr.HomeDir + "/geth/data/geth/chaindata"
	db, err := rawdb.NewLevelDBDatabase(dataDir, 128, 1024, "", false)

	if err != nil {
		fmt.Errorf("error opening database at %v: %v", dataDir, err)
	}

	fr := &ledger.FlashRead{}
	fr.Initialize(db)
	blocks, _ := fr.ReadTransactions(0, 12, false, nil)
	txStr, _ := json.Marshal(blocks)
	fmt.Println(string(txStr))
	db.Close()
}
