package reporter

import (
	"fmt"
	"os"
	"sync"
	"sync/atomic"
)

type cliReporter struct {
	mu           sync.Mutex
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
	r.mu.Lock()
	defer r.mu.Unlock()
	fmt.Printf("[PROGRESS] %s\n", message)
}

func (r *cliReporter) LogSuccess(dbPath string) {
	atomic.AddUint32(&r.successCount, 1)
	atomic.AddUint32(&r.totalCount, 1)

	r.mu.Lock()
	defer r.mu.Unlock()
	fmt.Printf("[OK] %s\n", dbPath)
}

func (r *cliReporter) LogError(dbPath string, err error) {
	atomic.AddUint32(&r.errorCount, 1)
	atomic.AddUint32(&r.totalCount, 1)

	r.mu.Lock()
	fmt.Printf("[ERROR] %s: %v\n", dbPath, err)
	r.mu.Unlock()

	os.Exit(1)
}

func (r *cliReporter) PrintSummary() {
	r.mu.Lock()
	defer r.mu.Unlock()

	fmt.Println("\n--- Verification Summary ---")
	fmt.Printf("Total Databases Checked: %d\n", r.totalCount)
	fmt.Printf("Successful Checks: %d\n", r.successCount)
	fmt.Printf("Failed Checks: %d\n", r.errorCount)

	if len(r.errorPaths) > 0 {
		fmt.Println("\nCorrupted Databases:")
		for _, path := range r.errorPaths {
			fmt.Printf("  - %s\n", path)
		}
	} else {
		fmt.Println("\nAll databases are healthy!")
	}
}
