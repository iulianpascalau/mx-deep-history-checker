package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/iulianpascalau/mx-deep-history-checker/factory"
	"github.com/iulianpascalau/mx-deep-history-checker/internal/config"
	"github.com/iulianpascalau/mx-deep-history-checker/internal/reporter"
	logger "github.com/multiversx/mx-chain-logger-go"
)

var log = logger.GetOrCreate("checker")

func main() {
	err := run()
	if err != nil {
		log.LogIfError(err)
		os.Exit(1)
	}
}

func run() error {
	cfgHandler := config.NewConfigHandler()
	cfg, err := cfgHandler.Parse()
	if err != nil {
		return fmt.Errorf("failed to parse configuration: %w", err)
	}

	rep := reporter.NewReporter()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		sigs := make(chan os.Signal, 1)
		signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
		<-sigs
		cancel()
	}()

	err = factory.DeepHistoryCheck(ctx, rep, cfg)
	if err != nil {
		return err
	}

	rep.PrintSummary()

	return nil
}
