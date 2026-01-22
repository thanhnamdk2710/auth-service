package usecase

import (
	"github.com/thanhnamdk2710/auth-service/internal/application/input"
	"github.com/thanhnamdk2710/auth-service/internal/application/output"
	"github.com/thanhnamdk2710/auth-service/internal/domain/repository"
)

type RegisterUseCase interface {
	Execute(input input.RegisterInput) (output.RegisterOutput, error)
}

type registerUseCase struct {
	userRepo repository.UserRepository
}

func NewRegisterUsecase(userRepo repository.UserRepository) *registerUseCase {
	return &registerUseCase{
		userRepo: userRepo,
	}
}

func (u registerUseCase) Execute(input input.RegisterInput) (output.RegisterOutput, error) {
	return output.RegisterOutput{
		Message: "Register API",
	}, nil
}
