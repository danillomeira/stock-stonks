package domain_test

import (
	"testing"

	"github.com/danillomeira/stock-stonks/backend/internal/domain"
)

func TestUser_Validate(t *testing.T) {
	cases := []struct {
		name string
		user domain.User
		want error
	}{
		{"empty google id", domain.User{Email: "a@example.com"}, domain.ErrInvalidGoogleID},
		{"empty email", domain.User{GoogleID: "google-123"}, domain.ErrInvalidEmail},
		{"valid", domain.User{GoogleID: "google-123", Email: "a@example.com"}, nil},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if err := tc.user.Validate(); err != tc.want {
				t.Fatalf("Validate() = %v, want %v", err, tc.want)
			}
		})
	}
}
