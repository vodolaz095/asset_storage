package repository

import (
	"context"
	"errors"
	"testing"

	"github.com/vodolaz095/asset_storage/internal/model"
)

// TestRepository используется для black-box тестирования репозиториев
func TestRepository(t *testing.T, user User, session Session, assets Assets) {
	const goodUsername, goodPassword = "alice", "secret"
	const badUsername, badPassword = "not_alice", "not_secret"
	const remoteAddr = "[::1]:43466"

	t.Parallel()

	t.Run("pinging", func(tt *testing.T) {
		var err error
		ctx := context.TODO()
		err = user.Ping(ctx)
		if err != nil {
			tt.Errorf("error pinging user: %s", err)
		}
		err = session.Ping(ctx)
		if err != nil {
			tt.Errorf("error pinging user: %s", err)
		}
		err = session.Ping(ctx)
		if err != nil {
			tt.Errorf("error pinging session: %s", err)
		}
		err = assets.Ping(ctx)
		if err != nil {
			tt.Errorf("error pinging assets: %s", err)
		}
	})

	t.Run("user.Login", func(tt *testing.T) {
		tt.Run("good login and password", func(ttt *testing.T) {
			ctx := context.TODO()
			ses, err := user.Login(ctx, goodUsername, goodPassword, remoteAddr)
			if err != nil {
				ttt.Errorf("error logging in as %s %s: %s", goodUsername, goodPassword, err)
				return
			}
			ttt.Logf("Session %v created", ses)
		})
		tt.Run("good login, bad password", func(ttt *testing.T) {
			ctx := context.TODO()
			ses, err := user.Login(ctx, goodUsername, badPassword, remoteAddr)
			if err != nil {
				if errors.Is(err, model.WrongUsernameOrPasswordError) {
					ttt.Logf("failed as expected")
					return
				}
				ttt.Errorf("error logging in as %s %s: %s", goodUsername, goodPassword, err)
				return
			}
			ttt.Errorf("not failed as expected! session: %v", ses)
		})
		tt.Run("bad login, good password", func(ttt *testing.T) {
			ctx := context.TODO()
			ses, err := user.Login(ctx, badUsername, goodPassword, remoteAddr)
			if err != nil {
				if errors.Is(err, model.WrongUsernameOrPasswordError) {
					ttt.Logf("failed as expected")
					return
				}
				ttt.Errorf("error logging in as %s %s: %s", badUsername, goodPassword, err)
				return
			}
			ttt.Errorf("not failed as expected! session: %v", ses)
		})
		tt.Run("bad login and bad password", func(ttt *testing.T) {
			ctx := context.TODO()
			ses, err := user.Login(ctx, badUsername, badPassword, remoteAddr)
			if err != nil {
				if errors.Is(err, model.WrongUsernameOrPasswordError) {
					ttt.Logf("failed as expected")
					return
				}
				ttt.Errorf("error logging in as %s %s: %s", badUsername, badPassword, err)
				return
			}
			ttt.Errorf("not failed as expected! session: %v", ses)
		})
	})

	t.Run("user.createAsset", func(tt *testing.T) {
		tt.Skipf("not implemented yet")
	})
	t.Run("user.loadAsset", func(tt *testing.T) {
		tt.Skipf("not implemented yet")
	})
}
