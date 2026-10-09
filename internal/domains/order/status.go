package order

type Status string

const (
	StatusReserved  Status = "reserved"
	StatusPaid      Status = "paid"
	StatusCancelled Status = "cancelled"
	StatusFailed    Status = "failed"
	StatusRefunded  Status = "refunded"
)
