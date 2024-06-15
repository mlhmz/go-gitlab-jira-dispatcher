package dispatcher

type Event struct {
	TicketNumber  string `json:"ticket_number"`
	StatusID      int    `json:"status_id"`
	ReviewerEmail string `json:"reviewer_email"`
}
