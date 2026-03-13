package reporter

import (
	"fmt"
	"os"
	"sync/atomic"

	logger "github.com/multiversx/mx-chain-logger-go"
)

var log = logger.GetOrCreate("reporter")

type cliReporter struct {
	successCount uint32
	errorCount   uint32
	totalCount   uint32
	errorPaths   []string
}

// NewReporter creates a new Reporter.
func NewReporter() *cliReporter {
	return &cliReporter{
		errorPaths: make([]string, 0),
	}
}

func (r *cliReporter) LogProgress(message string) {
	log.Debug(fmt.Sprintf("[PROGRESS] %s", message))
}

func (r *cliReporter) LogSuccess(dbPath string) {
	atomic.AddUint32(&r.successCount, 1)
	atomic.AddUint32(&r.totalCount, 1)

	log.Info(fmt.Sprintf("[OK] %s", dbPath))
}

func (r *cliReporter) LogError(dbPath string, err error) {
	atomic.AddUint32(&r.errorCount, 1)
	atomic.AddUint32(&r.totalCount, 1)

	log.Error(fmt.Sprintf("[ERROR] %s: %v", dbPath, err))

	os.Exit(1)
}

func (r *cliReporter) PrintSummary() {
	log.Info("--- Verification Summary ---")
	log.Info(fmt.Sprintf("Total Databases Checked: %d", r.totalCount))
	log.Info(fmt.Sprintf("Successful Checks: %d", r.successCount))
	log.Info(fmt.Sprintf("Failed Checks: %d", r.errorCount))

	if len(r.errorPaths) > 0 {
		log.Error("Corrupted Databases:")
		for _, path := range r.errorPaths {
			log.Error(fmt.Sprintf("  - %s", path))
		}
	} else {
		log.Info("All databases are healthy!")
	}
}
