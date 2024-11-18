package transport

import (
	"log"
	"net/http"

	"github.com/vodolaz095/asset_storage/internal/transport/http/dto"
)

func (s *WebServer) respond(writer http.ResponseWriter, status string, err error, code int) {
	writer.WriteHeader(code)
	r := dto.Response{
		Status: status,
	}
	if err != nil {
		r.Error = err.Error()
	}
	_, err = r.Write(writer)
	if err != nil {
		log.Printf("ошибка отправки запроса: %s", err)
		return
	}
	return

}
