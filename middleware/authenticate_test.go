package middleware_test

import (
	"buffalo_api_tokens/middleware"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/envy"
	"buffalo_api_tokens/models"
)

func (ms *MiddlewareSuite) Test_Authenticate() {
	ms.LoadFixture("users")

	ms.App.Use(middleware.Authenticate)

	ms.App.GET("/test/authenticate", func(c buffalo.Context) error {
		u, ok := c.Value("user").(*models.User)

		if !ok {
			ms.Fail("User not added to context")
		}

		ms.NotNil(u)

		c.Render(200, nil)

		return nil
	})

	u := &models.User{}

	if err := ms.DB.Where("email = ?", "bill@beexcellent.com").First(u); err != nil {
		ms.Fail(err.Error())
	}

	jwt := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject:   u.ID.String(),
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	})

	a, err := jwt.SignedString([]byte(envy.Get("APP_AUTH_KEY", "")))

	if err != nil {
		ms.Fail(err.Error())
	}

	ms.Willie.Headers["Authorization"] = a

	res := ms.JSON("/test/authenticate").Get()

	ms.Equal(200, res.Code)
}

func (ms *MiddlewareSuite) Test_AuthenticateExpiredToken() {
	ms.LoadFixture("users")

	ms.App.Use(middleware.Authenticate)

	ms.App.GET("/test/authenticate", func(c buffalo.Context) error {
		return nil
	})

	u := &models.User{}

	if err := ms.DB.Where("email = ?", "bill@beexcellent.com").First(u); err != nil {
		ms.Fail(err.Error())
	}

	jwt := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject:   u.ID.String(),
		ExpiresAt: time.Now().Add(time.Hour * -24).Unix(),
	})

	a, err := jwt.SignedString([]byte(envy.Get("APP_AUTH_KEY", "")))

	if err != nil {
		ms.Fail(err.Error())
	}

	ms.Willie.Headers["Authorization"] = a

	res := ms.JSON("/test/authenticate").Get()

	ms.Equal(401, res.Code)
}

func (ms *MiddlewareSuite) Test_AuthenticateRequiresExpiration() {
	ms.LoadFixture("users")

	ms.App.Use(middleware.Authenticate)

	ms.App.GET("/test/authenticate", func(c buffalo.Context) error {
		return nil
	})

	u := &models.User{}

	if err := ms.DB.Where("email = ?", "bill@beexcellent.com").First(u); err != nil {
		ms.Fail(err.Error())
	}

	jwt := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject:   u.ID.String(),
	})

	a, err := jwt.SignedString([]byte(envy.Get("APP_AUTH_KEY", "")))

	if err != nil {
		ms.Fail(err.Error())
	}

	ms.Willie.Headers["Authorization"] = a

	res := ms.JSON("/test/authenticate").Get()

	ms.Equal(401, res.Code)
}

func (ms *MiddlewareSuite) Test_AuthenticateRequiresToken() {
	ms.App.Use(middleware.Authenticate)

	ms.App.GET("/test/authenticate", func(c buffalo.Context) error {
		return nil
	})

	res := ms.JSON("/test/authenticate").Get()

	ms.Equal(401, res.Code)
}
