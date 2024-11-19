package transport

import (
	"encoding/json"
	"errors"
	"net/http"
)

func (s *WebServer) list(writer http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodGet {
		s.respond(writer, "", errors.New("only GET request allowed"), http.StatusBadRequest)
		return
	}
	assets, err := s.Assets.ListAll(request.Context())
	if err != nil {
		s.respond(writer, "error", err, http.StatusInternalServerError)
		return
	}
	writer.WriteHeader(http.StatusOK)
	enc := json.NewEncoder(writer)
	enc.SetIndent("", " ")
	err = enc.Encode(assets)
	if err != nil {
		s.Logger.Printf("ошибка кодировки в json: %s", err)
		return
	}
	s.Logger.Printf("Отправлен список из %v объектов", len(assets))
}
