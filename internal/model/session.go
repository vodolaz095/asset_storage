package model

import (
	"net/netip"
	"time"
)

type Session struct {
	ID        string
	UID       uint64
	ClientIP  netip.Addr
	CreatedAt time.Time
}
