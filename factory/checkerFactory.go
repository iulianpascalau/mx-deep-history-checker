package factory

import (
	"context"
	"fmt"
	"path/filepath"
	"sync"

	"github.com/iulianpascalau/mx-deep-history-checker/internal/checker"
	"github.com/iulianpascalau/mx-deep-history-checker/internal/config"
	"github.com/iulianpascalau/mx-deep-history-checker/internal/scanner"
)

const shardMeta = "Shard_metachain"

func DeepHistoryCheck(ctx context.Context, reporter Reporter, cfg *config.Config) error {
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, cfg.ParallelEpochs)

	epochScanner := scanner.NewTraverser()
	epochs, err := epochScanner.FindEpochs(cfg)
	if err != nil {
		return fmt.Errorf("failed to scan epochs: %w", err)
	}

	reporter.LogProgress(fmt.Sprintf("Found %d epochs matching criteria.", len(epochs)))
	dbChecker := checker.NewChecker(reporter)

	for _, epoch := range epochs {
		wg.Add(1)
		select {
		case semaphore <- struct{}{}:
		case <-ctx.Done():
			return fmt.Errorf("context done")
		}

		go func(epochPath string) {
			defer wg.Done()
			defer func() {
				<-semaphore
			}() // Release token

			path := filepath.Join(epochPath, cfg.Shard)
			dbChecker.CheckPath(ctx, path, cfg.MandatoryEpochDirs...)
		}(epoch)
	}

	wg.Wait()

	if cfg.CheckStatic {
		reporter.LogProgress("Checking Static directory databases...")
		staticPath := filepath.Join(cfg.NodeDir, "1", "Static", cfg.Shard)

		if cfg.Shard == shardMeta {
			dbChecker.CheckPath(ctx, staticPath, cfg.MandatoryStaticDirsForMeta...)
		} else {
			dbChecker.CheckPath(ctx, staticPath, cfg.MandatoryStaticDirsForShard...)
		}
	}

	return nil
}
