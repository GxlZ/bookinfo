package models

import (
	"time"
)

type Comments struct {
	ID        uint       `gorm:"primary_key"`
	BookId    uint       `sql:"index"`
	Content   string     `gorm:"type:text"` //书名
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`
}

var comments = []Comments{
	{ID: 1, BookId: 1, Content: "测试评论1"},
	{ID: 2, BookId: 1, Content: "测试评论2"},
	{ID: 3, BookId: 1, Content: "测试评论3"},
}
