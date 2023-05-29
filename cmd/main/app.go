package main

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net"
	"net/http"
	"time"
	"wifi-scaner-credentials/internal/config"
	"wifi-scaner-credentials/internal/credentials"
	"wifi-scaner-credentials/internal/devices/api"
	api2 "wifi-scaner-credentials/internal/passwords/api"
	"wifi-scaner-credentials/pkg/client/yacloud/IoT"
	"wifi-scaner-credentials/pkg/client/yacloud/autorisation"
	"wifi-scaner-credentials/pkg/logging"
)

func main() {
	logger := logging.GetLogger()

	logger.Info("create router")
	router := httprouter.New()

	cfg := config.GetConfig()

	key, err := autorisation.GetKeyFromFile(cfg.API.AuthorizedKeyPath)
	if err != nil {
		logger.Fatalf("unable to read authorized key file from %s", cfg.API.AuthorizedKeyPath)
	}

	apiClient := IoT.NewClient(
		key,
		cfg.API.RegistryID,
	)

	deviceRepo := api.NewRepository(apiClient, logger.GetLoggerWithField("api", "deviceRepo"))
	passRepo := api2.NewRepository(apiClient)

	logger.Info("register credentials handlers")
	handler := credentials.NewHandler(deviceRepo, passRepo, logger)
	handler.Register(router)

	start(router, cfg)
}

func start(router *httprouter.Router, cfg *config.Config) {
	logger := logging.GetLogger()
	logger.Info("start application")

	var listener net.Listener
	var err error

	if cfg.Listen.Type == "sock" {
		panic("not implemented")
	} else {
		listener, err = net.Listen("tcp", fmt.Sprintf("%s:%s", cfg.Listen.Host, cfg.Listen.Port))
		if err != nil {
			panic(err)
		}
	}

	server := http.Server{
		Handler:           router,
		WriteTimeout:      15 * time.Second,
		ReadHeaderTimeout: 15 * time.Second,
	}

	logger.Infof("listen tcp %s on port %s", cfg.Listen.Host, cfg.Listen.Port)

	logger.Fatal(server.ServeTLS(listener, "ssl/server.crt", "ssl/server.key"))
}
