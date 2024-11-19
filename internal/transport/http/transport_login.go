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
	// TODO - в случае, если используется reverse proxy, то в request.RemoteAddr
	// будет адрес прокси сервера, а адрес клиента будет в заголовке
	// `X-Forwarded-For` или в `CF-Connecting-IP` для Cloudflare.
	// Также бы было неплохо указать список IP адресов reverse-proxy
	// запросам с которых мы доверяем см. https://pkg.go.dev/github.com/gin-gonic/gin#Engine.SetTrustedProxies
	// мне кажется, это довольно сложно корректно реализовать без готовой библиотеки.
	s.Logger.Printf("Поступил запрос на авторизацию с %s", request.RemoteAddr)
	session, err := s.Authentication.Login(request.Context(),
		form.Username, form.Password, request.RemoteAddr,
	)
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
