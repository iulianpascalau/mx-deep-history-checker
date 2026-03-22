package integrationTests

import (
	"context"
	"math"
	"testing"

	"github.com/iulianpascalau/mx-deep-history-checker/factory"
	"github.com/iulianpascalau/mx-deep-history-checker/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func applyMandatoryDirs(cfg *config.Config) {
	cfg.MandatoryEpochDirs = []string{"AccountsTrie", "DbLookupExtensions/MiniblocksMetadata", "MiniBlocks"}
	cfg.MandatoryStaticDirsForShard = []string{"ShardHdrHashNonce0"}
	cfg.MandatoryStaticDirsForMeta = []string{"ShardHdrHashNonce0", "ShardHdrHashNonce1", "ShardHdrHashNonce2"}
}

func TestCheckerOk(t *testing.T) {
	t.Run("shard 0", func(t *testing.T) {
		expectedSuccessLogs := []string{
			"data/ok/1/Epoch_0/Shard_0/AccountsTrie",
			"data/ok/1/Epoch_0/Shard_0/DbLookupExtensions/MiniblocksMetadata",
			"data/ok/1/Epoch_0/Shard_0/MiniBlocks",
			"data/ok/1/Epoch_1/Shard_0/AccountsTrie/0",
			"data/ok/1/Epoch_1/Shard_0/AccountsTrie/1",
			"data/ok/1/Epoch_1/Shard_0/DbLookupExtensions/MiniblocksMetadata",
			"data/ok/1/Epoch_1/Shard_0/MiniBlocks",
			"data/ok/1/Static/Shard_0/ShardHdrHashNonce0",
		}

		testTestCheckerOk(t, "Shard_0", expectedSuccessLogs)
	})
	t.Run("shard metachain", func(t *testing.T) {
		expectedSuccessLogs := []string{
			"data/ok/1/Epoch_0/Shard_metachain/AccountsTrie",
			"data/ok/1/Epoch_0/Shard_metachain/DbLookupExtensions/MiniblocksMetadata",
			"data/ok/1/Epoch_0/Shard_metachain/MiniBlocks",
			"data/ok/1/Epoch_1/Shard_metachain/AccountsTrie/0",
			"data/ok/1/Epoch_1/Shard_metachain/AccountsTrie/1",
			"data/ok/1/Epoch_1/Shard_metachain/DbLookupExtensions/MiniblocksMetadata",
			"data/ok/1/Epoch_1/Shard_metachain/MiniBlocks",
			"data/ok/1/Static/Shard_metachain/ShardHdrHashNonce0",
			"data/ok/1/Static/Shard_metachain/ShardHdrHashNonce1",
			"data/ok/1/Static/Shard_metachain/ShardHdrHashNonce2",
		}

		testTestCheckerOk(t, "Shard_metachain", expectedSuccessLogs)
	})
}

func TestCheckerMissingEpoch(t *testing.T) {
	cfg := &config.Config{
		NodeDir:        "./data/missing-epoch",
		StartEpoch:     0,
		EndEpoch:       math.MaxUint64,
		CheckStatic:    true,
		ParallelEpochs: 1,
		Shard:          "Shard_0",
	}
	applyMandatoryDirs(cfg)

	rep := NewTestReporter()

	err := factory.DeepHistoryCheck(context.Background(), rep, cfg)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "epochs are not consecutive: 0, 2")
}

func TestCheckerMissingMandatoryMiniblocksDir(t *testing.T) {
	cfg := &config.Config{
		NodeDir:        "./data/missing-miniblocks-dir",
		StartEpoch:     0,
		EndEpoch:       math.MaxUint64,
		CheckStatic:    true,
		ParallelEpochs: 1,
		Shard:          "Shard_0",
	}
	applyMandatoryDirs(cfg)

	rep := NewTestReporter()

	err := factory.DeepHistoryCheck(context.Background(), rep, cfg)
	require.Nil(t, err)

	expectedMissingDirs := []string{
		"data/missing-miniblocks-dir/1/Epoch_1/Shard_0/MiniBlocks",
	}
	require.Equal(t, expectedMissingDirs, rep.GetErrorLogs())
}

func TestCheckerMissingMandatoryAccountsDBSubDir(t *testing.T) {
	cfg := &config.Config{
		NodeDir:        "./data/missing-accounts-subdir",
		StartEpoch:     0,
		EndEpoch:       math.MaxUint64,
		CheckStatic:    true,
		ParallelEpochs: 1,
		Shard:          "Shard_0",
	}
	applyMandatoryDirs(cfg)

	rep := NewTestReporter()

	err := factory.DeepHistoryCheck(context.Background(), rep, cfg)
	require.Nil(t, err)

	expectedMissingDirs := []string{
		"data/missing-accounts-subdir/1/Epoch_1/Shard_0/AccountsTrie/1",
	}
	require.Equal(t, expectedMissingDirs, rep.GetErrorLogs())
}

