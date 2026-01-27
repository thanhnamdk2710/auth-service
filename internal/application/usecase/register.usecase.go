package usecase

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"go.uber.org/zap"

	"github.com/thanhnamdk2710/auth-service/internal/application/input"
	"github.com/thanhnamdk2710/auth-service/internal/application/output"
	"github.com/thanhnamdk2710/auth-service/internal/application/service"
	"github.com/thanhnamdk2710/auth-service/internal/domain/entity"
	"github.com/thanhnamdk2710/auth-service/internal/domain/repository"
	"github.com/thanhnamdk2710/auth-service/internal/domain/vo"
	"github.com/thanhnamdk2710/auth-service/internal/pkg/correlationid"
	"github.com/thanhnamdk2710/auth-service/internal/pkg/logger"
)

type RegisterUseCase interface {
	Execute(ctx context.Context, input input.RegisterInput) (*output.RegisterOutput, error)
}

type registerUseCase struct {
	userRepo     repository.UserRepository
	auditService service.AuditService
	logger       *logger.Logger
}

func NewRegisterUsecase(userRepo repository.UserRepository, auditService service.AuditService, logger *logger.Logger) RegisterUseCase {
	return &registerUseCase{
		userRepo:     userRepo,
		auditService: auditService,
		logger:       logger,
	}
}

func (u *registerUseCase) Execute(ctx context.Context, input input.RegisterInput) (*output.RegisterOutput, error) {
	u.logger.InfoCtx(ctx, "Starting user registration",
		zap.String("username", input.Username),
		zap.String("email", input.Email),
	)

	username, err := vo.NewUsername(input.Username)
	if err != nil {
		return nil, err
	}

	exists, err := u.userRepo.ExistsByUsername(ctx, username.String())
	if err != nil {
		u.logger.ErrorCtx(ctx, "Failed to check username existence", zap.Error(err))
		return nil, err
	}
	if exists {
		return nil, errors.New("username already exists")
	}

	email, err := vo.NewEmail(input.Email)
	if err != nil {
		return nil, err
	}

	exists, err = u.userRepo.ExistsByEmail(ctx, email.String())
	if err != nil {
		u.logger.ErrorCtx(ctx, "Failed to check email existence", zap.Error(err))
		return nil, err
	}
	if exists {
		return nil, errors.New("email already exists")
	}

	userID, err := vo.NewUserID(uuid.New().String())
	if err != nil {
		return nil, err
	}

	user := entity.NewUser(userID, *username, email)

	if err := u.userRepo.Create(ctx, user); err != nil {
		u.logger.ErrorCtx(ctx, "Failed to create user", zap.Error(err))
		return nil, err
	}

	corrID := correlationid.FromContext(ctx)
	userIDStr := user.ID.String()
	auditLog, _ := entity.NewAuditLog(
		entity.AuditActionUserRegistered,
		&userIDStr,
		map[string]interface{}{
			"username": user.Username.String(),
			"email":    user.Email.String(),
		},
		input.IPAddress,
		corrID,
	)
	u.auditService.Log(ctx, auditLog)

	u.logger.InfoCtx(ctx, "User registration completed",
		zap.String("user_id", user.ID.String()),
	)

	return &output.RegisterOutput{
		UserID:  user.ID.String(),
		Message: "User registered successfully",
	}, nil
}
