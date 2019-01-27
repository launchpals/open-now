package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/lunchpals/open-now/backend/service"
	"go.uber.org/zap"
)

func main() {
	bareLogger, err := zap.NewDevelopment()
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
	var l = bareLogger.Sugar()

	// catch interrupts
	ctx, cancel := context.WithCancel(context.Background())
	var signals = make(chan os.Signal)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	go func() {
		<-signals
		cancel()
	}()

	// spin up service
	s, err := service.New(l)
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
	if err := s.Run(ctx, "127.0.0.1", "8081"); err != nil {
		l.Fatalw("server stopped",
			"error", err)
	}
}
