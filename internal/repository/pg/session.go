package pg

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/vodolaz095/asset_storage/internal/model"
)

type Session struct {
	baseRepo
}

func (s *Session) Extract(ctx context.Context, sessionID string) (user *model.User, err error) {
	var sessionFound model.Session
	var userFound model.User

	err = s.Conn.QueryRow(ctx, "SELECT * FROM sessions WHERE id=$1", sessionID).Scan(&sessionFound)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.SessionNotFoundError
		}
		return nil, err
	}
	err = s.Conn.QueryRow(ctx, "SELECT * FROM users WHERE id=$1", sessionFound.UID).Scan(&userFound)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, model.SessionNotFoundError
		}
		return nil, err
	}
	return &userFound, nil
}
