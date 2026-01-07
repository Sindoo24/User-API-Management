package service

import (
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestValidatePasswordStrength(t *testing.T) {
	service := &AuthService{}

	tests := []struct {
		name        string
		password    string
		expectError error
	}{
		{
			name:        "Valid strong password",
			password:    "SecurePass123!",
			expectError: nil,
		},
		{
			name:        "Password too short",
			password:    "Short1!",
			expectError: ErrPasswordTooShort,
		},
		{
			name:        "No uppercase letter",
			password:    "securepass123!",
			expectError: ErrPasswordNoUppercase,
		},
		{
			name:        "No lowercase letter",
			password:    "SECUREPASS123!",
			expectError: ErrPasswordNoLowercase,
		},
		{
			name:        "No digit",
			password:    "SecurePass!",
			expectError: ErrPasswordNoDigit,
		},
		{
			name:        "No special character",
			password:    "SecurePass123",
			expectError: ErrPasswordNoSpecial,
		},
		{
			name:        "All special characters work",
			password:    "Password1@",
			expectError: nil,
		},
		{
			name:        "Different special character",
			password:    "Password1#",
			expectError: nil,
		},
		{
			name:        "Minimum length with all requirements",
			password:    "Passw0rd!",
			expectError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := service.ValidatePasswordStrength(tt.password)
			if err != tt.expectError {
				t.Errorf("ValidatePasswordStrength(%q) = %v; want %v", tt.password, err, tt.expectError)
			}
		})
	}
}

func TestHashPassword(t *testing.T) {
	service := &AuthService{}

	t.Run("Successfully hashes password", func(t *testing.T) {
		password := "TestPassword123!"
		hash, err := service.HashPassword(password)

		if err != nil {
			t.Fatalf("HashPassword failed: %v", err)
		}

		if hash == "" {
			t.Error("HashPassword returned empty hash")
		}

		if hash == password {
			t.Error("HashPassword returned plain password instead of hash")
		}

		// Verify it's a valid bcrypt hash (starts with $2a$ or $2b$)
		if len(hash) < 60 {
			t.Errorf("Hash too short: %d characters", len(hash))
		}
	})

	t.Run("Different passwords produce different hashes", func(t *testing.T) {
		password1 := "Password123!"
		password2 := "DifferentPass456!"

		hash1, err1 := service.HashPassword(password1)
		hash2, err2 := service.HashPassword(password2)

		if err1 != nil || err2 != nil {
			t.Fatalf("HashPassword failed: %v, %v", err1, err2)
		}

		if hash1 == hash2 {
			t.Error("Different passwords produced same hash")
		}
	})

	t.Run("Same password produces different hashes (salt)", func(t *testing.T) {
		password := "SamePassword123!"

		hash1, _ := service.HashPassword(password)
		hash2, _ := service.HashPassword(password)

		// Bcrypt includes random salt, so same password should produce different hashes
		if hash1 == hash2 {
			t.Error("Same password produced identical hash (salt not working)")
		}
	})
}

func TestComparePassword(t *testing.T) {
	service := &AuthService{}

	t.Run("Correct password matches hash", func(t *testing.T) {
		password := "CorrectPassword123!"
		hash, _ := service.HashPassword(password)

		err := service.ComparePassword(hash, password)
		if err != nil {
			t.Errorf("ComparePassword failed for correct password: %v", err)
		}
	})

	t.Run("Incorrect password does not match hash", func(t *testing.T) {
		password := "CorrectPassword123!"
		wrongPassword := "WrongPassword456!"
		hash, _ := service.HashPassword(password)

		err := service.ComparePassword(hash, wrongPassword)
		if err == nil {
			t.Error("ComparePassword succeeded for incorrect password")
		}

		if err != bcrypt.ErrMismatchedHashAndPassword {
			t.Errorf("Expected ErrMismatchedHashAndPassword, got: %v", err)
		}
	})

	t.Run("Empty password does not match", func(t *testing.T) {
		password := "Password123!"
		hash, _ := service.HashPassword(password)

		err := service.ComparePassword(hash, "")
		if err == nil {
			t.Error("ComparePassword succeeded for empty password")
		}
	})
}

func TestPasswordHashCost(t *testing.T) {
	service := &AuthService{}

	password := "TestPassword123!"
	hash, err := service.HashPassword(password)

	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	// Extract and verify bcrypt cost
	cost, err := bcrypt.Cost([]byte(hash))
	if err != nil {
		t.Fatalf("Failed to get bcrypt cost: %v", err)
	}

	expectedCost := 12
	if cost != expectedCost {
		t.Errorf("Bcrypt cost = %d; want %d", cost, expectedCost)
	}
}
