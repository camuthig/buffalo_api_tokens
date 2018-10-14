package actions

import (
	"buffalo_api_tokens/models"
	"errors"

	"github.com/gobuffalo/buffalo"
)

// Me will return the logged in user's information
func Me(c buffalo.Context) error {
	u, ok := c.Value("user").(*models.User)

	if !ok {
		return errors.New("Internal Server Error")
	}

	c.Render(200, r.JSON(u))

	return nil
}
