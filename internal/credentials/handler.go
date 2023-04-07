package credentials

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"regexp"
	"wifi-scaner-credentials/internal/devices"
	"wifi-scaner-credentials/internal/handlers"
	"wifi-scaner-credentials/internal/passwords"
)

const (
	credentialURL = "/credentials/:mac"
)

var _ handlers.Handler = &handler{}

type handler struct {
	DeviceRepo   devices.Repository
	PasswordRepo passwords.Repository
}

func NewHandler(repository devices.Repository, repository2 passwords.Repository) handlers.Handler {
	return &handler{repository, repository2}
}

func (h *handler) Register(router *httprouter.Router) {
	router.GET(credentialURL, h.getCredentials)
}

func (h *handler) getCredentials(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	mac := params.ByName("mac")

	matched, err := regexp.MatchString(`^([0-9A-F]{2}-){5}([0-9A-F]{2})$`, mac)
	if err != nil || !matched {
		log.Printf("")
		http.Error(w, "invalid mac address", http.StatusBadRequest)
		return
	}

	device, err := h.DeviceRepo.FindByMAC(context.Background(), mac)
	if err != nil {
		log.Printf("")
		http.Error(w, fmt.Sprintf("api error: %e", err), http.StatusInternalServerError)
		return
	}

	if device == nil {
		log.Printf("device with mac %s not found\n", mac)

		log.Println("creating device")

		deviceID, err := h.DeviceRepo.Create(context.Background(), mac)
		if err != nil {
			log.Printf("")
			http.Error(w, fmt.Sprintf("api error: %e", err), http.StatusInternalServerError)
			return
		}

		device = &devices.Device{ID: deviceID, MAC: mac}

		log.Printf("device created, id: %s", deviceID)
	}

	// generate pass
	password := "QQQQQQ1111notimplemented"

	pass := &passwords.Password{
		DeviceID: device.ID,
		Password: password,
	}

	err = h.PasswordRepo.Create(context.Background(), pass)
	if err != nil {
		log.Printf("")
		http.Error(w, fmt.Sprintf("api error: %e", err), http.StatusInternalServerError)
		return
	}

	log.Println("password created")

	cred := &Credentials{
		DeviceID: device.ID,
		Password: password,
	}

	marshal, err := json.Marshal(cred)
	if err != nil {
		log.Printf("")
		http.Error(w, fmt.Sprintf("internal json error: %e", err), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(marshal)
	if err != nil {
		log.Printf("")
		http.Error(w, fmt.Sprintf("internal error: %e", err), http.StatusInternalServerError)
		return
	}
}
