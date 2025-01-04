package pending_approved_cancelled_completed_enum

var PendingApprovedCancelledCompleted = newPendingApprovedCancelledCompleted()

func newPendingApprovedCancelledCompleted() *pendingApprovedCancelledCompleted {
	return &pendingApprovedCancelledCompleted{
		PENDING:   "PENDING",
		APPROVED:  "APPROVED",
		CANCELLED: "CANCELLED",
		COMPLETED: "COMPLETED",
	}
}

type pendingApprovedCancelledCompleted struct {
	PENDING   string
	APPROVED  string
	CANCELLED string
	COMPLETED string
}
