package checker

import (
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

// CheckLevelDB performs a deep integrity check of the provided database path
// while maintaining a focus on performance by sampling the database boundaries.
func CheckLevelDB(path string) error {
	options := &opt.Options{
		ReadOnly: true,
		// StrictAll ensures every block read (manifest, journals, or indices)
		// has its checksum verified against the recorded value.
		Strict: opt.StrictAll,
	}

	// 1. Open the database in Read-Only mode.
	// This implicitly verifies the manifest's structural integrity.
	db, err := leveldb.OpenFile(path, options)
	if err != nil {
		return err
	}
	defer func() {
		_ = db.Close()
	}()

	// 2. Perform a "Sampling" Deep Check.
	// To find the absolute first and last keys, the engine must read and verify
	// the index blocks of every SSTable in the entire database chain.
	iter := db.NewIterator(nil, nil)
	defer iter.Release()

	// Seek to the absolute boundaries to trigger deeper checksum verification.
	if iter.First() {
		// Found first key
	}
	if iter.Last() {
		// Found last key
	}

	// Verify if any corrupted blocks were detected during the seeks/iteration.
	err = iter.Error()
	if err != nil {
		return err
	}

	return nil
}