func TestCheckerMissingMandatoryStaticShardDir(t *testing.T) {
	cfg := &config.Config{
		NodeDir:        "./data/missing-static-shard-dir",
		StartEpoch:     0,
		EndEpoch:       math.MaxUint64,
		CheckStatic:    true,
		ParallelEpochs: 1,
		Shard:          "Shard_0",
	}
	applyMandatoryDirs(cfg)

	rep := NewTestReporter()

	err := factory.DeepHistoryCheck(context.Background(), rep, cfg)
	require.Nil(t, err)

	expectedMissingDirs := []string{
		"data/missing-static-shard-dir/1/Static/Shard_0/ShardHdrHashNonce0",
	}
	require.Equal(t, expectedMissingDirs, rep.GetErrorLogs())
}

func TestCheckerMissingMandatoryStaticMetaDir(t *testing.T) {
	cfg := &config.Config{
		NodeDir:        "./data/missing-static-meta-dir",
		StartEpoch:     0,
		EndEpoch:       math.MaxUint64,
		CheckStatic:    true,
		ParallelEpochs: 1,
		Shard:          "Shard_metachain",
	}
	applyMandatoryDirs(cfg)

	rep := NewTestReporter()

	err := factory.DeepHistoryCheck(context.Background(), rep, cfg)
	require.Nil(t, err)

	expectedMissingDirs := []string{
		"data/missing-static-meta-dir/1/Static/Shard_metachain/ShardHdrHashNonce1",
	}
	require.Equal(t, expectedMissingDirs, rep.GetErrorLogs())
}

func TestCheckerCorruptedMandatoryAccountsDir(t *testing.T) {
	cfg := &config.Config{
		NodeDir:        "./data/corrupted-accounts-dir",
		StartEpoch:     0,
		EndEpoch:       math.MaxUint64,
		CheckStatic:    true,
		ParallelEpochs: 1,
		Shard:          "Shard_0",
	}
	applyMandatoryDirs(cfg)

	rep := NewTestReporter()

	err := factory.DeepHistoryCheck(context.Background(), rep, cfg)
	require.Nil(t, err)

	// Expect an error log indicating the DB is corrupted
	errorLogs := rep.GetErrorLogs()
	require.Len(t, errorLogs, 1, "Expected exactly one error log")
	assert.Contains(t, errorLogs[0], "data/corrupted-accounts-dir/1/Epoch_0/Shard_0/AccountsTrie")
	assert.Contains(t, errorLogs[0], "corrupted") // Or whatever the missing file error looks like
}

func TestCheckerCorruptedMandatoryAccountsDir2(t *testing.T) {
	cfg := &config.Config{
		NodeDir:        "./data/corrupted-accounts-dir-2",
		StartEpoch:     0,
		EndEpoch:       math.MaxUint64,
		CheckStatic:    true,
		ParallelEpochs: 1,
		Shard:          "Shard_0",
	}
	applyMandatoryDirs(cfg)

	rep := NewTestReporter()

	err := factory.DeepHistoryCheck(context.Background(), rep, cfg)
	require.Nil(t, err)

	// Expect an error log indicating the DB is corrupted
	errorLogs := rep.GetErrorLogs()
	require.Len(t, errorLogs, 1, "Expected exactly one error log")
	assert.Contains(t, errorLogs[0], "data/corrupted-accounts-dir-2/1/Epoch_0/Shard_0/DbLookupExtensions/MiniblocksMetadata")
	assert.Contains(t, errorLogs[0], "corrupted") // Or whatever the missing file error looks like
}

func testTestCheckerOk(tb testing.TB, shard string, expectedSuccessLogs []string) {
	cfg := &config.Config{
		NodeDir:        "./data/ok",
		StartEpoch:     0,
		EndEpoch:       math.MaxUint64,
		CheckStatic:    true,
		ParallelEpochs: 1,
		Shard:          shard,
	}
	applyMandatoryDirs(cfg)

	rep := NewTestReporter()

	err := factory.DeepHistoryCheck(context.Background(), rep, cfg)
	require.NoError(tb, err)

	require.Empty(tb, rep.GetErrorLogs())
	require.Equal(tb, expectedSuccessLogs, rep.GetSuccessLogs())
}
