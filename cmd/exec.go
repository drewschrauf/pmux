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

var forceExec bool

func init() {
	execCmd.Flags().BoolVarP(&forceExec, "force", "f", false, "continue on errors")
	RootCmd.AddCommand(execCmd)
}

var execCmd = &cobra.Command{
	Use:   "exec [command]",
	Short: "Run an arbitrary command against all projects in workspace",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.Load(Workspace)
		cfg.Filter(Project, Projects)

		script := args[0]
		var commands []util.Cmd

		for projectName, project := range cfg.Projects {
			dir, err := homedir.Expand(project.Dir)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Couldn't expand project path", err)
				os.Exit(1)
			}
			commands = append(commands, util.Cmd{
				Project:  projectName,
				Dir:      dir,
				Script:   script,
				Colorize: util.Colors[len(commands)%len(util.Colors)],
			})
		}

		var wg sync.WaitGroup
		wg.Add(len(commands))

		for _, cmd := range commands {
			go func(cmd util.Cmd) {
				err := util.Run(cmd)
				if err != nil {
					fmt.Fprintf(os.Stderr, "Error running command for '%v'\n", cmd.Project)
					if !forceExec {
						fmt.Fprintf(os.Stderr, "Exiting. Use flag --force to ignore.\n")
						os.Exit(1)
					}
				}
				wg.Done()
			}(cmd)
		}

		wg.Wait()
	},
}
