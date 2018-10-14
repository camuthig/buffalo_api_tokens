package actions

import (
	"golang.org/x/crypto/bcrypt"
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
	RememberMe bool `json:"remember_me"`
}

// Register a new user into the system and retrieve initial API tokens.
func Register(c buffalo.Context) error {
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

	rt, at, err := createTokens(c, req.RememberMe)

	c.Render(201 ,r.JSON(map[string]interface{}{
		"refresh_token": rt,
		"access_token": at,
		"user": u,
	}))

	return nil
}

type loginRequest struct {
	Email string
	Password string
	RememberMe bool
}

// Login using user credentials and return user information along with API tokens
func Login (c buffalo.Context) error {
	req := &loginRequest{}

	if err := c.Bind(req); err != nil {
		return err
	}

	tx, ok := c.Value("tx").(*pop.Connection)

	if !ok {
		return errors.New("Internal Server Error")
	}

	u := &models.User{}

	if err := tx.Where("email = ?", req.Email).First(u); err != nil {
		return err;
	}

	if err := bcrypt.CompareHashAndPassword(u.PasswordHash, []byte(req.Password)); err != nil {
		c.Render(404, r.JSON(map[string]interface{}{
			"error": map[string]interface{}{
				"message": "The email and password did not match any users in our system.",
			},
		}))

		return nil
	}

	c.Set("user", u)

	rt, at, err := createTokens(c, req.RememberMe)

	if err != nil {
		return err
	}

	c.Render(200 ,r.JSON(map[string]interface{}{
		"refresh_token": rt,
		"access_token": at,
		"user": u,
	}))

	return nil
}

type refreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// Refresh the current user's API token
func Refresh(c buffalo.Context) error {
	req := &refreshRequest{}

	if err := c.Bind(req); err != nil {
		return err
	}

	tx, ok := c.Value("tx").(*pop.Connection)

	if !ok {
		return errors.New("No Transaction")
	}

	// Match the refresh token to the user
	e := &models.RefreshToken{}

	if err := tx.Eager().Where("id = ?", req.RefreshToken).First(e); err != nil {
		return err
	}

	// Destroy the refresh token
	tx.Destroy(e)

	// Create a new API token and refresh token
	c.Set("user", &e.User)

	rt, at, err := createTokens(c, true)

	if err != nil {
		return err
	}

	c.Render(200 ,r.JSON(map[string]interface{}{
		"refresh_token": rt,
		"access_token": at,
	}))

	return nil
}

type logoutRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// Logout a user from the system
func Logout(c buffalo.Context) error {
	req := &refreshRequest{}

	if err := c.Bind(req); err != nil {
		return err
	}

	tx, ok := c.Value("tx").(*pop.Connection)

	if !ok {
		return errors.New("Internal Server Error")
	}

	rt := &models.RefreshToken{
		ID: req.RefreshToken,
	}

	if err := tx.Destroy(rt); err != nil {
		return err
	}

	c.Render(204, nil)

	return nil
}

func createTokens(c buffalo.Context, rm bool) (*string, *string, error) {
	var rt *models.RefreshToken

	u, ok := c.Value("user").(*models.User)

	if !ok {
		return nil, nil, errors.New("Session does not include user")
	}

	if rm {
		tx, ok := c.Value("tx").(*pop.Connection)

		if !ok {
			return nil, nil, errors.New("No transaction found")
		}

		rt = &models.RefreshToken{
			User:   *u,
			UserID: u.ID,
		}

		if err := tx.Create(rt); err != nil {
			return nil, nil, err
		}
	}


	jwt := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject:   u.ID.String(),
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	})

	a, err := jwt.SignedString([]byte(envy.Get("APP_AUTH_KEY", "")))

	if err != nil {
		return nil, nil, err
	}

	if rt != nil {
		return &rt.ID, &a, nil
	}

	return nil, &a, nil
}
