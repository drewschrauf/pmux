package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"pmux/config"
	"sync"

	"github.com/fatih/color"
	shellwords "github.com/mattn/go-shellwords"
	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
)

func init() {
	RootCmd.AddCommand(RunCmd)
}

type Command struct {
	project string
	dir     string
	cmd     string
}

type colorFunc func(format string, a ...interface{}) string

var colors = [...]colorFunc{color.BlueString, color.YellowString, color.MagentaString, color.CyanString, color.GreenString, color.RedString}

var RunCmd = &cobra.Command{
	Use:   "run",
	Short: "Run a command against all projects in workspace",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.Load()
		commandName := args[0]
		var commands []Command

		for projectName, project := range cfg.Projects {
			command, ok := project.Commands[commandName]
			if ok {
				dir, err := homedir.Expand(project.Dir)
				if err != nil {
					fmt.Fprintln(os.Stderr, "Couldn't expand project path", err)
					os.Exit(1)
				}
				commands = append(commands, Command{
					project: projectName,
					dir:     dir,
					cmd:     command,
				})
			}
		}

		var wg sync.WaitGroup
		wg.Add(len(commands))
		for i, cmd := range commands {
			go run(cmd, colors[i%len(colors)], &wg)
		}
		wg.Wait()
	},
}

func run(command Command, colorize colorFunc, wg *sync.WaitGroup) {
	defer wg.Done()
	parts, err := shellwords.Parse(command.cmd)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Couldn't parse command", err)
		os.Exit(1)
	}
	cmd := exec.Command(parts[0], parts[1:]...)
	cmd.Dir = command.dir

	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating StdoutPipe", err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			fmt.Printf("%v - %v\n", colorize(command.project), scanner.Text())
		}
	}()

	err = cmd.Start()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error starting command", err)
		os.Exit(1)
	}

	err = cmd.Wait()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error waiting for command", err)
		os.Exit(1)
	}
}
