package usecase

import (
	"context"

	"github.com/thanhnamdk2710/auth-service/internal/application/input"
	"github.com/thanhnamdk2710/auth-service/internal/application/output"
	"github.com/thanhnamdk2710/auth-service/internal/application/port"
	"github.com/thanhnamdk2710/auth-service/internal/domain/entity"
	"github.com/thanhnamdk2710/auth-service/internal/domain/exception"
	"github.com/thanhnamdk2710/auth-service/internal/domain/repository"
	"github.com/thanhnamdk2710/auth-service/internal/domain/vo"
	"github.com/thanhnamdk2710/auth-service/internal/pkg/correlationid"
)

type registerUseCase struct {
	userRepo      repository.UserRepository
	auditLogger   port.AuditLogger
	logger        port.Logger
	uuidGenerator port.UUIDGenerator
}

func NewRegisterUsecase(
	userRepo repository.UserRepository,
	auditLogger port.AuditLogger,
	logger port.Logger,
	uuidGenerator port.UUIDGenerator,
) port.RegisterUseCase {
	return &registerUseCase{
		userRepo:      userRepo,
		auditLogger:   auditLogger,
		logger:        logger,
		uuidGenerator: uuidGenerator,
	}
}

func (u *registerUseCase) Execute(ctx context.Context, input input.RegisterInput) (*output.RegisterOutput, error) {
	u.logger.InfoCtx(ctx, "Starting user registration",
		"username", input.Username,
		"email", input.Email,
	)

	username, err := vo.NewUsername(input.Username)
	if err != nil {
		return nil, err
	}

	exists, err := u.userRepo.ExistsByUsername(ctx, username.String())
	if err != nil {
		u.logger.ErrorCtx(ctx, "Failed to check username existence", "error", err)
		return nil, err
	}
	if exists {
		return nil, exception.ErrUsernameAlreadyExists
	}

	email, err := vo.NewEmail(input.Email)
	if err != nil {
		return nil, err
	}

	exists, err = u.userRepo.ExistsByEmail(ctx, email.String())
	if err != nil {
		u.logger.ErrorCtx(ctx, "Failed to check email existence", "error", err)
		return nil, err
	}
	if exists {
		return nil, exception.ErrEmailAlreadyExists
	}

	userID, err := vo.NewUserID(u.uuidGenerator.Generate())
	if err != nil {
		return nil, err
	}

	user := entity.NewUser(userID, *username, email)

	if err := u.userRepo.Create(ctx, user); err != nil {
		u.logger.ErrorCtx(ctx, "Failed to create user", "error", err)
		return nil, err
	}

	u.logAudit(ctx, user, input.IPAddress)

	u.logger.InfoCtx(ctx, "User registration completed",
		"user_id", user.ID.String(),
	)

	return &output.RegisterOutput{
		UserID:  user.ID.String(),
		Message: "User registered successfully",
	}, nil
}

func (u *registerUseCase) logAudit(ctx context.Context, user *entity.User, ipAddress string) {
	corrID := correlationid.FromContext(ctx)
	userIDStr := user.ID.String()

	auditLog, err := entity.NewAuditLog(
		entity.AuditActionUserRegistered,
		&userIDStr,
		map[string]interface{}{
			"username": user.Username.String(),
			"email":    user.Email.String(),
		},
		ipAddress,
		corrID,
	)
	if err != nil {
		u.logger.ErrorCtx(ctx, "Failed to create audit log", "error", err)
		return
	}

	u.auditLogger.Log(ctx, auditLog)
}
