package factory

type Reporter interface {
	LogProgress(message string)
	LogSuccess(dbPath string)
	LogError(dbPath string, err error)
	PrintSummary()
}
