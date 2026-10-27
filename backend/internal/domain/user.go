package domain

import (
	"errors"
	"strings"
	"time"
)

// User is a person who logged in via Google. GoogleID (the OAuth "sub"
// claim) is the stable identity; Email/Name/AvatarURL are a cache of what
// Google reported at last login, refreshed on every login.
type User struct {
	ID        uint   `gorm:"primaryKey"`
	GoogleID  string `gorm:"uniqueIndex;size:255;not null"`
	Email     string `gorm:"size:255;not null"`
	Name      string `gorm:"size:255;not null"`
	AvatarURL string `gorm:"size:1024"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

var (
	ErrInvalidGoogleID = errors.New("domain: google id must not be empty")
	ErrInvalidEmail    = errors.New("domain: email must not be empty")
)

// Validate enforces the invariants a User must hold before it is persisted.
func (u *User) Validate() error {
	if strings.TrimSpace(u.GoogleID) == "" {
		return ErrInvalidGoogleID
	}
	if strings.TrimSpace(u.Email) == "" {
		return ErrInvalidEmail
	}
	return nil
}
