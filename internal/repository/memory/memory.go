package memory

import (
	"context"
	"fmt"

	"github.com/vodolaz095/asset_storage/internal/model"
)

// это мок, который может работать без базы данных
// возможно, я его когда-нибудь и доделаю...

type Repository struct {
	Users    map[string]model.User
	Sessions map[string]model.Session
	Assets   map[string]model.Asset
	Active   bool
}

func (r *Repository) Ping(ctx context.Context) error {
	if r.Active {
		return nil
	}
	return fmt.Errorf("ping error")
}

func (r *Repository) Close(ctx context.Context) error {
	return fmt.Errorf("not implemented")
}

func (r *Repository) Login(ctx context.Context, username, password, address string) (session *model.Session, err error) {
	return nil, fmt.Errorf("not implemented")
}

func (r *Repository) Extract(ctx context.Context, sessionID string) (user *model.User, err error) {
	return nil, fmt.Errorf("not implemented")
}

func (r *Repository) LoadAssetForUser(ctx context.Context, assetID string, user *model.User) (asset *model.Asset, err error) {
	return nil, fmt.Errorf("not implemented")
}

func (r *Repository) CreateAsset(ctx context.Context, author *model.User, name, content string) error {
	return fmt.Errorf("not implemented")
}

func (r *Repository) ListAll(ctx context.Context) ([]model.ListOfAssets, error) {
	return nil, fmt.Errorf("not implemented")
}

func (r *Repository) DeleteMyAsset(ctx context.Context, author *model.User, assetID string) (err error) {
	return fmt.Errorf("not implemented")
}
