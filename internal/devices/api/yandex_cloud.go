package api

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"wifi-scaner-credentials/internal/devices"
	"wifi-scaner-credentials/pkg/client/yacloud/IoT"
	"wifi-scaner-credentials/pkg/logging"
)

const createURL = "https://iot-devices.api.cloud.yandex.net/iot-devices/v1/devices"

type repository struct {
	client *IoT.Client
	logger *logging.Logger
}

func (r *repository) Create(ctx context.Context, mac string) (string, error) {
	body, err := json.Marshal(createRequest{
		r.client.RegisterID,
		mac,
	})
	if err != nil {
		r.logger.Error("unable to serialize json")
		return "", err
	}

	bodyReader := bytes.NewReader(body)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, createURL, bodyReader)
	if err != nil {
		r.logger.Error("unable to create request")
		return "", err
	}

	res, err := r.client.MakeRequest(req)
	if err != nil {
		r.logger.Error("unable to make request")
		return "", err
	}

	r.logger.Infof("HTTP code: %d", res.StatusCode)

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("status code: %d", res.StatusCode)
	}

	resBody, err := IoT.ParseResponse(res)
	if err != nil {
		return "", err
	}

	var result response
	err = json.Unmarshal(resBody, &result)
	if err != nil {
		return "", err
	}

	return result.Device.ID, nil
}

func (r *repository) FindAll(ctx context.Context) ([]devices.Device, error) {
	//TODO implement me
	panic("implement me")
}

func (r *repository) FindByMAC(ctx context.Context, mac string) (*devices.Device, error) {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		fmt.Sprintf("%s:getByName?registryId=%s&deviceName=%s", createURL, r.client.RegisterID, mac),
		nil,
	)
	if err != nil {
		return nil, err
	}

	res, err := r.client.MakeRequest(req)
	if err != nil {
		return nil, err
	}

	switch res.StatusCode {
	case http.StatusNotFound:
		return nil, nil
	}

	resBody, err := IoT.ParseResponse(res)
	if err != nil {
		return nil, err
	}

	var result devices.Device
	err = json.Unmarshal(resBody, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (r *repository) Delete(ctx context.Context, deviceID string) error {
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodDelete,
		fmt.Sprintf("%s/%s", createURL, deviceID),
		nil,
	)
	if err != nil {
		return err
	}

	res, err := r.client.MakeRequest(req)
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		return fmt.Errorf("status code: %d", res.StatusCode)
	}

	return nil
}

func NewRepository(client *IoT.Client, logger *logging.Logger) devices.Repository {
	return &repository{client: client, logger: logger}
}

type createRequest struct {
	RegistryID string `json:"registryId"`
	Name       string `json:"name"`
}

type response struct {
	Device devices.Device `json:"response"`
}
