package gitlab

type MergeRequestEvent struct {
	ObjectKind       string           `json:"object_kind"`
	Repository       Repository       `json:"repository"`
	ObjectAttributes ObjectAttributes `json:"object_attributes"`
	Changes          Changes          `json:"changes"`
}

type ObjectAttributes struct {
	ID     int    `json:"id"`
	Action string `json:"action"`
	Title  string `json:"title"`
}

type Changes struct {
	Reviewers UserChange `json:"reviewers"`
}

type UserChange struct {
	Previous []User `json:"previous"`
	Current  []User `json:"current"`
}

type Repository struct {
	Name string `json:"name"`
}

type User struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Username  string `json:"username"`
	AvatarURL string `json:"avatar_url"`
	Email     string `json:"email"`
}
