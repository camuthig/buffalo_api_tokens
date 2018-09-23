package actions

import (
	"encoding/json"
)

func (as *ActionSuite) Test_RegisterHandler() {
	res := as.JSON("/auth/register").Post(map[string]interface{}{
		"email":    "bob.barker@thepriceisright.com",
		"password": "thepriceiswrong",
		"name":     "Bob Barker",
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
	as.Equal(u["nickname"], "Bob Barker")
	as.Contains(u, "name")
	as.Equal(u["name"], "Bob Barker")
	as.Contains(u, "email")
	as.Equal(u["email"], "bob.barker@thepriceisright.com")

	as.NotContains(u, "password")
}
