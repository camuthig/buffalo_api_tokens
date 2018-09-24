package actions

import (
	"encoding/json"
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
