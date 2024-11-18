package dto

import (
	"encoding/json"
	"io"
)

type Response struct {
	Status string `json:"status,omitempty"`
	Error  string `json:"error,omitempty"`
}

func (r *Response) Write(to io.Writer) (bytesWritten int, err error) {
	data, err := json.MarshalIndent(r, "", " ")
	if err != nil {
		return 0, err
	}
	return to.Write(data)
}
