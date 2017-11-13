package pmux_test

import (
	"testing"

	"github.com/drewschrauf/pmux/util"
	"github.com/stretchr/testify/assert"
)

func TestNoDepsSingleGroup(t *testing.T) {
	assert := assert.New(t)
	cmds := []util.Cmd{util.Cmd{}, util.Cmd{}, util.Cmd{}}

	cmdGroups := util.BuildCmdGroups(cmds)
	assert.Len(cmdGroups, 1)
	assert.Len(cmdGroups[0], 3)
}

func TestDepsMultipleGroups(t *testing.T) {
	assert := assert.New(t)
	cmds := []util.Cmd{
		util.Cmd{
			Project:  "foo",
			Requires: []string{"baz"},
		},
		util.Cmd{
			Project:  "bar",
			Requires: []string{"baz"},
		},
		util.Cmd{
			Project: "baz",
		},
	}

	cmdGroups := util.BuildCmdGroups(cmds)
	assert.Len(cmdGroups, 2)
	assert.Len(cmdGroups[0], 1)
	assert.Len(cmdGroups[1], 2)
}
