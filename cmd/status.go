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

				var dirtyStr = "No"
				dirty, err := util.GitDirty(dir)
				if err != nil {
					dirtyStr = color.RedString("Error")
				} else if dirty {
					dirtyStr = color.RedString("Yes")
				}

				var aheadStr = "No"
				ahead, err := util.GitAhead(dir)
				if err != nil {
					aheadStr = color.RedString("Error")
				} else if ahead {
					aheadStr = color.RedString("Yes")
				}

				var behindStr = "No"
				behind, err := util.GitBehind(dir)
				if err != nil {
					behindStr = color.RedString("Error")
				} else if behind {
					behindStr = color.RedString("Yes")
				}

				statuses = append(statuses, Status{
					Project: projectName,
					Branch:  branch,
					Dirty:   dirtyStr,
					Ahead:   aheadStr,
					Behind:  behindStr,
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
