package cmd

import (
	"fmt"
	"os"
	"sort"
	"sync"

	"github.com/drewschrauf/pmux/config"
	"github.com/drewschrauf/pmux/util"
	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
)

type Status struct {
	Project string
	Branch  string
	Dirty   string
	Ahead   string
	Behind  string
}

func init() {
	RootCmd.AddCommand(statusCmd)
}

type gitFlagCmd func(dir string) (bool, error)

func runGitFlagCmd(cmd gitFlagCmd, dir string) string {
	var resultStr = "No"
	dirty, err := cmd(dir)
	if err != nil {
		resultStr = color.RedString("Error")
	} else if dirty {
		resultStr = color.RedString("Yes")
	}
	return resultStr
}

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Display git status of all projects in workspace",
	Args:  cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.Load(Workspace)
		cfg.Filter(Project, Projects)

		var wg sync.WaitGroup
		wg.Add(len(cfg.Projects))

		var statuses []Status
		for projectName, project := range cfg.Projects {
			go func(projectName string, project config.Project) {
				defer wg.Done()

				dir, err := homedir.Expand(project.Dir)
				if err != nil {
					fmt.Fprintln(os.Stderr, "Couldn't expand project path", err)
					os.Exit(1)
				}

				branch, err := util.GitBranch(dir)
				if err != nil {
					branch = "Error"
				}

				var dirty = runGitFlagCmd(util.GitDirty, dir)
				var ahead = runGitFlagCmd(util.GitAhead, dir)
				var behind = runGitFlagCmd(util.GitBehind, dir)

				statuses = append(statuses, Status{
					Project: projectName,
					Branch:  branch,
					Dirty:   dirty,
					Ahead:   ahead,
					Behind:  behind,
				})
			}(projectName, project)
		}
		wg.Wait()

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Project", "Branch", "Dirty", "Ahead", "Behind"})

		sort.Slice(statuses, func(i, j int) bool {
			return statuses[i].Project < statuses[j].Project
		})
		for _, status := range statuses {
			table.Append([]string{
				status.Project, status.Branch, status.Dirty, status.Ahead, status.Behind,
			})
		}

		table.Render()
	},
}
