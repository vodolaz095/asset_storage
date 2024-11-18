package model

import "time"

type User struct {
	ID           uint64
	Login        string
	PasswordHash string
	CreatedAt    time.Time
}
