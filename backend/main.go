package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/launchpals/open-now/backend/env"
	"github.com/launchpals/open-now/backend/maps"
	"github.com/launchpals/open-now/backend/service"
	"go.uber.org/zap"
)

func main() {
	// load env vars
	godotenv.Load()
	var vals = env.Load()

	// init logger
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

	// connect to maps API
	m, err := maps.NewClient(l.Named("maps"), vals.GCPKey)
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}

	// spin up service
	s, err := service.New(l.Named("service"), m)
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
	if err := s.Run(ctx, vals.Host, "8081"); err != nil {
		l.Fatalw("server stopped",
			"error", err)
	}
}
