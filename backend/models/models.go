package models

import (
	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	gorm.Model
	// uniqueIndex: 用户名必须唯一
	// varchar(100): 指定数据库字段类型和长度
	Username string `gorm:"type:varchar(100);uniqueIndex;not null" json:"username"`

	// json:"-": 序列化 JSON 时忽略该字段，防止密码泄露给前端
	Password string `gorm:"type:varchar(255);not null" json:"-"`

	// Role 角色: admin / user
	Role string `gorm:"type:varchar(20);default:'user'" json:"role"`

	// 关联: 一个用户可以有多篇文章
	Posts []Post `gorm:"foreignKey:AuthorID" json:"posts,omitempty"`
}

// Category 分类模型
type Category struct {
	gorm.Model
	// uniqueIndex: 分类名唯一
	Name string `gorm:"type:varchar(100);uniqueIndex;not null" json:"name"`

	// 关联: 一个分类下有多篇文章
	Posts []Post `gorm:"foreignKey:CategoryID" json:"posts,omitempty"`
}

// Post 文章模型
type Post struct {
	gorm.Model
	// index: 标题虽然不需要唯一，但可能会用于搜索，加索引可以加快模糊查询（视具体数据库优化而定，这里先保留基础定义）
	Title string `gorm:"type:varchar(200);not null" json:"title"`

	// type:text: 文章内容通常较长，使用 text 类型
	Content string `gorm:"type:text;not null" json:"content"`

	// 摘要，可为空
	Summary string `gorm:"type:varchar(500)" json:"summary"`

	// 浏览量，默认为 0
	ViewCount uint64 `gorm:"default:0" json:"view_count"`

	// 外键: 关联分类
	// index: 经常需要查询"某个分类下的文章"，所以必须加索引
	CategoryID uint     `gorm:"index;not null" json:"category_id"`
	Category   Category `json:"category"`

	// 外键: 关联作者
	// index: 经常需要查询"某个用户的文章"
	AuthorID uint `gorm:"index;not null" json:"author_id"`
	Author   User `json:"-"` // 在文章详情中通常不需要嵌套返回完整的作者所有信息，视情况调整
}
