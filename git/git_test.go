// Package git will manage all requirements to clone repository
package git

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClone_fail(t *testing.T) {
	assert := assert.New(t)
	c := &Repository{}

	_, _, err := c.Clone()
	assert.Error(err)
}

func TestClone_fail_branch(t *testing.T) {
	assert := assert.New(t)
	c := &Repository{}

	c.Repository = "https://github.com/cypress-io/cypress-example-kitchensink.git"
	c.Ref = "test"

	z, _, err := c.Clone()
	defer os.RemoveAll(z)
	assert.Error(err)
}

func TestClone_fail_user(t *testing.T) {
	assert := assert.New(t)
	c := &Repository{}

	c.Repository = "https://github.com/cypress-io/cypress-example-kitchensink.git"
	c.Username = "test"

	z, _, err := c.Clone()
	defer os.RemoveAll(z)
	assert.Error(err)
}

func TestClone_fail_user_branch(t *testing.T) {
	assert := assert.New(t)
	c := &Repository{}

	c.Repository = "https://github.com/cypress-io/cypress-example-kitchensink.git"
	c.Username = "test"
	c.Ref = "test"

	z, _, err := c.Clone()
	defer os.RemoveAll(z)
	assert.Error(err)
}

func TestClone_success(t *testing.T) {
	assert := assert.New(t)
	c := &Repository{}

	c.Repository = "https://github.com/cypress-io/cypress-example-kitchensink.git"

	z, _, err := c.Clone()
	defer os.RemoveAll(z)
	assert.Nil(err)
}

func TestClone_success_master(t *testing.T) {
	assert := assert.New(t)
	c := &Repository{}

	c.Repository = "https://github.com/cypress-io/cypress-example-kitchensink.git"
	c.Ref = "master"

	z, _, err := c.Clone()
	defer os.RemoveAll(z)
	assert.Nil(err)
}
