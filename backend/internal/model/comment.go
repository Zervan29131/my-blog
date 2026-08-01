package model

import "time"

const (
	CommentStatusPending  = "pending"
	CommentStatusApproved = "approved"
	CommentStatusRejected = "rejected"
)

type Comment struct {
	ID        uint64    `gorm:"primaryKey;autoIncrement" json:"id"`
	ArticleID uint64    `gorm:"not null;index" json:"article_id"`
	Article   Article   `gorm:"foreignKey:ArticleID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE" json:"-"`
	Nickname  string    `gorm:"type:varchar(50);not null" json:"nickname"`
	Email     *string   `gorm:"type:varchar(255)" json:"-"`
	Content   string    `gorm:"type:varchar(1000);not null" json:"content"`
	Status    string    `gorm:"type:varchar(20);not null;default:pending;index;check:comments_status_check,status IN ('pending','approved','rejected')" json:"status"`
	CreatedAt time.Time `gorm:"not null;index" json:"created_at"`
	UpdatedAt time.Time `gorm:"not null" json:"updated_at"`
}

func (Comment) TableName() string {
	return "comments"
}
