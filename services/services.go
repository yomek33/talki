package services

import (
	"github.com/yomek33/talki/stores"
)

type Services struct {
	UserService    UserService
	ArticleService ArticleService
}

func NewServices(s *stores.Stores) *Services {
	return &Services{
		UserService:    &userService{store: s.UserStore},
		ArticleService: &articleService{store: s.ArticleStore},
	}
}
