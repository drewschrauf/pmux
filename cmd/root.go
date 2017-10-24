package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var Workspace string
var Project string
var Projects []string

func init() {
	workspace := os.Getenv("PMUX_WORKSPACE")
	if workspace == "" {
		workspace = "default"
	}
	RootCmd.PersistentFlags().StringVarP(&Workspace, "workspace", "w", workspace, "set workspace to use")
	RootCmd.PersistentFlags().StringVarP(&Project, "project", "p", "", "run command against single project")
	RootCmd.PersistentFlags().StringArrayVarP(&Projects, "projects", "m", []string{}, "run command against multiple project")
}

var RootCmd = &cobra.Command{
	Use:   "pmux",
	Short: "pmux is a manager for projects spread across multiple repositories",
}
