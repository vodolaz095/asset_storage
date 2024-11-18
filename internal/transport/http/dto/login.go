package dto

import (
	"encoding/json"
	"io"
)

type LoginForm struct {
	Username string `json:"login"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func (r *LoginResponse) Write(to io.Writer) (bytesWritten int, err error) {
	data, err := json.Marshal(r)
	if err != nil {
		return 0, err
	}
	return to.Write(data)
}
