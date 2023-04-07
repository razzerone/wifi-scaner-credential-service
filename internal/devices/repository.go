package devices

import "context"

type Repository interface {
	Create(ctx context.Context, mac string) (string, error)
	FindAll(ctx context.Context) ([]Device, error)
	FindByMAC(ctx context.Context, mac string) (*Device, error)
	Delete(ctx context.Context, deviceID string) error
}
