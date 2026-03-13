package checker

import (
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/opt"
)

func CheckLevelDB(path string) error {
	options := &opt.Options{
		// disable internal cache
		BlockCacheCapacity:     -1,
		OpenFilesCacheCapacity: 10,
		ReadOnly:               true,
	}

	db, errOpen := leveldb.OpenFile(path, options)
	if errOpen != nil {
		return errOpen
	}

	_ = db.Close()
	return nil
}
