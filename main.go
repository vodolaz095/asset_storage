package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/jackc/pgx/v5"
	"github.com/vodolaz095/asset_storage/config"
	"github.com/vodolaz095/asset_storage/internal/repository/pg"
	"github.com/vodolaz095/asset_storage/internal/service"
	transport "github.com/vodolaz095/asset_storage/internal/transport/http"
)

func main() {
	wg := sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	logger := log.New(os.Stdout, "ASSET: ", log.Lshortfile|log.Lmsgprefix|log.Ltime)
	config.Load()

	// handle signals
	wg.Add(1)
	sigc := make(chan os.Signal, 1)
	signal.Notify(sigc, syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGABRT,
	)
	go func() {
		s := <-sigc
		logger.Printf("Signal %s is received", s.String())
		wg.Done()
		cancel()
	}()

	wg.Add(1)
	conn, err := pgx.Connect(ctx, config.DSN)
	if err != nil {
		log.Printf("ошибка соединения с postgres: %s", err)
		return
	}
	var user pg.User
	var session pg.Session
	var asset pg.Asset

	user.Conn, session.Conn, asset.Conn = conn, conn, conn

	authenticationService := service.Authentication{
		UserRepo:    &user,
		SessionRepo: &session,
		Logger:      logger,
	}

	assetService := service.Assets{
		AssetsRepo: &asset,
		Logger:     logger,
	}

	server := transport.WebServer{
		Authentication: &authenticationService,
		Assets:         &assetService,
		Logger:         logger,
	}

	go func() {
		// надо бы закрыть и другие репозитории, но так как они используют один и тот же
		// пул соединений с одной и той же базой, то закрыть репозиторий пользователей будет в
		// данном случае будет достаточно
		<-ctx.Done()
		err = user.Close(context.Background())
		if err != nil {
			logger.Printf("Ошибка завершения репозитория User: %s", err)
		}
		wg.Done()
	}()

	logger.Printf("Запускаем приложение на %s:%s...",
		config.Address, config.Port)
	err = server.ListenHTTP(ctx, config.Address+":"+config.Port)
	if err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			logger.Printf("ошибка запуска вэб-сервера на %s:%s - %s",
				config.Address, config.Port, err)
		}
	}
	wg.Wait()
	logger.Printf("Приложение завершено успешно!")
}
