package jirav2

type RestClient struct {
}

func NewRestClient() *RestClient {
	return &RestClient{}
}

func (r *RestClient) TransitionIssue(ticketNumber string, statusID int, reviewerEmail string) error {
	return nil
}
