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

func (s *WebServer) ListenHTTP(ctx context.Context, addr string) (err error) {
	handler := http.NewServeMux()
	handler.HandleFunc("/api/auth", s.login)
	handler.HandleFunc("/api/upload-asset/", s.upload)
	handler.HandleFunc("/api/asset/", s.get)
	s.Srv = &http.Server{
		Addr:    addr,
		Handler: handler,
	}
	go func() {
		<-ctx.Done()
		s.Logger.Printf("Stopping http server...")
		cwto, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
		err = s.Srv.Shutdown(cwto)
		if err != nil {
			s.Logger.Printf("Error terminating HTTP server: %s", err)
		} else {
			s.Logger.Printf("HTTP server stopped...")
		}
	}()
	return s.Srv.ListenAndServe()
}

func (s *WebServer) ListenHTTPS(ctx context.Context, addr string, pathToCert, pathToKey string) (err error) {
	handler := http.NewServeMux()
	handler.HandleFunc("/api/auth", s.login)
	handler.HandleFunc("/api/upload-asset/", s.upload)
	handler.HandleFunc("/api/asset/", s.get)
	s.Srv = &http.Server{
		Addr:    addr,
		Handler: handler,
	}
	go func() {
		<-ctx.Done()
		s.Logger.Printf("Stopping https server...")
		cwto, cancel := context.WithTimeout(ctx, 10*time.Second)
		defer cancel()
		err = s.Srv.Shutdown(cwto)
		if err != nil {
			s.Logger.Printf("Error terminating HTTPS server: %s", err)
		} else {
			s.Logger.Printf("HTTPS server stopped...")
		}
	}()
	return s.Srv.ListenAndServeTLS(pathToCert, pathToKey)
}
