package models_test

import (
	"buffalo_api_tokens/models"

	"github.com/gobuffalo/uuid"
	"golang.org/x/crypto/bcrypt"
)

func (ms *ModelSuite) Test_BeforeSave() {
	u := &models.User{
		ID:       uuid.Must(uuid.NewV4()),
		Name:     "Napolean Bonaparte",
		Nickname: "Napolean",
		Email:    "napolean.bonaparte@ziggypig.com",
		Password: "ziggypiggy",
	}

	err := u.BeforeSave(ms.DB)

	ms.Assert().Nil(err)

	ms.Assert().Nil(bcrypt.CompareHashAndPassword(u.PasswordHash, []byte("ziggypiggy")))
}

func (ms *ModelSuite) Test_Validate() {
	u := &models.User{
		ID:       uuid.Must(uuid.NewV4()),
		Name:     "",
		Nickname: "",
		Email:    "garbage",
		Password: "short",
	}

	verr, err := u.Validate(ms.DB)

	ms.Assert().Nil(err)
	ms.Assert().NotNil(verr)
	ms.Assert().Len(verr.Errors, 4)
	ms.Assert().Contains(verr.Errors, "email")
	ms.Assert().Contains(verr.Errors, "name")
	ms.Assert().Contains(verr.Errors, "nickname")
	ms.Assert().Contains(verr.Errors, "password")
}
