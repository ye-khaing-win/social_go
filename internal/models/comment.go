package models

type Comment struct {
	ID        string `json:"id"`
	PostID    string `json:"post_id"`
	UserID    string `json:"user_id"`
	Content   string `json:"content_id"`
	CreatedAt string `json:"created_at"`
	User      User   `json:"user"`
}
