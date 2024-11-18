package model

import (
	"time"
)

type Asset struct {
	Name      string
	UID       uint64
	Data      []byte
	CreatedAt time.Time
}
