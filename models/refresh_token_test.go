package models_test

import (
	"buffalo_api_tokens/models"

	"github.com/gobuffalo/uuid"
)

func (ms *ModelSuite) Test_RefreshTokensCreateSecureIDs() {
	u := &models.User{
		ID: uuid.Must(uuid.NewV4()),
	}

	rt := &models.RefreshToken{
		User: *u,
	}

	err := rt.BeforeCreate(ms.DB)

	ms.Assert().Nil(err)

	ms.Assert().NotNil(rt.ID)
	ms.Assert().Len(rt.ID, 80)
}
