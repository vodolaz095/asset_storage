package service

import (
	"context"
	"log"

	"github.com/vodolaz095/asset_storage/internal/model"
	"github.com/vodolaz095/asset_storage/internal/repository"
)

type Assets struct {
	AssetsRepo repository.Assets
	Logger     *log.Logger
}

func (a *Assets) LoadAssetForUser(ctx context.Context, assetsID string, user *model.User) (*model.Asset, error) {
	a.Logger.Printf("Пользователь %s пытается получить данные по ключу %s", user.Login, assetsID)
	asset, err := a.AssetsRepo.LoadAssetForUser(ctx, assetsID, user)
	if err != nil {
		a.Logger.Printf("Пользователь %s не смог получить данные по ключу %s: %s",
			user.Login, assetsID, err)

		return nil, err
	}
	a.Logger.Printf("Пользователь %s получил данные по ключу %s (%v байт)",
		user.Login, assetsID, len(asset.Data))

	return asset, nil
}

func (a *Assets) CreateAsset(ctx context.Context, author *model.User, assetsID, data string) error {
	a.Logger.Printf("Пользователь %s пытается создать данные (%v байт) по ключу %s",
		author.Login, len(data), assetsID)
	err := a.AssetsRepo.CreateAsset(ctx, author, assetsID, data)
	if err != nil {
		a.Logger.Printf("Пользователь %s не смог создать данные по ключу %s: %s",
			author.Login, assetsID, err)

		return err
	}
	a.Logger.Printf("Пользователь %s создал данные по ключу %s (%v байт)",
		author.Login, assetsID, len(data))

	return nil
}

func (a *Assets) ListAll(ctx context.Context) ([]model.ListOfAssets, error) {
	a.Logger.Printf("Загружаем список всех объектов")
	return a.AssetsRepo.ListAll(ctx)
}
