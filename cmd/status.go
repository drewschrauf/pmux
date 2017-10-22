package cmd

import (
	"fmt"
	"os"
	"pmux/config"
	"pmux/git"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
)

func init() {
	RootCmd.AddCommand(StatusCmd)
}

var StatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Display git status of all projects in workspace",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.Load()

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Project", "Branch", "Dirty", "Ahead", "Behind"})

		for k, v := range cfg.Projects {
			dir, err := homedir.Expand(v.Dir)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Couldn't expand project path", err)
				os.Exit(1)
			}

			branch, err := git.Branch(dir)
			if err != nil {
				branch = "Error"
			}

			var dirtyStr = "No"
			dirty, err := git.Dirty(dir)
			if err != nil {
				dirtyStr = color.RedString("Error")
			} else if dirty {
				dirtyStr = color.RedString("Yes")
			}

			var aheadStr = "No"
			ahead, err := git.Ahead(dir)
			if err != nil {
				aheadStr = color.RedString("Error")
			} else if ahead {
				aheadStr = color.RedString("Yes")
			}

			var behindStr = "No"
			behind, err := git.Behind(dir)
			if err != nil {
				behindStr = color.RedString("Error")
			} else if behind {
				behindStr = color.RedString("Yes")
			}

			table.Append([]string{k, branch, dirtyStr, aheadStr, behindStr})
		}

		table.Render()
	},
}
