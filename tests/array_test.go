package pmux_test

import (
	"testing"

	"github.com/drewschrauf/pmux/util"
	"github.com/stretchr/testify/assert"
)

func TestArrayContainsTrue(t *testing.T) {
	assert := assert.New(t)

	slice := []string{"foo", "bar", "baz"}
	assert.True(util.ArrayContains(slice, "bar"))
}

func TestArrayContainsFalse(t *testing.T) {
	assert := assert.New(t)

	slice := []string{"foo", "bar", "baz"}
	assert.False(util.ArrayContains(slice, "foobar"))
}

func TestArrayContainsAllTrue(t *testing.T) {
	assert := assert.New(t)

	slice := []string{"foo", "bar", "baz"}
	sub := []string{"foo", "baz"}
	assert.True(util.ArrayContainsAll(slice, sub))
}

func TestArrayContainsAllFalse(t *testing.T) {
	assert := assert.New(t)

	slice := []string{"foo", "bar", "baz"}
	sub := []string{"foo", "baz", "foobar"}
	assert.False(util.ArrayContainsAll(slice, sub))
}
