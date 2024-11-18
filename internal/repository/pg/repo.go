package pg

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type baseRepo struct {
	Conn *pgx.Conn
}

func (br *baseRepo) Ping(ctx context.Context) error {
	return br.Conn.Ping(ctx)
}

func (br *baseRepo) Close(ctx context.Context) error {
	return br.Conn.Close(ctx)
}
