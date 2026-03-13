package config

import (
	"flag"
	"fmt"
)

// Config holds the application configuration parsed from command line flags.
type Config struct {
	NodeDir                     string
	StartEpoch                  uint32
	EndEpoch                    *uint32
	CheckStatic                 bool
	ParallelEpochs              uint
	Shard                       string
	MandatoryEpochDirs          []string
	MandatoryStaticDirsForShard []string
	MandatoryStaticDirsForMeta  []string
}

type cliConfigHandler struct{}

// NewConfigHandler creates a new instance of ConfigHandler.
func NewConfigHandler() *cliConfigHandler {
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
	shard := flag.String("shard", "Shard_0", "The shard to be checked. Example: Shard_0, Shard_1, Shard_2, Shard_metachain")

	flag.Parse()

	if *nodeDir == "" {
		return nil, fmt.Errorf("--node-dir is required")
	}

	cfg := &Config{
		NodeDir:                     *nodeDir,
		StartEpoch:                  uint32(*startEpoch),
		CheckStatic:                 *checkStatic,
		ParallelEpochs:              *parallelEpochs,
		Shard:                       *shard,
		MandatoryEpochDirs:          h.getMandatoryEpochDirs(),
		MandatoryStaticDirsForShard: h.getMandatoryStaticDirsForShard(),
		MandatoryStaticDirsForMeta:  h.getMandatoryStaticDirsForMeta(),
	}

	if *endEpochRaw != -1 {
		val := uint32(*endEpochRaw)
		cfg.EndEpoch = &val
	}

	return cfg, nil
}

func (h *cliConfigHandler) getMandatoryEpochDirs() []string {
	return []string{
		"AccountsTrie",
		"AccountsTrieCheckpoints",
		"BlockHeaders",
		"BootstrapData",
		"DbLookupExtensions",
		"DbLookupExtensions_ResultsHashesByTx",
		"Logs",
		"MetaBlock",
		"MiniBlocks",
		"PeerAccountsTrie",
		"PeerAccountsTrieCheckpoints",
		"Receipts",
		"RewardTransactions",
		"ScheduledSCRs",
		"Transactions",
		"UnsignedTransactions",
	}
}

func (h *cliConfigHandler) getMandatoryStaticDirsForShard() []string {
	return []string{
		"DbLookupExtensions_EpochByHash",
		"DbLookupExtensions_ESDTSupplies",
		"DbLookupExtensions_MiniblockHashByTxHash",
		"DbLookupExtensions_RoundHash",
		"MetaHdrHashNonce",
		"ShardHdrHashNonce0",
		"StatusMetricsStorageDB",
	}
}

func (h *cliConfigHandler) getMandatoryStaticDirsForMeta() []string {
	return []string{
		"DbLookupExtensions_EpochByHash",
		"DbLookupExtensions_ESDTSupplies",
		"DbLookupExtensions_MiniblockHashByTxHash",
		"DbLookupExtensions_RoundHash",
		"MetaHdrHashNonce",
		"ShardHdrHashNonce0",
		"ShardHdrHashNonce1",
		"ShardHdrHashNonce2",
		"StatusMetricsStorageDB",
	}
}
