package model

import (
	"errors"
)

var (
	WrongUsernameOrPasswordError = errors.New("invalid login/password")
	SessionExpiredError          = errors.New("session expired")
	SessionNotFoundError         = errors.New("session not found")
	AssetNotFoundError           = errors.New("asset not found")
	ForbiddenError               = errors.New("access restricted")
)
