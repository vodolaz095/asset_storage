package pg

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/vodolaz095/asset_storage/internal/model"
)

func makeHash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}

type User struct {
	baseRepo
}

func (ur *User) Login(ctx context.Context, username, password string) (session *model.Session, err error) {
	var user model.User
	var sessionID string
	err = ur.Conn.
		QueryRow(ctx, "SELECT id, login, password_hash, created_at FROM users WHERE login=$1", username).
		Scan(&user.ID, &user.Login, &user.PasswordHash, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.WrongUsernameOrPasswordError
		}
		// unexpected error
		return nil, err
	}
	if user.PasswordHash != makeHash(password) {
		return nil, model.WrongUsernameOrPasswordError
	}
	err = ur.Conn.QueryRow(ctx, "INSERT INTO sessions (uid) values ($1) returning id", user.ID).
		Scan(&sessionID)
	if err != nil {
		return nil, err
	}
	return &model.Session{
		ID:        sessionID,
		UID:       user.ID,
		CreatedAt: time.Now(),
	}, nil
}
