package pg_test

import (
	"context"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/vodolaz095/asset_storage/internal/repository"
	"github.com/vodolaz095/asset_storage/internal/repository/pg"
)

func TestRepo(t *testing.T) {
	conn, err := pgx.Connect(context.TODO(), "user=assets password=assets dbname=assets host=localhost port=5432 sslmode=disable")
	if err != nil {
		t.Errorf("error dialing postgres: %s", err)
		return
	}
	var user pg.User
	user.Conn = conn

	var session pg.Session
	session.Conn = conn

	var asset pg.Asset
	asset.Conn = conn

	repository.TestRepository(t, &user, &session, &asset)
}
