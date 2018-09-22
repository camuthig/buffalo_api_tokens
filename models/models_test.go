package models_test

import (
	"testing"

	"github.com/gobuffalo/packr"
	"github.com/gobuffalo/suite"
	"github.com/gobuffalo/validate"
)

type ModelSuite struct {
	*suite.Model
}

// AssertContainsErrorKeys will ensure that the collection of validation errors includes the given keys.
func (ms *ModelSuite) AssertContainsErrorKeys(e validate.Errors, keys []string) bool {
	f := make([]string, 0, len(e.Errors))
	for k := range e.Errors {
		println(k)
		f = append(f, k)
	}

	return ms.Assert().ElementsMatch(f, keys)
}

func Test_ModelSuite(t *testing.T) {
	model, err := suite.NewModelWithFixtures(packr.NewBox("../fixtures"))
	if err != nil {
		t.Fatal(err)
	}

	as := &ModelSuite{
		Model: model,
	}
	suite.Run(t, as)
}
