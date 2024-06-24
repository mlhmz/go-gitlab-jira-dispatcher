package jira

type TransitionPayload struct {
	Transition Transition `json:"transition"`
}

type Transition struct {
	ID string `json:"id"`
}
