package main

import (
	"context"
	"fmt"
	"github.com/julienschmidt/httprouter"
	postgresql2 "main/internal/adapters/postgresql"
	"main/internal/config"
	v1 "main/internal/controller/http/v1"
	"main/internal/domain/service"
	product_usecase "main/internal/domain/usecase/product"
	user_usecase "main/internal/domain/usecase/user"
	tokenManager2 "main/internal/tokenManager"
	postgresql "main/pkg/client/postgresql"
	"main/pkg/logging"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"
)

func main() {
	logger := logging.GetLogger()
	logger.Infoln("create router")
	router := httprouter.New()

	logger.Info("get config")
	cfg := config.GetConfig()

	logger.Info("create clients")
	cfgPostgre := cfg.Storage
	postgreClient, err := postgresql.NewClient(context.Background(), 5, cfgPostgre)
	if err != nil {
		panic(err)
	}

	userStorage := postgresql2.NewUserStorage(postgreClient)
	userService := service.NewUserService(userStorage)
	tokenManager := tokenManager2.NewTokenManager(cfg.JWTConfig.Secret, userStorage)
	userUsecase := user_usecase.NewUserUsecase(userService, tokenManager)

	productStorage := postgresql2.NewProductStorage(postgreClient)
	productService := service.NewProductService(productStorage)
	productUsecase := product_usecase.NewProductUsecase(productService, tokenManager)

	logger.Info("register handlers")
	handler := v1.NewUserHandler(userUsecase, tokenManager, logger)
	productHandler := v1.NewProductHandler(productUsecase, tokenManager, logger)
	productHandler.Register(router)
	handler.Register(router)

	start(router, cfg, logger)
}

func start(router *httprouter.Router, cfg *config.Config, logger *logging.Logger) {
	logger.Infoln("start application")

	var listener net.Listener
	var listenErr error

	if cfg.Listen.Type == "sock" {
		logger.Info("detect app path")
		appDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
		if err != nil {
			logger.Fatal(err)
		}
		logger.Info("create socket")
		socketPath := path.Join(appDir, "app.sock")

		logger.Info("listen unix socket")
		listener, listenErr = net.Listen("unix", socketPath)
		logger.Infof("server is listening unix socket: %s", socketPath)
	} else {
		logger.Info("listen tcp")
		listener, listenErr = net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.Listen.BindIP, cfg.Listen.Port))
		logger.Infof("server is listening port %s:%s", cfg.Listen.BindIP, cfg.Listen.Port)
	}

	if listenErr != nil {
		logger.Fatal(listenErr)
	}

	server := &http.Server{
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	logger.Fatal(server.Serve(listener))
}
