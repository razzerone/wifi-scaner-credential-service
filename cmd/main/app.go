package main

import (
	"github.com/julienschmidt/httprouter"
	"log"
	"net"
	"net/http"
	"time"
	"wifi-scaner-credentials/internal/credentials"
	"wifi-scaner-credentials/internal/devices/api"
	api2 "wifi-scaner-credentials/internal/passwords/api"
	"wifi-scaner-credentials/pkg/client/yacloud/IoT"
	"wifi-scaner-credentials/pkg/client/yacloud/autorisation"
)

func main() {
	log.Println("create router")
	router := httprouter.New()

	apiClient := IoT.NewClient(
		autorisation.GetKeyFromFile("authorized_key.json"),
		"arehtnb60pd5dgjvfrel",
	)

	deviceRepo := api.NewRepository(apiClient)
	passRepo := api2.NewRepository(apiClient)

	log.Println("register credentials handlers")
	handler := credentials.NewHandler(deviceRepo, passRepo)
	handler.Register(router)

	start(router)
}

func start(router *httprouter.Router) {
	log.Println("start application")

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	server := http.Server{
		Handler:           router,
		WriteTimeout:      15 * time.Second,
		ReadHeaderTimeout: 15 * time.Second,
	}

	log.Fatal(server.Serve(listener))
}
