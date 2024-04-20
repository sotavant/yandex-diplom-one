package order

const (
	STATUS_NEW        = "NEW"
	STATUS_PROCESSING = "PROCESSING"
	STATUS_INVALID    = "INVALID"
	STATUS_PROCESSED  = "PROCESSED"
	STATUS_REGISTERED = "REGISTERED"
)

func GetNotProcessedStates() []string {
	return []string{STATUS_NEW, STATUS_PROCESSING}
}
