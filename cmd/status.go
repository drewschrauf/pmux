package cmd

import (
	"os"
	"pmux/config"
	"pmux/git"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
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
			branch, err := git.Branch(v.Dir)
			if err != nil {
				branch = "Error"
			}

			var dirtyStr = "No"
			dirty, err := git.Dirty(v.Dir)
			if err != nil {
				dirtyStr = "Error"
			} else if dirty {
				dirtyStr = "Yes"
			}

			var aheadStr = "No"
			ahead, err := git.Ahead(v.Dir)
			if err != nil {
				aheadStr = "Error"
			} else if ahead {
				aheadStr = "Yes"
			}

			var behindStr = "No"
			behind, err := git.Behind(v.Dir)
			if err != nil {
				behindStr = "Error"
			} else if behind {
				behindStr = "Yes"
			}

			table.Append([]string{k, branch, dirtyStr, aheadStr, behindStr})
		}

		table.Render()
	},
}
