package ledger

import (
	"fmt"

	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/ethdb"
)

type DB struct {
	ethdb.Database
}

func (db *DB) Initialize(dataDir string) {
	ethdb, err := rawdb.NewLevelDBDatabase(dataDir, 128, 1024, "", false)
	db.Database = ethdb
	if err != nil {
		fmt.Printf("error opening database at %v: %v", dataDir, err)
	}
}

func (db *DB) Close() {
	err := db.Database.Close()
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("db closed")
	}
}
