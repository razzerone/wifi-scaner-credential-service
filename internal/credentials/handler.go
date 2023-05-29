package credentials

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"regexp"
	"wifi-scaner-credentials/internal/devices"
	"wifi-scaner-credentials/internal/handlers"
	"wifi-scaner-credentials/internal/passwords"
	"wifi-scaner-credentials/pkg/logging"
)

const (
	credentialURL = "/credentials/:mac"
)

var _ handlers.Handler = &handler{}

type handler struct {
	DeviceRepo   devices.Repository
	PasswordRepo passwords.Repository
	logger       *logging.Logger
}

func NewHandler(repository devices.Repository, repository2 passwords.Repository, logger *logging.Logger) handlers.Handler {
	return &handler{repository, repository2, logger}
}

func (h *handler) Register(router *httprouter.Router) {
	router.GET(credentialURL, h.getCredentials)
}

func (h *handler) getCredentials(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
	mac := params.ByName("mac")

	h.logger.Infof("request from: %s", mac)

	matched, err := regexp.MatchString(`^([0-9A-F]{2}-){5}([0-9A-F]{2})$`, mac)
	if err != nil || !matched {
		h.logger.Error("invalid mac address")
		http.Error(w, "invalid mac address", http.StatusBadRequest)
		return
	}

	device, err := h.DeviceRepo.FindByMAC(context.Background(), mac)
	if err != nil {
		h.logger.Error("api error")
		http.Error(w, fmt.Sprintf("api error: %e", err), http.StatusInternalServerError)
		return
	}

	if device == nil {
		h.logger.Info("device with mac %s not found\n", mac)

		h.logger.Info("creating device")

		deviceID, err := h.DeviceRepo.Create(context.Background(), mac)
		if err != nil {
			h.logger.Error("api error")
			http.Error(w, fmt.Sprintf("api error: %e", err), http.StatusInternalServerError)
			return
		}

		device = &devices.Device{ID: deviceID, MAC: mac}

		h.logger.Infof("device created, id: %s", deviceID)
	}

	// TODO: generate pass
	password := "QQQQQQ1111notimplemented"

	pass := &passwords.Password{
		DeviceID: device.ID,
		Password: password,
	}

	err = h.PasswordRepo.Create(context.Background(), pass)
	if err != nil {
		h.logger.Error("api error")
		http.Error(w, fmt.Sprintf("api error: %e", err), http.StatusInternalServerError)
		return
	}

	h.logger.Info("password created")

	cred := &Credentials{
		DeviceID: device.ID,
		Password: password,
	}

	marshal, err := json.Marshal(cred)
	if err != nil {
		h.logger.Error("json error")
		http.Error(w, fmt.Sprintf("internal json error: %e", err), http.StatusInternalServerError)
		return
	}

	_, err = w.Write(marshal)
	if err != nil {
		h.logger.Error("write error")
		http.Error(w, fmt.Sprintf("internal error: %e", err), http.StatusInternalServerError)
		return
	}
}
