package main

import (
	"context"
	"os"

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

	s, err := service.New(l)
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
	s.Run(context.Background(), "127.0.0.1", "8081")
}
