package transport

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/vodolaz095/asset_storage/internal/model"
	"github.com/vodolaz095/asset_storage/internal/transport/http/dto"
)

func (s *WebServer) login(writer http.ResponseWriter, request *http.Request) {
	var err error
	var form dto.LoginForm
	if request.Method != http.MethodPost {
		s.respond(writer,
			"error",
			errors.New("HTTP POST request expected"),
			http.StatusBadRequest,
		)
		return
	}
	reader := json.NewDecoder(request.Body)
	defer request.Body.Close()

	err = reader.Decode(&form)
	if err != nil {
		s.respond(writer,
			"error",
			err,
			http.StatusBadRequest,
		)
		return
	}

	session, err := s.Authentication.Login(request.Context(), form.Username, form.Password)
	if err != nil {
		if errors.Is(err, model.WrongUsernameOrPasswordError) {
			writer.WriteHeader(http.StatusUnauthorized)
		} else {
			writer.WriteHeader(http.StatusInternalServerError)
		}
		r := dto.Response{
			Status: "error",
			Error:  err.Error(),
		}
		_, err = r.Write(writer)
		if err != nil {
			s.Logger.Printf("ошибка отправки запроса: %s", err)
			return
		}
		return
	}
	loginResponse := dto.LoginResponse{Token: session.ID}

	writer.WriteHeader(http.StatusOK)
	_, err = loginResponse.Write(writer)
	if err != nil {
		s.Logger.Printf("ошибка отправки запроса: %s", err)
		return
	}
}
