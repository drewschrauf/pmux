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

func init() {
	RootCmd.AddCommand(RunCmd)
}

var RunCmd = &cobra.Command{
	Use:   "run",
	Short: "Run a command against all projects in workspace",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.Load()
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
					Project: projectName,
					Dir:     dir,
					Script:  command.Script,
				})
			}
		}

		var wg sync.WaitGroup
		wg.Add(len(commands))
		for i, cmd := range commands {
			go func(cmd util.Cmd, colorize util.ColorFunc) {
				util.Run(cmd, colorize)
				wg.Done()
			}(cmd, util.Colors[i%len(util.Colors)])
		}
		wg.Wait()
	},
}
