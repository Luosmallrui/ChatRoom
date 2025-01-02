package dao

import (
	"chatroom/model"
	"gorm.io/gorm"
)

type Article struct {
	Repo[model.Article]
}

func NewArticle(db *gorm.DB) *Article {
	return &Article{Repo: NewRepo[model.Article](db)}
}
