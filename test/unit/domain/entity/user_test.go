package entity_test

import (
	"testing"

	"github.com/thanhnamdk2710/auth-service/internal/domain/entity"
	"github.com/thanhnamdk2710/auth-service/internal/domain/exception"
	"github.com/thanhnamdk2710/auth-service/internal/domain/vo"
)

func createValidUser(t *testing.T) *entity.User {
	t.Helper()

	userID, err := vo.NewUserID("0190a5b0-7e1c-7b3d-8f4e-9a1b2c3d4e5f")
	if err != nil {
		t.Fatalf("failed to create UserID: %v", err)
	}

	username, err := vo.NewUsername("testuser")
	if err != nil {
		t.Fatalf("failed to create Username: %v", err)
	}

	email, err := vo.NewEmail("test@example.com")
	if err != nil {
		t.Fatalf("failed to create Email: %v", err)
	}

	return entity.NewUser(userID, *username, email)
}

func TestNewUser(t *testing.T) {
	user := createValidUser(t)

	if user == nil {
		t.Fatal("NewUser returned nil")
	}

	// New users should be active by default
	if !user.IsActive {
		t.Error("NewUser should create active user by default")
	}

	// New users should not have verified email
	if user.IsEmailVerified {
		t.Error("NewUser should create user with unverified email")
	}

	// Verify value objects
	if user.ID.String() != "0190a5b0-7e1c-7b3d-8f4e-9a1b2c3d4e5f" {
		t.Errorf("User.ID = %q, want %q", user.ID.String(), "0190a5b0-7e1c-7b3d-8f4e-9a1b2c3d4e5f")
	}

	if user.Username.String() != "testuser" {
		t.Errorf("User.Username = %q, want %q", user.Username.String(), "testuser")
	}

	if user.Email.String() != "test@example.com" {
		t.Errorf("User.Email = %q, want %q", user.Email.String(), "test@example.com")
	}
}

func TestUser_Activate(t *testing.T) {
	t.Run("activate inactive user", func(t *testing.T) {
		user := createValidUser(t)
		// Deactivate first
		_ = user.Deactivate()

		err := user.Activate()
		if err != nil {
			t.Errorf("Activate() unexpected error: %v", err)
		}

		if !user.IsActive {
			t.Error("Activate() should set IsActive to true")
		}
	})

	t.Run("activate already active user", func(t *testing.T) {
		user := createValidUser(t)
		// User is active by default

		err := user.Activate()
		if err != exception.ErrUserAlreadyActive {
			t.Errorf("Activate() expected error %v, got %v", exception.ErrUserAlreadyActive, err)
		}
	})
}

func TestUser_Deactivate(t *testing.T) {
	t.Run("deactivate active user", func(t *testing.T) {
		user := createValidUser(t)

		err := user.Deactivate()
		if err != nil {
			t.Errorf("Deactivate() unexpected error: %v", err)
		}

		if user.IsActive {
			t.Error("Deactivate() should set IsActive to false")
		}
	})

	t.Run("deactivate already inactive user", func(t *testing.T) {
		user := createValidUser(t)
		_ = user.Deactivate()

		err := user.Deactivate()
		if err != exception.ErrUserAlreadyInactive {
			t.Errorf("Deactivate() expected error %v, got %v", exception.ErrUserAlreadyInactive, err)
		}
	})
}

func TestUser_VerifyEmail(t *testing.T) {
	t.Run("verify email for active user", func(t *testing.T) {
		user := createValidUser(t)

		err := user.VerifyEmail()
		if err != nil {
			t.Errorf("VerifyEmail() unexpected error: %v", err)
		}

		if !user.IsEmailVerified {
			t.Error("VerifyEmail() should set IsEmailVerified to true")
		}
	})

	t.Run("verify email for inactive user", func(t *testing.T) {
		user := createValidUser(t)
		_ = user.Deactivate()

		err := user.VerifyEmail()
		if err != exception.ErrUserAlreadyInactive {
			t.Errorf("VerifyEmail() expected error %v, got %v", exception.ErrUserAlreadyInactive, err)
		}

		if user.IsEmailVerified {
			t.Error("VerifyEmail() should not set IsEmailVerified for inactive user")
		}
	})

	t.Run("verify already verified email", func(t *testing.T) {
		user := createValidUser(t)
		_ = user.VerifyEmail()

		err := user.VerifyEmail()
		if err != exception.ErrEmailAlreadyVerified {
			t.Errorf("VerifyEmail() expected error %v, got %v", exception.ErrEmailAlreadyVerified, err)
		}
	})
}

func TestUser_StateTransitions(t *testing.T) {
	t.Run("activate -> deactivate -> activate", func(t *testing.T) {
		user := createValidUser(t)

		// Initial state: active
		if !user.IsActive {
			t.Fatal("initial state should be active")
		}

		// Deactivate
		if err := user.Deactivate(); err != nil {
			t.Fatalf("Deactivate() failed: %v", err)
		}
		if user.IsActive {
			t.Fatal("user should be inactive after Deactivate")
		}

		// Activate again
		if err := user.Activate(); err != nil {
			t.Fatalf("Activate() failed: %v", err)
		}
		if !user.IsActive {
			t.Fatal("user should be active after Activate")
		}
	})
}
