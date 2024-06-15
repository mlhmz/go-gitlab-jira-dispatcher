package store

type Transitions struct {
	ReadyForReview  int
	InReview        int
	ReviewOK        int
	ReviewNotOK     int
	DevelopmentDone int
}
