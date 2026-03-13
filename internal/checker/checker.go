package checker

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/multiversx/mx-chain-storage-go/leveldb"
)

type structureChecker struct {
	rep Reporter
}

// NewChecker creates a new DatabaseChecker instance.
func NewChecker(rep Reporter) *structureChecker {
	return &structureChecker{
		rep: rep,
	}
}

func (c *structureChecker) CheckEpoch(epochPath string, ctx context.Context) error {
	entries, err := os.ReadDir(epochPath)
	if err != nil {
		return fmt.Errorf("failed to read epoch directory %s: %w", epochPath, err)
	}

	for _, entry := range entries {
		select {
		case <-ctx.Done():
			return nil
		default:
		}

		if !entry.IsDir() || !strings.HasPrefix(entry.Name(), "Shard_") {
			continue
		}
		shardPath := filepath.Join(epochPath, entry.Name())
		c.checkShard(shardPath, ctx)
	}

	return nil
}

func (c *structureChecker) CheckStatic(staticPath string, ctx context.Context) error {
	entries, err := os.ReadDir(staticPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // Static folder might not exist, which is fine
		}
		return fmt.Errorf("failed to read static directory %s: %w", staticPath, err)
	}

	for _, entry := range entries {
		if !entry.IsDir() || !strings.HasPrefix(entry.Name(), "Shard_") {
			continue
		}
		shardPath := filepath.Join(staticPath, entry.Name())
		c.checkShard(shardPath, ctx)
	}

	return nil
}

func (c *structureChecker) checkShard(shardPath string, ctx context.Context) {
	entries, err := os.ReadDir(shardPath)
	if err != nil {
		c.rep.LogError(shardPath, fmt.Errorf("failed to read shard dir: %w", err))
		return
	}

	for _, entry := range entries {
		select {
		case <-ctx.Done():
			return
		default:
		}

		if !entry.IsDir() {
			continue
		}

		dbBasePath := filepath.Join(shardPath, entry.Name())
		c.processDatabaseDirectory(dbBasePath)
	}
}

func (c *structureChecker) processDatabaseDirectory(dbBasePath string) {
	configPath := filepath.Join(dbBasePath, "config.toml")

	cfg, err := readDBConfig(configPath)
	if err != nil {
		c.rep.LogError(dbBasePath, fmt.Errorf("config error: %w", err))
		return
	}

	if cfg != nil && cfg.Type == "LvlDBSerial" && cfg.NumShards > 0 {
		// Sharded DB handling
		for i := uint(0); i < cfg.NumShards; i++ {
			shardedDBPath := filepath.Join(dbBasePath, fmt.Sprintf("%d", i))

			// Validate that the folder actually exists
			if _, err := os.Stat(shardedDBPath); os.IsNotExist(err) {
				c.rep.LogError(shardedDBPath, fmt.Errorf("expected shard %d directory missing", i))
				continue
			}

			c.verifyLevelDB(shardedDBPath)
		}
	} else {
		// Standard DB handling
		c.verifyLevelDB(dbBasePath)
	}
}

func (c *structureChecker) verifyLevelDB(dbPath string) {
	c.rep.LogProgress(fmt.Sprintf("Checking %s...", dbPath))

	// The mx-chain-storage-go library provides ways to interact with level DB.
	// Signature of leveldb.NewDB is (filePath string, batchDelaySeconds int, maxBatchSize int, maxOpenFiles int)
	db, err := leveldb.NewDB(dbPath, 1, 100, 10)
	if err != nil {
		c.rep.LogError(dbPath, fmt.Errorf("failed to open leveldb: %w", err))
		return
	}

	// Successfully opened, so we immediately close it to free locks.
	defer func() {
		_ = db.Close()
	}()

	// Since NewDB succeeds, we consider this DB to be readable and primarily uncorrupted.
	c.rep.LogSuccess(dbPath)
}
