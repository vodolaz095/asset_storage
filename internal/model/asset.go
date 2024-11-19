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

type ListOfAssets struct {
	Name      string    `json:"name"`
	Author    string    `json:"author"`
	Size      int       `json:"size"`
	CreatedAt time.Time `json:"created_at"`
}
