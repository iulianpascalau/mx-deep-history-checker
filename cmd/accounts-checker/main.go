package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/iulianpascalau/mx-deep-history-checker/internal/accountschecker"
)

const (
	baseURL      = "https://mvx-deep-history.jls-software.net/"
	token        = "<token>"
	increment    = 5000
	startNonce   = 1
	reportWindow = 50000 // Report after processing a block of a size of reportWindow
)

var addresses = []string{
	"erd103r4tfg6x00jtcyzvara4nwjegrs4mvzmtfvxen3c3688z44d7yqmfs0gj",
	"erd105vcnjmzaaw2pd6awpfj6amwcnkm67wljucka7vefnwsk4e62cys4pydst",
	"erd100z5n2u5gre3fqzdhnu3twhjjt4m0ym8v8s5t2hg803ujrfhn7dqjs9fnf",
	"erd1qqqqqqqqqqqqqqqpqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqllls0lczs7",
}

func main() {
	client := accountschecker.NewAPIClient(baseURL, token)
	checker := accountschecker.NewChecker(client, addresses, startNonce, increment, reportWindow)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		log.Println("Received termination signal, shutting down...")
		cancel()
	}()

	log.Println("Starting accounts checker...")
	err := checker.Run(ctx)
	if err != nil {
		log.Fatalf("Checker exited with error: %v", err)
	}
	log.Println("Done.")
}
