package gitlab

type MergeRequestEvent struct {
	ObjectKind       string           `json:"object_kind"`
	ObjectAttributes ObjectAttributes `json:"object_attributes"`
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

type User struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Username  string `json:"username"`
	AvatarURL string `json:"avatar_url"`
	Email     string `json:"email"`
}
