package cmd

import (
	"fmt"
	"os"
	"sync"

	"github.com/drewschrauf/pmux/config"
	"github.com/drewschrauf/pmux/util"

	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
)

var skipRequires bool
var force bool

func init() {
	runCmd.Flags().BoolVarP(&skipRequires, "skip-requires", "s", false, "don't enforce requires")
	runCmd.Flags().BoolVarP(&force, "force", "f", false, "continue on errors")
	RootCmd.AddCommand(runCmd)
}

var runCmd = &cobra.Command{
	Use:   "run [command]",
	Short: "Run a configured command against all projects in workspace",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.Load(Workspace)
		cfg.Filter(Project, Projects)

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
			cmdGroups = util.BuildCmdGroups(commands)
		}

		for _, group := range cmdGroups {
			var wg sync.WaitGroup
			wg.Add(len(group))

			for _, cmd := range group {
				go func(cmd util.Cmd) {
					err := util.Run(cmd)
					if err != nil {
						fmt.Fprintf(os.Stderr, "Error running command for '%v'\n", cmd.Project)
						if !force {
							fmt.Fprintf(os.Stderr, "Exiting. Use flag --force to ignore.\n")
							os.Exit(1)
						}
					}
					wg.Done()
				}(cmd)
			}

			wg.Wait()
		}
	},
}
