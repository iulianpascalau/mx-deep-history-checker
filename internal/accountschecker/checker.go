package accountschecker

import (
	"context"
	"fmt"
	"log"
	"time"
)

type checkerImpl struct {
	client       APIClient
	addresses    []string
	startNonce   uint64
	increment    uint64
	reportWindow uint64
}

// NewChecker creates a new Checker instance that loops through addresses and nonces.
func NewChecker(client APIClient, addresses []string, startNonce, increment, reportWindow uint64) Checker {
	return &checkerImpl{
		client:       client,
		addresses:    addresses,
		startNonce:   startNonce,
		increment:    increment,
		reportWindow: reportWindow,
	}
}

func (c *checkerImpl) Run(ctx context.Context) error {
	log.Println("Fetching highest nonce...")
	highestNonce, err := c.client.GetHighestNonce(ctx)
	if err != nil {
		return fmt.Errorf("failed to get highest nonce: %w", err)
	}
	log.Printf("Highest nonce is: %d\n", highestNonce)

	var reqCount uint64
	startTime := time.Now()
	windowStart := time.Now()
	var windowReqCount uint64

	for nonce := c.startNonce; nonce <= highestNonce; nonce += c.increment {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		for _, addr := range c.addresses {
			acc, errGet := c.client.GetAccountData(ctx, addr, nonce)
			if errGet != nil {
				return fmt.Errorf("failed fetching account %s at nonce %d: %w", addr, nonce, errGet)
			}

			if !acc.IsValid() {
				// Specifications state:
				// "If the response is not valid, the program will halt, displaying the address, nonce and the returned fields."
				log.Fatalf("HALT! Invalid response for address %s at nonce %d.\nReturned fields:\nNonce: %d\nBalance: %s\nRootHash: %v\n",
					addr, nonce, acc.Nonce, acc.Balance, formatRootHash(acc.RootHash))
			}

			reqCount++
			windowReqCount++
		}

		// Calculate when to report progress, e.g. every 100,000 nonces range.
		if (nonce-c.startNonce) > 0 && (nonce-c.startNonce)%c.reportWindow == 0 {
			elapsed := time.Since(windowStart)
			if elapsed.Seconds() > 0 {
				reqPerSec := float64(windowReqCount) / elapsed.Seconds()
				log.Printf("Progress: nonce %d / %d. Checked %d nonces range. Speed: %.2f req/s (average over last window)",
					nonce, highestNonce, c.reportWindow, reqPerSec)
			}
			// Reset window for next batch
			windowStart = time.Now()
			windowReqCount = 0
		}
	}

	totalElapsed := time.Since(startTime)
	if totalElapsed.Seconds() > 0 {
		overallSpeed := float64(reqCount) / totalElapsed.Seconds()
		log.Printf("Finished. Total requests: %d. Overall Speed: %.2f req/s", reqCount, overallSpeed)
	} else {
		log.Printf("Finished. Total requests: %d. Time taken was very short.", reqCount)
	}

	return nil
}

func formatRootHash(r *string) string {
	if r == nil {
		return "null"
	}
	return *r
}
