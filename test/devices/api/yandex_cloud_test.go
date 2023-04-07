package api

import (
	"context"
	"fmt"
	"testing"
	"wifi-scaner-credentials/internal/devices"
	"wifi-scaner-credentials/internal/devices/api"
	"wifi-scaner-credentials/pkg/client/yacloud/IoT"
	"wifi-scaner-credentials/pkg/client/yacloud/autorisation"
)

const registryID = "arehtnb60pd5dgjvfrel"
const mac = "autotestmac3"

var deviceID string

var client = IoT.NewClient(
	autorisation.GetKeyFromFile("/home/razzerone/GolandProjects/wifi-scaner-credentials/authorized_key.json"),
	registryID,
)

func TestAPICreateDeviceOK(t *testing.T) {
	repo := api.NewRepository(client)

	create, err := repo.Create(context.Background(), mac)
	if err != nil {
		t.Error(err)
		fmt.Println()
		return
	}

	if create == "" {
		t.Error("empty device id")
	}

	fmt.Printf("deviceID: %s\n", create)

	deviceID = create
}

func TestAPIGetDeviceByMACOK(t *testing.T) {
	repo := api.NewRepository(client)

	device := &devices.Device{MAC: mac}
	byMAC, err := repo.FindByMAC(context.Background(), device)
	if err != nil {
		t.Error(err)
		fmt.Println()
		return
	}

	fmt.Println(byMAC)

	if byMAC.ID == "" {
		t.Error("empty device id")
	}

}

func TestAPIDeleteDeviceOK(t *testing.T) {
	deviceID = "areufp3s6kjmbssomt0c"

	repo := api.NewRepository(client)

	err := repo.Delete(context.Background(), deviceID)
	if err != nil {
		t.Error(err)
		fmt.Println()
		return
	}
}
