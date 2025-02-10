package contracts

type Notification struct {
	UserId      int64  `json:"user_id"`
	Date        string `json:"date"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Type        string `json:"type"`
}
