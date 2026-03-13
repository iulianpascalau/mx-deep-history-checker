package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/iulianpascalau/mx-deep-history-checker/factory"
	"github.com/iulianpascalau/mx-deep-history-checker/internal/config"
	"github.com/iulianpascalau/mx-deep-history-checker/internal/reporter"
	logger "github.com/multiversx/mx-chain-logger-go"
)

// appVersion should be populated at build time using ldflags
// Usage examples:
// Linux/macOS:
//
//	go build -v -ldflags="-X main.appVersion=$(git describe --all | cut -c7-32)
var appVersion = "undefined"
var log = logger.GetOrCreate("checker")

func main() {
	appVersion = fmt.Sprintf("%s/%s/%s-%s", appVersion, runtime.Version(), runtime.GOOS, runtime.GOARCH)
	log.Info("Starting deep history checker", "version", appVersion, "pid", os.Getpid())

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
