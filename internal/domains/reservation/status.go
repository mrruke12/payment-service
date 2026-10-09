package reservation

type Status string

const (
	StatusActive    Status = "active"
	StatusExpired   Status = "expired"
	StatusCancelled Status = "cancelled"
	StatusConfirmed Status = "confirmed"
)
