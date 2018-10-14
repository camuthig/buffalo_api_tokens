package actions

import (
	"buffalo_api_tokens/models"
	"encoding/json"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gobuffalo/envy"
)

func (as *ActionSuite) Test_Register() {
	res := as.JSON("/auth/register").Post(map[string]interface{}{
		"email":       "napoleon.bonaparte@ziggypig.com",
		"password":    "ziggypiggy",
		"name":        "Napolean Bonaparte",
		"remember_me": true,
	})

	as.Equal(201, res.Code)

	j := map[string]interface{}{}

	if err := json.Unmarshal(res.Body.Bytes(), &j); err != nil {
		as.Fail(err.Error())
	}

	as.Contains(j, "refresh_token")
	as.Contains(j, "access_token")
	as.Contains(j, "user")

	u := j["user"].(map[string]interface{})

	as.Contains(u, "id")
	as.Contains(u, "created_at")
	as.Contains(u, "nickname")
	as.Equal(u["nickname"], "Napolean Bonaparte")
	as.Contains(u, "name")
	as.Equal(u["name"], "Napolean Bonaparte")
	as.Contains(u, "email")
	as.Equal(u["email"], "napoleon.bonaparte@ziggypig.com")

	as.NotContains(u, "password")
}

func (as *ActionSuite) Test_RegisterWithoutRememberMe() {
	res := as.JSON("/auth/register").Post(map[string]interface{}{
		"email":       "napoleon.bonaparte@ziggypig.com",
		"password":    "ziggypiggy",
		"name":        "Napolean Bonaparte",
		"remember_me": false,
	})

	as.Equal(201, res.Code)

	j := map[string]interface{}{}

	if err := json.Unmarshal(res.Body.Bytes(), &j); err != nil {
		as.Fail(err.Error())
	}

	as.Contains(j, "refresh_token")
	as.Nil(j["refresh_token"])
	as.Contains(j, "access_token")
	as.Contains(j, "user")
}

func (as *ActionSuite) Test_Login() {
	as.LoadFixture("users")

	res := as.JSON("/auth/login").Post(map[string]interface{}{
		"email":       "bill@beexcellent.com",
		"password":    "beexcellent",
		"remember_me": true,
	})

	as.Equal(200, res.Code)

	j := map[string]interface{}{}

	if err := json.Unmarshal(res.Body.Bytes(), &j); err != nil {
		as.Fail(err.Error())
	}

	as.Contains(j, "refresh_token")
	as.Contains(j, "access_token")
	as.Contains(j, "user")
}

func (as *ActionSuite) Test_LoginWithoutRememberMe() {
	as.LoadFixture("users")

	res := as.JSON("/auth/login").Post(map[string]interface{}{
		"email":       "bill@beexcellent.com",
		"password":    "beexcellent",
		"remember_me": false,
	})

	as.Equal(200, res.Code)

	j := map[string]interface{}{}

	if err := json.Unmarshal(res.Body.Bytes(), &j); err != nil {
		as.Fail(err.Error())
	}

	as.Contains(j, "refresh_token")
	as.Nil(j["refresh_token"])
	as.Contains(j, "access_token")
	as.Contains(j, "user")
}

func (as *ActionSuite) Test_Refresh() {
	as.LoadFixture("users")

	res := as.JSON("/auth/refresh").Post(map[string]interface{}{
		"refresh_token": "supersecretrefreshtoken",
	})

	as.Equal(200, res.Code)

	j := map[string]interface{}{}

	if err := json.Unmarshal(res.Body.Bytes(), &j); err != nil {
		as.Fail(err.Error())
	}

	as.Contains(j, "refresh_token")
	as.Contains(j, "access_token")

	// Ensure the column has been deleted
	r := &models.RefreshToken{}

	err := as.DB.Where("id = ?", "supersecretrefreshtoken").First(r)

	as.NotNil(err)
}

func (as *ActionSuite) Test_RefreshWithBadToken() {
	as.LoadFixture("users")

	res := as.JSON("/auth/refresh").Post(map[string]interface{}{
		"refresh_token": "junk",
	})

	as.Equal(404, res.Code)
}

func (as *ActionSuite) Test_Logout() {
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

	res := as.JSON("/auth/logout").Post(map[string]interface{}{
		"refresh_token": "supersecretrefreshtoken",
	})

	as.Equal(204, res.Code)

	// Ensure the column has been deleted
	r := &models.RefreshToken{}

	err = as.DB.Where("id = ?", "supersecretrefreshtoken").First(r)

	as.NotNil(err)
}

func (as *ActionSuite) Test_LogoutWithoutFindingToken() {
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

	res := as.JSON("/auth/logout").Post(map[string]interface{}{
		"refresh_token": "anonexistingtoken",
	})

	as.Equal(204, res.Code)
}

func (as *ActionSuite) Test_LogoutAllDevices() {
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

	res := as.JSON("/auth/logout").Post(map[string]interface{}{})

	as.Equal(204, res.Code)

	// Ensure the column has been deleted
	r := []models.RefreshToken{}

	err = as.DB.Where("user_id = ?", u.ID).All(r)

	as.NotNil(err)
}

func (as *ActionSuite) Test_LogoutWithOtherUsersToken() {
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

	res := as.JSON("/auth/logout").Post(map[string]interface{}{
		"refresh_token": "tedssupersecretrefreshtoken",
	})

	as.Equal(204, res.Code)

	// Ensure the column has been deleted
	r := &models.RefreshToken{}

	err = as.DB.Where("id = ?", "supersecretrefreshtoken").First(r)

	as.Nil(err)

	count, _ := as.DB.Count(models.RefreshToken{})

	as.Equal(2, count)

	err = as.DB.Where("id = ?", "tedssupersecretrefreshtoken").First(r)

	as.Nil(err)
}
