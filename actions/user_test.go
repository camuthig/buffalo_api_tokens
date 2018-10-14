package actions

import (
	"buffalo_api_tokens/models"
	"encoding/json"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gobuffalo/envy"
)

func (as *ActionSuite) Test_Me() {
	as.LoadFixture("users")

	u := &models.User{}

	if err := as.DB.Where("email = ?", "bill@beexcellent.com").First(u); err != nil {
		as.Fail(err.Error())
	}

	jwt := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Subject:   u.ID.String(),
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
	})

	a, err := jwt.SignedString([]byte(envy.Get("APP_AUTH_KEY", "")))

	if err != nil {
		as.Fail(err.Error())
	}

	as.Willie.Headers["Authorization"] = a

	res := as.JSON("/me").Get()

	as.Equal(200, res.Code)

	j := map[string]interface{}{}

	if err := json.Unmarshal(res.Body.Bytes(), &j); err != nil {
		as.Fail(err.Error())
	}

	as.Equal(j["id"], u.ID.String())
}

func (as *ActionSuite) Test_MeWithoutToken() {
	res := as.JSON("/me").Get()

	as.Equal(401, res.Code)
}
