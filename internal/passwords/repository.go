package passwords

import (
	"context"
)

type Repository interface {
	Create(ctx context.Context, password *Password) error
	List(ctx context.Context, deviceID string) ([]Password, error)
	Delete(ctx context.Context, password *Password) error
}
