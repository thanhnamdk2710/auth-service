package port

import (
	"context"

	"github.com/thanhnamdk2710/auth-service/internal/application/input"
	"github.com/thanhnamdk2710/auth-service/internal/application/output"
)

type RegisterUseCase interface {
	Execute(ctx context.Context, input input.RegisterInput) (*output.RegisterOutput, error)
}
