package transport

import (
	"errors"
	"net/http"
	"strings"

	"github.com/vodolaz095/asset_storage/internal/model"
)

func (s *WebServer) get(writer http.ResponseWriter, request *http.Request) {
	var err error
	var user *model.User
	var bytesWritten int
	if request.Method != http.MethodGet {
		s.respond(writer, "", errors.New("only GET request allowed"), http.StatusBadRequest)
		return
	}

	assetID := strings.TrimPrefix(request.URL.Path, "/api/asset/")
	if assetID == "" {
		s.respond(writer, "", model.AssetNotFoundError, http.StatusNotFound)
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
	asset, err := s.Assets.LoadAssetForUser(request.Context(), assetID, user)
	if err != nil {
		if errors.Is(err, model.AssetNotFoundError) {
			s.Logger.Printf("Пользователь %s запросил несуществующие данные по ключу %s",
				user.Login, assetID)

			s.respond(writer, "", model.AssetNotFoundError, http.StatusNotFound)
			return
		}
		if errors.Is(err, model.ForbiddenError) {
			s.Logger.Printf("Пользователь %s запросил чужие данные по ключу %s",
				user.Login, assetID)

			s.respond(writer, "", model.AssetNotFoundError, http.StatusForbidden)
			return
		}

		s.respond(writer, "", err, http.StatusInternalServerError)
		return
	}
	writer.WriteHeader(http.StatusOK)
	bytesWritten, err = writer.Write(asset.Data)
	if err != nil {
		s.Logger.Printf("Ошибка отправки данных по ключу %s: %s", asset.Name, err)
		return
	}
	s.Logger.Printf("Пользователь %s запросил данные по ключу %s и получил %v байт",
		user.Login, assetID, bytesWritten)
}
