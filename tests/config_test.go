package config

import (
	"testing"

	"github.com/drewschrauf/pmux/config"
	"github.com/stretchr/testify/assert"
)

func newConfig() *config.Config {
	return &config.Config{
		Projects: map[string]config.Project{
			"foo": config.Project{},
			"bar": config.Project{},
			"baz": config.Project{},
		},
	}
}

func TestConfigFilterSingle(t *testing.T) {
	assert := assert.New(t)

	config := newConfig()
	config.Filter("foo", []string{})

	assert.Contains(config.Projects, "foo")
	assert.Len(config.Projects, 1)
}

func TestConfigFilterArray(t *testing.T) {
	assert := assert.New(t)

	config := newConfig()
	config.Filter("", []string{"foo", "bar"})

	assert.Contains(config.Projects, "foo")
	assert.Contains(config.Projects, "bar")
	assert.Len(config.Projects, 2)
}

func TestConfigFilterBoth(t *testing.T) {
	assert := assert.New(t)

	config := newConfig()
	config.Filter("foo", []string{"bar"})

	assert.Contains(config.Projects, "foo")
	assert.Contains(config.Projects, "bar")
	assert.Len(config.Projects, 2)
}

func TestConfigFilterNone(t *testing.T) {
	assert := assert.New(t)

	config := newConfig()
	config.Filter("", []string{})

	assert.Contains(config.Projects, "foo")
	assert.Contains(config.Projects, "bar")
	assert.Contains(config.Projects, "baz")
	assert.Len(config.Projects, 3)
}
