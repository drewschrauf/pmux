package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var Workspace string

func init() {
	workspace := os.Getenv("PMUX_WORKSPACE")
	if workspace == "" {
		workspace = "default"
	}
	RootCmd.PersistentFlags().StringVarP(&Workspace, "workspace", "w", workspace, "Set workspace to use")
}

var RootCmd = &cobra.Command{
	Use:   "pmux",
	Short: "pmux is a manager for projects spread across multiple repositories",
}
