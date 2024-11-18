package repository

import (
	"context"

	"github.com/vodolaz095/asset_storage/internal/model"
)

type Repository interface {
	Ping(ctx context.Context) error
	Close(ctx context.Context) error
}

type User interface {
	Repository
	Login(ctx context.Context, login, password string) (session *model.Session, err error)
}

type Session interface {
	Repository
	Extract(ctx context.Context, sessionID string) (user *model.User, err error)
}

type Assets interface {
	Repository
	LoadAssetForUser(ctx context.Context, assetID string, user *model.User) (asset *model.Asset, err error)
	CreateAsset(ctx context.Context, author *model.User, name, content string) error
}
