package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/OK0X/ethereum-chaindata-flashreader/src/config"
	"github.com/OK0X/ethereum-chaindata-flashreader/src/handler"
)

func main() {

	handler.Run()
	fmt.Printf("server start success, listening at: %v \n", config.Addr)
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)
	<-exit
	handler.Stop()

}
