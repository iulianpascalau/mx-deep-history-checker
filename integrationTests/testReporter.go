package integrationTests

import "sync"

type testReporter struct {
	mut          sync.Mutex
	progressLogs []string
	successLogs  []string
	errorLogs    []string
}

func NewTestReporter() *testReporter {
	return &testReporter{
		progressLogs: make([]string, 0, 1000),
		successLogs:  make([]string, 0, 1000),
		errorLogs:    make([]string, 0, 1000),
	}
}

func (rep *testReporter) LogProgress(message string) {
	rep.mut.Lock()
	rep.progressLogs = append(rep.progressLogs, message)
	rep.mut.Unlock()
}

func (rep *testReporter) LogSuccess(dbPath string) {
	rep.mut.Lock()
	rep.successLogs = append(rep.successLogs, dbPath)
	rep.mut.Unlock()
}

func (rep *testReporter) LogError(dbPath string, _ error) {
	rep.mut.Lock()
	rep.errorLogs = append(rep.errorLogs, dbPath)
	rep.mut.Unlock()
}

func (rep *testReporter) PrintSummary() {
}

func (rep *testReporter) GetProgressLogs() []string {
	rep.mut.Lock()
	defer rep.mut.Unlock()

	return rep.progressLogs
}

func (rep *testReporter) GetSuccessLogs() []string {
	rep.mut.Lock()
	defer rep.mut.Unlock()

	return rep.successLogs
}

func (rep *testReporter) GetErrorLogs() []string {
	rep.mut.Lock()
	defer rep.mut.Unlock()

	return rep.errorLogs
}
