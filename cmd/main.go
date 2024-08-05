// Package main is the etry point
package main

import (
	"fmt"

	"github.com/vtpl1/logutilgo/pkg/logger"
)

func main() {
	log, err := logger.New(
		logger.LogConfig{
			AppName:       "",
			SessionFolder: "",
			ConsoleLog:    true,
		},
	)
	if err != nil {
		fmt.Println("Unexpected logger error")
		return
	}
	log.Info().Msg("Hello")
}
