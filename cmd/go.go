package cmd

import (
	"fmt"
	"os"
	"pmux/config"

	"github.com/spf13/cobra"
)

// func init() {
// 	RootCmd.AddCommand(GoCmd)
// }

var GoCmd = &cobra.Command{
	Use:   "go",
	Short: "Go to a project in the workspace",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.Load()
		projectName := args[0]

		project := cfg.Projects[projectName]
		err := os.Chdir(project.Dir)
		if err != nil {
			fmt.Println("Unable to change to project")
			os.Exit(1)
		}
	},
}
