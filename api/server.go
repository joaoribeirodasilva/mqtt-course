package main

import (
	"context"
	"fmt"
	"net/http"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joaoribeirodasilva/mqtt-course/api/configuration"
	"github.com/joaoribeirodasilva/wait_signals"
)

type Server struct {
	conf   *configuration.Configuration
	Router *gin.Engine
}

func NewServer(conf *configuration.Configuration) *Server {

	s := &Server{}

	s.conf = conf
	s.Router = gin.Default()

	return s
}

func (s *Server) Listen() error {

	addr := fmt.Sprintf("%s:%d", s.conf.Server.Address, s.conf.Server.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: s.Router,
	}

	go func() {
		fmt.Printf("INFO: [HTTP SERVER] listening at %s ", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			panic(err)
		}

	}()

	wait_signals.Wait(syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		return fmt.Errorf("ERROR: [HTTP SERVER] failed to stop server REASON: %v", err)
	}

	fmt.Printf("INFO: [HTTP SERVER] stopped")

	return nil
}
