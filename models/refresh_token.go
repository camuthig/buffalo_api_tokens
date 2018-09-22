package models

import (
	"strings"
	"time"

	"crypto/rand"
	"encoding/base64"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
)

// A RefreshToken allows the user to stay logged in for longer periods of time.
type RefreshToken struct {
	ID        string    `json:"id" db:"id"`
	UserID    uuid.UUID `json:"-" db:"user_id"`
	User      User      `json:"user" db:"-"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

func (r *RefreshToken) BeforeCreate(ts *pop.Connection) error {
	b := make([]byte, 36)

	_, err := rand.Read(b)

	if err != nil {
		return err
	}

	t := base64.URLEncoding.EncodeToString(b)

	r.ID = t + strings.Replace(r.User.ID.String(), "-", "", -1)

	return nil
}
