package cmd

import (
	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:   "pmux",
	Short: "pmux is a manager for projects spread across multiple repositories",
}
