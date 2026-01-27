package entity

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type AuditAction string

const (
	AuditActionUserRegistered  AuditAction = "USER_REGISTERED"
	AuditActionUserLogin       AuditAction = "USER_LOGIN"
	AuditActionUserLoginFailed AuditAction = "USER_LOGIN_FAILED"
	AuditActionPasswordChanged AuditAction = "PASSWORD_CHANGED"
	AuditActionPasswordReset   AuditAction = "PASSWORD_RESET"
	AuditActionEmailVerified   AuditAction = "EMAIL_VERIFIED"
)

type AuditLog struct {
	ID            string
	Timestamp     time.Time
	UserID        *string
	Action        AuditAction
	Details       json.RawMessage
	IPAddress     string
	CorrelationID string
}

func NewAuditLog(action AuditAction, userID *string, details map[string]interface{}, ipAddress, correlationID string) (*AuditLog, error) {
	detailsJSON, err := json.Marshal(details)
	if err != nil {
		return nil, err
	}

	return &AuditLog{
		ID:            uuid.New().String(),
		Timestamp:     time.Now().UTC(),
		UserID:        userID,
		Action:        action,
		Details:       detailsJSON,
		IPAddress:     ipAddress,
		CorrelationID: correlationID,
	}, nil
}
