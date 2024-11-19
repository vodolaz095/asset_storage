package transport

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/vodolaz095/asset_storage/internal/service"
)

type WebServer struct {
	Srv            *http.Server
	Authentication *service.Authentication
	Assets         *service.Assets
	Logger         *log.Logger
}

func (s *WebServer) makeServer(ctx context.Context, addr string) {
	var err error
	handler := http.NewServeMux()
	handler.HandleFunc("/api/auth", s.login)
	handler.HandleFunc("/api/list", s.list)
	handler.HandleFunc("/api/upload-asset/", s.upload)
	handler.HandleFunc("/api/asset/", s.get)
	handler.HandleFunc("/api/delete/", s.deleteMyAsset)
	s.Srv = &http.Server{
		Addr:    addr,
		Handler: handler,
	}
	go func() {
		<-ctx.Done()
		s.Logger.Printf("Stopping web server...")
		cwto, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
		err = s.Srv.Shutdown(cwto)
		if err != nil {
			s.Logger.Printf("Error terminating web server: %s", err)
		} else {
			s.Logger.Printf("web server stopped...")
		}
	}()
}

func (s *WebServer) ListenHTTP(ctx context.Context, addr string) (err error) {
	s.makeServer(ctx, addr)
	return s.Srv.ListenAndServe()
}

func (s *WebServer) ListenHTTPS(ctx context.Context, addr string, pathToCert, pathToKey string) (err error) {
	s.makeServer(ctx, addr)
	return s.Srv.ListenAndServeTLS(pathToCert, pathToKey)
}
