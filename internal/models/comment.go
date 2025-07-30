package models

type Comment struct {
	ID        int64  `json:"id"`
	PostID    int64  `json:"post_id"`
	UserID    int64  `json:"user_id"`
	Content   string `json:"content_id"`
	CreatedAt string `json:"created_at"`
	User      User   `json:"user"`
}
