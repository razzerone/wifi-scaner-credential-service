package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"wifi-scaner-credentials/internal/passwords"
	"wifi-scaner-credentials/pkg/client/yacloud/IoT"
)

const (
	createURL = "https://iot-devices.api.cloud.yandex.net/iot-devices/v1/devices"
)

type repository struct {
	client *IoT.Client
}

//TODO целостность данных

func (r repository) Create(ctx context.Context, password *passwords.Password) error {
	body, err := json.Marshal(password)
	if err != nil {
		return err
	}

	bodyReader := bytes.NewReader(body)
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		fmt.Sprintf("%s/%s/passwords", createURL, password.DeviceID),
		bodyReader,
	)
	if err != nil {
		return err
	}

	_, err = r.client.MakeRequest(req)
	if err != nil {
		return err
	}

	return nil
}

func (r repository) List(ctx context.Context, deviceID string) ([]passwords.Password, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s/%s/passwords", createURL, deviceID),
		nil,
	)
	if err != nil {
		return nil, err
	}

	res, err := r.client.MakeRequest(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code: %d", res.StatusCode)
	}

	respBody, err := IoT.ParseResponse(res)
	if err != nil {
		return nil, err
	}

	var result response
	err = json.Unmarshal(respBody, &result)
	if err != nil {
		return nil, err
	}

	return result.Passwords, nil
}

func (r repository) Delete(ctx context.Context, password *passwords.Password) error {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodDelete,
		fmt.Sprintf("%s/%s/passwords/%s", createURL, password.DeviceID, password.ID),
		nil,
	)
	if err != nil {
		return err
	}

	res, err := r.client.MakeRequest(req)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("status code: %d", res.StatusCode)
	}

	return nil
}

func NewRepository(client *IoT.Client) passwords.Repository {
	return &repository{client: client}
}

type response struct {
	Passwords []passwords.Password `json:"passwords"`
}
