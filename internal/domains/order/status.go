package order

type Status string

const (
	StatusPendingPayment Status = "pending_payment" // created but not paid yet
	StatusPaid           Status = "paid"            // successfully paid
	StatusCancelled      Status = "cancelled"       // cancelled before payment
	StatusRefunded       Status = "refunded"        // cancelled after payment
)
