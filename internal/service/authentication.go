package service

import (
	"context"
	"log"

	"github.com/vodolaz095/asset_storage/internal/model"
	"github.com/vodolaz095/asset_storage/internal/repository"
)

type Authentication struct {
	UserRepo    repository.User
	SessionRepo repository.Session
	Logger      *log.Logger
}

func (a *Authentication) Login(ctx context.Context, username, password string) (*model.Session, error) {
	a.Logger.Printf("Пользователь %s пытается авторизоваться...", username)
	session, err := a.UserRepo.Login(ctx, username, password)
	if err != nil {
		a.Logger.Printf("Пользователь %s не смог авторизоваться: %s.", username, err)
		return nil, err
	}
	a.Logger.Printf("Пользователь %s создал сессию.", username)
	return session, nil
}

func (a *Authentication) Extract(ctx context.Context, sessionID string) (*model.User, error) {
	user, err := a.SessionRepo.Extract(ctx, sessionID)
	if err != nil {
		a.Logger.Printf("Ошибка загрузки профиля пользователя из сессии: %s", err)
		return nil, err
	}
	a.Logger.Printf("Сессия востановлена для пользователя: %s", user.Login)
	return user, nil
}
