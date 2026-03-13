package config

import (
	"flag"
	"fmt"
)

// Config holds the application configuration parsed from command line flags.
type Config struct {
	NodeDir        string
	StartEpoch     uint32
	EndEpoch       *uint32
	CheckStatic    bool
	ParallelEpochs uint
}

// ConfigHandler defines the interface for managing application configuration.
type ConfigHandler interface {
	// Parse loads the configuration from command line arguments.
	Parse() (*Config, error)
}

type cliConfigHandler struct{}

// NewConfigHandler creates a new instance of ConfigHandler.
func NewConfigHandler() ConfigHandler {
	return &cliConfigHandler{}
}

// Parse implements ConfigHandler.Parse by reading the CLI flags.
func (h *cliConfigHandler) Parse() (*Config, error) {
	nodeDir := flag.String("node-dir", "", "The root path to the node data directory")
	startEpoch := flag.Uint("start-epoch", 0, "The starting epoch number to check (inclusive)")

	// Default to nil, but we can't easily do a pointer via standard flag without a custom Var.
	// Alternative: int64 initialized to -1.
	endEpochRaw := flag.Int64("end-epoch", -1, "The ending epoch number to check (inclusive). If omitted, goes to the highest one.")

	checkStatic := flag.Bool("check-static", true, "Check the Static directory databases")
	parallelEpochs := flag.Uint("parallel-epochs", 4, "The number of epochs to process in parallel")

	flag.Parse()

	if *nodeDir == "" {
		return nil, fmt.Errorf("--node-dir is required")
	}

	cfg := &Config{
		NodeDir:        *nodeDir,
		StartEpoch:     uint32(*startEpoch),
		CheckStatic:    *checkStatic,
		ParallelEpochs: uint(*parallelEpochs),
	}

	if *endEpochRaw != -1 {
		val := uint32(*endEpochRaw)
		cfg.EndEpoch = &val
	}

	return cfg, nil
}
