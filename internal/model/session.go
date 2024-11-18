package model

import "time"

type Session struct {
	ID        string
	UID       uint64
	CreatedAt time.Time
}
