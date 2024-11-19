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

func (a *Asset) ListAll(ctx context.Context) (list []model.ListOfAssets, err error) {
	query := `select assets.name as "name",
u.login as "author",
length(assets.data) as "size",
assets.created_at as "created_at"
from assets
left join users u on assets.uid = u.id
order by assets.name asc;
`
	var item model.ListOfAssets
	rows, err := a.Conn.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		err = rows.Scan(&item.Name, &item.Author, &item.Size, &item.CreatedAt)
		if err != nil {
			return nil, err
		}
		list = append(list, item)
	}
	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return list, nil
}
