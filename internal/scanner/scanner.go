package scanner

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/iulianpascalau/mx-deep-history-checker/internal/config"
)

type fileSystemTraverser struct{}

// NewTraverser creates a new DirectoryTraverser.
func NewTraverser() *fileSystemTraverser {
	return &fileSystemTraverser{}
}

func (t *fileSystemTraverser) FindEpochs(cfg *config.Config) ([]string, error) {
	basePath := filepath.Join(cfg.NodeDir, "1")

	entries, err := os.ReadDir(basePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read node directory %s: %w", basePath, err)
	}

	var validEpochs []string
	epochsNum := make([]uint64, 0, len(entries))

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		name := entry.Name()
		if !strings.HasPrefix(name, "Epoch_") {
			continue
		}

		epochStr := strings.TrimPrefix(name, "Epoch_")
		epochNum := uint64(0)
		epochNum, err = strconv.ParseUint(epochStr, 10, 32)
		if err != nil {
			// Skip directories that don't match the exact pattern
			continue
		}

		epochsNum = append(epochsNum, epochNum)

		if epochNum >= cfg.StartEpoch && epochNum <= cfg.EndEpoch {
			validEpochs = append(validEpochs, filepath.Join(basePath, name))
		}
	}

	sort.SliceStable(epochsNum, func(i, j int) bool { return epochsNum[i] < epochsNum[j] })
	for i := 1; i < len(epochsNum); i++ {
		if epochsNum[i] != epochsNum[i-1]+1 {
			return nil, fmt.Errorf("epochs are not consecutive: %d, %d", epochsNum[i-1], epochsNum[i])
		}
	}

	return validEpochs, nil
}
