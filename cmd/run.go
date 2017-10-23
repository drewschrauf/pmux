package cmd

import (
	"fmt"
	"os"
	"pmux/config"
	"pmux/util"
	"sync"

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
)

var skipRequires bool

func init() {
	runCmd.Flags().BoolVarP(&skipRequires, "skip-requires", "s", false, "Don't enforce requires")
	RootCmd.AddCommand(runCmd)
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Run a command against all projects in workspace",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.Load(Workspace)
		commandName := args[0]
		var commands []util.Cmd

		for projectName, project := range cfg.Projects {
			command, ok := project.Commands[commandName]
			if ok {
				dir, err := homedir.Expand(project.Dir)
				if err != nil {
					fmt.Fprintln(os.Stderr, "Couldn't expand project path", err)
					os.Exit(1)
				}
				commands = append(commands, util.Cmd{
					Project:  projectName,
					Dir:      dir,
					Script:   command.Script,
					Requires: command.Requires,
					Colorize: util.Colors[len(commands)%len(util.Colors)],
				})
			}
		}

		var cmdGroups [][]util.Cmd
		if skipRequires {
			cmdGroups = append(cmdGroups, commands)
		} else {
			cmdGroups = buildCmdGroups(commands)
		}

		for _, group := range cmdGroups {
			var wg sync.WaitGroup
			wg.Add(len(group))

			for _, cmd := range group {
				go func(cmd util.Cmd) {
					err := util.Run(cmd)
					if err != nil {
						fmt.Fprintf(os.Stderr, "Error running command for '%v'", cmd.Project)
						os.Exit(1)
					}
					wg.Done()
				}(cmd)
			}

			wg.Wait()
		}
	},
}

func buildCmdGroups(commands []util.Cmd) [][]util.Cmd {
	var processed []string
	var groups [][]util.Cmd

	for len(processed) < len(commands) {
		var group []util.Cmd
		var thisProcessed []string
		for _, cmd := range commands {
			if !contains(processed, cmd.Project) && containsAll(processed, cmd.Requires) {
				group = append(group, cmd)
				thisProcessed = append(thisProcessed, cmd.Project)
			}
		}
		processed = append(processed, thisProcessed...)
		groups = append(groups, group)
	}

	return groups
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func containsAll(set []string, subset []string) bool {
	var found []string
	for _, e := range subset {
		if contains(set, e) {
			found = append(found, e)
		}
	}
	return len(found) == len(subset)
}
