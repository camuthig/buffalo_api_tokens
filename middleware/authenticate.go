package middleware

import (
	"errors"
	"fmt"
	"time"

	"buffalo_api_tokens/models"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/envy"
	"github.com/gobuffalo/pop"
)

// Authenticate a user based on passed in API credentials and add it to the context
func Authenticate(next buffalo.Handler) buffalo.Handler {
	return func(c buffalo.Context) error {
		j, err := request.ParseFromRequest(
			c.Request(),
			request.AuthorizationHeaderExtractor,
			func(t *jwt.Token) (interface{}, error) {
				// Validate the algorithm matches the expected
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("Unexpected signing method: %v", t.Header["alg"])
				}

				return []byte(envy.Get("APP_AUTH_KEY", "")), nil
			},
			request.WithClaims(&jwt.StandardClaims{}),
		)

		if err != nil {
			return c.Error(401, err)
		}

		claims, ok := j.Claims.(*jwt.StandardClaims)

		if !ok {
			return c.Error(401, errors.New("Invalid access token"))
		}

		// Explicitly check expiration date and require it
		if !claims.VerifyExpiresAt(time.Now().Unix(), true) {
			return c.Error(401, errors.New("Invalid access token"))
		}

		tx, ok := c.Value("tx").(*pop.Connection)

		if !ok {
			return errors.New("Internal Server Error")
		}

		u := &models.User{}

		if err := tx.Find(u, claims.Subject); err != nil {
			return c.Error(401, errors.New("Invalid access token"))
		}

		c.Set("user", u)

		return next(c)
	}
}
