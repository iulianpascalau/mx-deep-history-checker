package checker

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
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

func (c *structureChecker) CheckPath(ctx context.Context, path string, mandatoryDirectories ...string) {
	for _, dir := range mandatoryDirectories {
		select {
		case <-ctx.Done():
			return
		default:
		}

		dbBasePath := filepath.Join(path, dir)
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

			_, err = os.Stat(shardedDBPath)
			// Validate that the folder actually exists
			if os.IsNotExist(err) {
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

	err := CheckLevelDB(dbPath)
	if err != nil {
		c.rep.LogError(dbPath, fmt.Errorf("failed to open leveldb: %w", err))
		return
	}

	// Since NewDB succeeds, we consider this DB to be readable and primarily uncorrupted.
	c.rep.LogSuccess(dbPath)
}
