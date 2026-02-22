package pos_sync_status

type PosSyncStatus string

const (
	PosSyncStatusNotSent PosSyncStatus = "NOT_SENT"
	PosSyncStatusSent    PosSyncStatus = "SENT"
	PosSyncStatusError   PosSyncStatus = "ERROR"
)
