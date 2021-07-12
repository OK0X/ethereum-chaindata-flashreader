package handler

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/OK0X/ethereum-chaindata-flashreader/src/config"
)

type Server struct {
	*http.Server
}

func (s *Server) Initialize() {

	if !config.Debug {
		gin.SetMode(gin.ReleaseMode)
	}
	e := gin.New()
	AddRouter(e)
	s.Server = &http.Server{
		Addr:    config.Addr,
		Handler: e,
	}
}

func (s *Server) Start() {
	go func() {
		if err := s.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}
	}()
}

// Stop stops service
func (s *Server) Stop() {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := s.Server.Shutdown(ctx); err != nil {
		fmt.Printf("Error: Server forced to shutdown: %s", err.Error())
	} else {
		fmt.Println("server stoped")
	}
}
