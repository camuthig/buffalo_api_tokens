package models

import (
	"time"

	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
	"github.com/gobuffalo/validate"
	"github.com/gobuffalo/validate/validators"
	"golang.org/x/crypto/bcrypt"
)

// A User gives our application a way to track accounts in the application.
type User struct {
	ID           uuid.UUID `json:"id" db:"id"`
	Email        string    `json:"email" db:"email"`
	Name         string    `json:"name" db:"name"`
	Nickname     string    `json:"nickname" db:"nickname"`
	Password     string    `json:"-" db:"-"`
	PasswordHash []byte    `json:"-" db:"password_hash"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

func (u *User) Validate(tx *pop.Connection) (*validate.Errors, error) {
	v := []validate.Validator{
		&validators.EmailIsPresent{Field: u.Email, Name: "Email"},
		&validators.StringIsPresent{Field: u.Name, Name: "Name"},
		&validators.StringIsPresent{Field: u.Nickname, Name: "Nickname"},
	}

	if u.Password != "" {
		v = append(v, &validators.StringLengthInRange{Field: u.Password, Name: "Password", Min: 8})
	}

	return validate.Validate(
		v...,
	), nil
}

func (u *User) BeforeSave(tx *pop.Connection) error {
	if u.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)

		if err != nil {
			return err
		}

		u.PasswordHash = hash
	}

	return nil
}
