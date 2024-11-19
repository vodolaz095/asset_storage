package transport

import (
	"errors"
	"net/http"
	"strings"

	"github.com/vodolaz095/asset_storage/internal/model"
)

func (s *WebServer) deleteMyAsset(writer http.ResponseWriter, request *http.Request) {
	var err error
	var user *model.User
	if request.Method != http.MethodDelete {
		s.respond(writer, "", errors.New("only DELETE request allowed"), http.StatusBadRequest)
		return
	}

	assetID := strings.TrimPrefix(request.URL.Path, "/api/delete/")
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

	err = s.Assets.DeleteMyAsset(request.Context(), user, assetID)
	if err != nil {
		if errors.Is(err, model.AssetNotFoundError) {
			s.respond(writer, "asset not found", err, http.StatusNotFound)
			return
		}
		if errors.Is(err, model.ForbiddenError) {
			s.respond(writer, "you cannot delete asset belonging to different user", err, http.StatusForbidden)
			return
		}
		s.respond(writer, "", err, http.StatusInternalServerError)
		return
	}
	writer.WriteHeader(http.StatusNoContent)
	return
}
