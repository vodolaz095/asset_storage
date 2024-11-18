package pg

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/vodolaz095/asset_storage/internal/model"
)

type Asset struct {
	baseRepo
}

func (a *Asset) LoadAssetForUser(ctx context.Context, assetID string, user *model.User) (*model.Asset, error) {
	var asset model.Asset
	err := a.Conn.
		QueryRow(ctx, "SELECT name, uid, data, created_at FROM assets WHERE name=$1 and uid=$2", assetID, user.ID).
		Scan(&asset.Name, &asset.UID, &asset.Data, &asset.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.AssetNotFoundError
		}
		return nil, err
	}
	return &asset, nil
}

func (a *Asset) CreateAsset(ctx context.Context, author *model.User, name, content string) (err error) {
	_, err = a.Conn.Exec(ctx, "INSERT INTO assets(name,uid,data) VALUES ($1,$2,$3)",
		name, author.ID, content,
	)
	return err
}
