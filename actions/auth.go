package actions

import (
	"errors"
	"time"

	"buffalo_api_tokens/models"

	"github.com/dgrijalva/jwt-go"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/pop"
	"github.com/gobuffalo/uuid"
)

type registerRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

// RegisterHandler supports creating new users.
func RegisterHandler(c buffalo.Context) error {
	req := &registerRequest{}

	if err := c.Bind(req); err != nil {
		return err
	}

	u := &models.User{
		ID:       uuid.Must(uuid.NewV4()),
		Email:    req.Email,
		Password: req.Password,
		Name:     req.Name,
		Nickname: req.Name,
	}

	tx, ok := c.Value("tx").(*pop.Connection)

	if !ok {
		return errors.New("Internal Server Error")
	}

	verr, err := tx.ValidateAndCreate(u)

	if verr != nil {
		r.Auto(c, r.JSON(verr.Errors))
	}

	if err != nil {
		return err
	}

	c.Set("user", u)

	rt, at, err := createTokens(c)

	c.Render(201 ,r.JSON(map[string]interface{}{
		"refresh_token": rt.ID,
		"access_token": at,
		"user": u,
	}))

	return nil
}

func createTokens(c buffalo.Context) (*models.RefreshToken, string, error) {
	u, ok := c.Value("user").(*models.User)

	if !ok {
		c.Render(204, nil)
	}

	tx, ok := c.Value("tx").(*pop.Connection)

	if !ok {
		return nil, "", errors.New("No transaction found")
	}

	rt := &models.RefreshToken{
		User:   *u,
		UserID: u.ID,
	}

	if err := tx.Create(rt); err != nil {
		return nil, "", err
	}

	jwt := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject:   u.ID.String(),
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	})

	a, err := jwt.SignedString([]byte(envy.Get("APP_AUTH_KEY", "")))

	if err != nil {
		return nil, "", err
	}

	return rt, a, nil
}
