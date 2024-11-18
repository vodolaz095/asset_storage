package model

import "time"

type Asset struct {
	Name      string
	UID       uint64
	Data      string
	CreatedAt time.Time
}
