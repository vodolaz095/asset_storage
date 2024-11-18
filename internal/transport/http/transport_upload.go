package transport

import (
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/vodolaz095/asset_storage/internal/model"
)

func (s *WebServer) upload(writer http.ResponseWriter, request *http.Request) {
	var err error
	var user *model.User
	var body []byte

	if request.Method != http.MethodPost {
		s.respond(writer, "", errors.New("only POST request allowed"), http.StatusBadRequest)
		return
	}

	assetID := strings.TrimPrefix(request.URL.Path, "/api/upload-asset/")
	if assetID == "" {
		s.respond(writer, "", errors.New("asset's name is not provided"), http.StatusBadRequest)
		return
	}

	bearer := request.Header.Get("Authorization")
	if !strings.HasPrefix(bearer, "Bearer") {
		s.respond(writer, "",
			errors.New("authorization header with Bearer strategy is required"),
			http.StatusBadRequest,
		)
	}
	bearer = strings.TrimPrefix(bearer, "Bearer ")
	user, err = s.Authentication.Extract(request.Context(), bearer)
	if err != nil {
		if errors.Is(err, model.SessionNotFoundError) {
			s.respond(writer, "", err, http.StatusUnauthorized)
			return
		}
		s.respond(writer, "", err, http.StatusInternalServerError)
		return
	}
	s.Logger.Printf("Пользователь %s восстановлен из сессии", user.Login)
	if request.Body != nil {
		body, err = io.ReadAll(request.Body)
		if err != nil {
			s.respond(writer, "", err, http.StatusInternalServerError)
			return
		}
		defer request.Body.Close()
		s.Logger.Printf("Пользователь %s пытается загрузить %v байт в ключ %s",
			user.Login, len(body), assetID)

		err = s.Assets.CreateAsset(request.Context(), user, assetID, string(body))
		if err != nil {
			s.respond(writer, "", err, http.StatusBadRequest)
			return
		}

		s.Logger.Printf("Пользователь %s успешно загрузил %v байт в ключ %s",
			user.Login, len(body), assetID)

		// https://restapitutorial.ru/lessons/httpmethods/
		// 201 (Created), заголовок 'Location' ссылается на /customers/{id}, где ID - идентификатор нового экземпляра.
		writer.Header().Add("Location", "/api/asset/"+assetID)
		s.respond(writer, "ok", nil, http.StatusCreated)
		return
	}
	s.respond(writer, "",
		errors.New("body is empty"),
		http.StatusBadRequest,
	)
}
