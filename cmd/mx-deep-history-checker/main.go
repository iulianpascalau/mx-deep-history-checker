package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"

	"github.com/iulianpascalau/mx-deep-history-checker/internal/checker"
	"github.com/iulianpascalau/mx-deep-history-checker/internal/config"
	"github.com/iulianpascalau/mx-deep-history-checker/internal/reporter"
	"github.com/iulianpascalau/mx-deep-history-checker/internal/scanner"
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
	epochScanner := scanner.NewTraverser()

	epochs, err := epochScanner.FindEpochs(cfg)
	if err != nil {
		return fmt.Errorf("failed to scan epochs: %w", err)
	}

	rep.LogProgress(fmt.Sprintf("Found %d epochs matching criteria.", len(epochs)))

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, cfg.ParallelEpochs)
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for _, epoch := range epochs {
		wg.Add(1)
		select {
		case semaphore <- struct{}{}:
		case <-sigs:
			return fmt.Errorf("interrupted by user")
		}

		go func(epochPath string) {
			defer wg.Done()
			defer func() {
				<-semaphore
			}() // Release token

			dbChecker := checker.NewChecker(rep)
			err = dbChecker.CheckEpoch(epochPath, ctx)
			if err != nil {
				rep.LogError(epochPath, fmt.Errorf("epoch check failed: %w", err))
			}
		}(epoch)
	}

	wg.Wait()

	if cfg.CheckStatic {
		rep.LogProgress("Checking Static directory databases...")
		staticPath := filepath.Join(cfg.NodeDir, "1", "Static")
		dbChecker := checker.NewChecker(rep)
		err = dbChecker.CheckStatic(staticPath, ctx)
		if err != nil {
			rep.LogError(staticPath, fmt.Errorf("static check failed: %w", err))
		}
	}

	rep.PrintSummary()

	return nil
}
