package handler

import (
	"os/user"

	"github.com/OK0X/ethereum-chaindata-flashreader/src/ledger"
)

var (
	DB        *ledger.DB
	FlashRead *ledger.FlashRead
	Serv      *Server
)

func Run() {
	usr, _ := user.Current()
	dataDir := usr.HomeDir + "/geth/data/geth/chaindata"
	DB = &ledger.DB{}
	DB.Initialize(dataDir)

	FlashRead = &ledger.FlashRead{}
	FlashRead.Initialize(DB.Database)

	Serv = &Server{}
	Serv.Initialize()
	Serv.Start()

}

func Stop() {
	DB.Close()
	Serv.Stop()
}
