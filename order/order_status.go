package order

const (
	StatusNew        = "NEW"
	StatusProcessing = "PROCESSING"
	StatusInvalid    = "INVALID"
	StatusProcessed  = "PROCESSED"
	StatusRegistered = "REGISTERED"
)

func GetNotProcessedStates() []string {
	return []string{StatusNew, StatusProcessing}
}
