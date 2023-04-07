package api

import (
	"context"
	"fmt"
	"testing"
	deviceapi "wifi-scaner-credentials/internal/devices/api"
	"wifi-scaner-credentials/internal/passwords"
	"wifi-scaner-credentials/internal/passwords/api"
	"wifi-scaner-credentials/pkg/client/yacloud/IoT"
	"wifi-scaner-credentials/pkg/client/yacloud/autorisation"
)

const registryID = "arehtnb60pd5dgjvfrel"
const mac = "autotestmac"

var client = IoT.NewClient(
	autorisation.GetKeyFromFile("/home/razzerone/GolandProjects/wifi-scaner-credentials/authorized_key.json"),
	registryID,
)

func TestAPICreatePasswordOK(t *testing.T) {
	repo := api.NewRepository(client)
	drepo := deviceapi.NewRepository(client)

	create, err := drepo.Create(context.Background(), mac)
	if err != nil {
		t.Fatal(err)
		return
	}

	password := &passwords.Password{DeviceID: create, Password: "Qwerty1234567890"}
	err = repo.Create(context.Background(), password)
	if err != nil {
		t.Error(err)
		return
	}
}

func TestAPIListPasswordOK(t *testing.T) {
	repo := api.NewRepository(client)

	list, err := repo.List(context.Background(), "arebpeusi8tuu0np8bhb")
	if err != nil {
		t.Error(err)
		return
	}

	fmt.Println(list)
}

func TestAPIDeletePasswordOK(t *testing.T) {
	repo := api.NewRepository(client)

	pass := &passwords.Password{DeviceID: "arebpeusi8tuu0np8bhb", ID: "arefv7j8tqo1b2n5221m"}
	err := repo.Delete(context.Background(), pass)
	if err != nil {
		t.Error(err)
		return
	}
}
