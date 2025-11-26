package domain

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AccountID struct {
	uuid uuid.UUID
}

func ParseAccountID(s string) (AccountID, error) {
	u, err := uuid.Parse(s)
	if err != nil {
		return AccountID{}, fmt.Errorf("failed to parse uuid: %w", err)
	}

	return AccountID{uuid: u}, nil
}

func AccountIDFromUuid(u uuid.UUID) AccountID {
	return AccountID{uuid: u}
}

func (a AccountID) String() string {
	return a.uuid.String()
}

type UserName struct {
	userName string
}

var userNameRegExp = regexp.MustCompile("^[a-zA-Z0-9]{1,32}$")

var (
	ErrInvalidUserName = errors.New("invalid username")
)

func NewUserName(s string) (UserName, error) {
	if !userNameRegExp.MatchString(s) {
		return UserName{}, ErrInvalidUserName
	}

	return UserName{userName: s}, nil
}

func (u UserName) String() string {
	return u.userName
}

// https://pkg.go.dev/golang.org/x/crypto@v0.45.0/bcrypt#GenerateFromPassword
// GenerateFromPassword does not accept passwords longer than 72 bytes, which is the longest password bcrypt will operate on.
var passwordRegExp = regexp.MustCompile("^[a-zA-Z0-9!@#$%^&*]{8,72}$")

const passwordHashCost = 10

var (
	ErrInvalidPassword = errors.New("invalid password")
)

type HashedPassword struct {
	hash []byte
}

func NewHashedPassword(r RawPassword) (HashedPassword, error) {
	h, err := bcrypt.GenerateFromPassword(r.raw, passwordHashCost)
	if err != nil {
		return HashedPassword{}, fmt.Errorf("failed to generate password hash: %w", err)
	}
	return HashedPassword{hash: h}, nil
}

func NewHashedPasswordFromHash(hash []byte) HashedPassword {
	return HashedPassword{hash}
}

func (h HashedPassword) String() string {
	return string(h.hash)
}

func (h HashedPassword) Bytes() []byte {
	return h.hash
}

func (h HashedPassword) Match(r RawPassword) (bool, error) {
	if err := bcrypt.CompareHashAndPassword(h.hash, r.raw); err != nil {
		return false, nil
	}
	return true, nil
}

type RawPassword struct {
	raw []byte
}

func NewRawPassword(raw []byte) (RawPassword, error) {
	if !passwordRegExp.Match(raw) {
		return RawPassword{}, ErrInvalidPassword
	}

	return RawPassword{raw}, nil
}
