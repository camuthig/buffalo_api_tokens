package middleware_test

import (
	"testing"

	"buffalo_api_tokens/actions"

	"github.com/gobuffalo/packr"
	"github.com/gobuffalo/suite"
)

type MiddlewareSuite struct {
	*suite.Action
}

func Test_MiddlewareSuite(t *testing.T) {
	action, err := suite.NewActionWithFixtures(actions.App(), packr.NewBox("../fixtures"))
	if err != nil {
		t.Fatal(err)
	}

	ms := &MiddlewareSuite{
		Action: action,
	}
	suite.Run(t, ms)
}
