package util

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"

	"github.com/fatih/color"
	shellwords "github.com/mattn/go-shellwords"
)

// Cmd : Settings to run command
type Cmd struct {
	Project  string
	Dir      string
	Script   string
	Requires []string
	Colorize ColorFunc
}

// ColorFunc : Function to colorize output
type ColorFunc func(format string, a ...interface{}) string

// Colors : Static list of available output colors
var Colors = [...]ColorFunc{color.BlueString, color.YellowString, color.MagentaString, color.CyanString, color.GreenString, color.RedString}

// Run : Run a command
func Run(command Cmd) error {
	stop := make(chan bool)
	parts, err := shellwords.Parse(command.Script)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Couldn't parse command", err)
		return err
	}
	cmd := exec.Command(parts[0], parts[1:]...)
	cmd.Dir = command.Dir

	cmdOutReader, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating StdoutPipe", err)
		return err
	}
	outScanner := bufio.NewScanner(cmdOutReader)
	go func() {
		for outScanner.Scan() {
			fmt.Printf("%v - %v\n", command.Colorize(command.Project), outScanner.Text())
		}
		stop <- true
	}()

	cmdErrReader, err := cmd.StderrPipe()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating StdoutPipe", err)
		return err
	}
	errScanner := bufio.NewScanner(cmdErrReader)
	go func() {
		for errScanner.Scan() {
			fmt.Fprintf(os.Stderr, "%v - %v\n", command.Colorize(command.Project), errScanner.Text())
		}
	}()

	err = cmd.Start()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error starting command", err)
		return err
	}

	<-stop
	err = cmd.Wait()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error waiting for command", err)
		return err
	}

	return nil
}

// BuildCmdGroups : Convert array of Cmds into ordered command groups
func BuildCmdGroups(commands []Cmd) [][]Cmd {
	var processed []string
	var groups [][]Cmd

	for len(processed) < len(commands) {
		var group []Cmd
		var thisProcessed []string
		for _, cmd := range commands {
			if !ArrayContains(processed, cmd.Project) && ArrayContainsAll(processed, cmd.Requires) {
				group = append(group, cmd)
				thisProcessed = append(thisProcessed, cmd.Project)
			}
		}
		processed = append(processed, thisProcessed...)
		groups = append(groups, group)
	}

	return groups
}
