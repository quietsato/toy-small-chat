package domain_test

import (
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/quietsato/toy-small-chat/api/internal/domain"
	"github.com/stretchr/testify/require"
)

// AccountID tests
func TestParseAccountID(t *testing.T) {
	t.Parallel()

	t.Run("valid UUID", func(t *testing.T) {
		t.Parallel()
		validUUID := "550e8400-e29b-41d4-a716-446655440000"
		accountID, err := domain.ParseAccountID(validUUID)
		require.NoError(t, err)
		require.Equal(t, validUUID, accountID.String())
	})

	t.Run("invalid UUID", func(t *testing.T) {
		t.Parallel()

		invalidUUID := "not-a-uuid"
		_, err := domain.ParseAccountID(invalidUUID)
		require.Error(t, err)
	})

	t.Run("empty string", func(t *testing.T) {
		t.Parallel()

		_, err := domain.ParseAccountID("")
		require.Error(t, err)
	})
}

func TestAccountIDFromUuid(t *testing.T) {
	t.Parallel()

	t.Run("create from UUID", func(t *testing.T) {
		t.Parallel()
		u := uuid.New()
		accountID := domain.AccountIDFromUuid(u)
		require.Equal(t, u.String(), accountID.String())
	})
}

func TestAccountID_String(t *testing.T) {
	t.Parallel()

	t.Run("convert to string", func(t *testing.T) {
		t.Parallel()
		validUUID := "550e8400-e29b-41d4-a716-446655440000"
		accountID, _ := domain.ParseAccountID(validUUID)
		require.Equal(t, validUUID, accountID.String())
	})
}

// UserName tests
func TestNewUserName(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		input     string
		wantError bool
	}{
		{"valid lowercase", "user123", false},
		{"valid uppercase", "USER123", false},
		{"valid mixed case", "User123", false},
		{"valid single char", "a", false},
		{"valid max length 32", "abcdefghijklmnopqrstuvwxyz123456", false},
		{"empty string", "", true},
		{"too long 33 chars", "abcdefghijklmnopqrstuvwxyz1234567", true},
		{"contains space", "user name", true},
		{"contains special chars", "user@name", true},
		{"contains hyphen", "user-name", true},
		{"contains underscore", "user_name", true},
		{"contains unicode", "user名前", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			userName, err := domain.NewUserName(tt.input)
			if tt.wantError {
				require.Error(t, err)
				require.ErrorIs(t, err, domain.ErrInvalidUserName)
			} else {
				require.NoError(t, err)
				require.Equal(t, tt.input, userName.String())
			}
		})
	}
}

func TestUserName_String(t *testing.T) {
	t.Parallel()

	t.Run("convert to string", func(t *testing.T) {
		t.Parallel()

		input := "testuser"
		userName, _ := domain.NewUserName(input)
		require.Equal(t, input, userName.String())
	})
}

// Password tests
func TestNewPassword(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		input     string
		wantError bool
	}{
		{"valid password min length", "abcd1234", false},
		{"valid password with special chars", "Abcd123!@#$%", false},
		{"valid password 72 chars (bcrypt max)", "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*", false},
		{"too short 7 chars", "abcd123", true},
		{"too long 73 chars", strings.Repeat("a", 73), true},
		{"empty string", "", true},
		{"contains space", "abcd 1234", true},
		{"contains invalid special char", "abcd1234-test", true},
		{"contains unicode", "password日本語", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			rawPassword, err := domain.NewRawPassword([]byte(tt.input))
			if tt.wantError {
				require.Error(t, err)
				require.ErrorIs(t, err, domain.ErrInvalidPassword)
			} else {
				require.NoError(t, err)
				// Create hashed password from raw password
				hashedPassword, err := domain.NewHashedPassword(rawPassword)
				require.NoError(t, err)
				// Verify that hash was generated
				require.NotEmpty(t, hashedPassword.String())
				// Verify that hash is different from original password
				require.NotEqual(t, tt.input, hashedPassword.String())
			}
		})
	}
}

func TestNewPasswordFromHash(t *testing.T) {
	t.Parallel()

	t.Run("create from existing hash", func(t *testing.T) {
		t.Parallel()
		hash := []byte("$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy")
		password := domain.NewHashedPasswordFromHash(hash)
		require.Equal(t, string(hash), password.String())
	})
}

func TestPassword_HashString(t *testing.T) {
	t.Parallel()

	t.Run("get hash string", func(t *testing.T) {
		t.Parallel()
		rawPasswordStr := "testpass123"
		rawPassword, _ := domain.NewRawPassword([]byte(rawPasswordStr))
		hashedPassword, _ := domain.NewHashedPassword(rawPassword)
		hash := hashedPassword.String()
		require.NotEmpty(t, hash)
		require.NotEqual(t, rawPasswordStr, hash)
	})
}

func TestPassword_Match(t *testing.T) {
	t.Parallel()

	t.Run("matching password", func(t *testing.T) {
		t.Parallel()
		rawPasswordBytes := []byte("testpass123")
		rawPassword, _ := domain.NewRawPassword(rawPasswordBytes)
		hashedPassword, _ := domain.NewHashedPassword(rawPassword)

		matched, err := hashedPassword.Match(rawPassword)
		require.NoError(t, err)
		require.True(t, matched)
	})

	t.Run("non-matching password", func(t *testing.T) {
		t.Parallel()

		rawPasswordBytes := []byte("testpass123")
		wrongPasswordBytes := []byte("wrongpass456")
		rawPassword, _ := domain.NewRawPassword(rawPasswordBytes)
		wrongPassword, _ := domain.NewRawPassword(wrongPasswordBytes)
		hashedPassword, _ := domain.NewHashedPassword(rawPassword)

		matched, err := hashedPassword.Match(wrongPassword)
		require.NoError(t, err)
		require.False(t, matched)
	})

	t.Run("empty password", func(t *testing.T) {
		t.Parallel()

		rawPasswordBytes := []byte("testpass123")
		rawPassword, _ := domain.NewRawPassword(rawPasswordBytes)
		hashedPassword, _ := domain.NewHashedPassword(rawPassword)

		emptyPassword, err := domain.NewRawPassword([]byte(""))
		require.Error(t, err)
		require.ErrorIs(t, err, domain.ErrInvalidPassword)

		// Test with zero-value RawPassword
		matched, _ := hashedPassword.Match(emptyPassword)
		require.False(t, matched)
	})
}
