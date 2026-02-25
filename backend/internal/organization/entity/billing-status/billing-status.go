package billing_status

type BillingStatus string

const (
	BillingStatusActive  BillingStatus = "ACTIVE"
	BillingStatusPastDue BillingStatus = "PAST_DUE"
	BillingStatusBlocked BillingStatus = "BLOCKED"
)
