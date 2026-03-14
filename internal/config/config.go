package config

// Config holds the application configuration parsed from command line flags.
type Config struct {
	NodeDir                     string
	StartEpoch                  uint64
	EndEpoch                    uint64
	CheckStatic                 bool
	ParallelEpochs              uint
	Shard                       string
	MandatoryEpochDirs          []string
	MandatoryStaticDirsForShard []string
	MandatoryStaticDirsForMeta  []string
}
