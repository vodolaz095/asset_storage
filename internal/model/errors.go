package model

import (
	"errors"
)

var (
	WrongUsernameOrPasswordError = errors.New("invalid login/password")
	SessionNotFoundError         = errors.New("session not found")
	AssetNotFoundError           = errors.New("asset not found")
	RestrictedError              = errors.New("access restricted")
)
