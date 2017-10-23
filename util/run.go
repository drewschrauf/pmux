package util

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"

	"github.com/fatih/color"
	shellwords "github.com/mattn/go-shellwords"
)

type Cmd struct {
	Project string
	Dir     string
	Script  string
}

type ColorFunc func(format string, a ...interface{}) string

var Colors = [...]ColorFunc{color.BlueString, color.YellowString, color.MagentaString, color.CyanString, color.GreenString, color.RedString}

func Run(command Cmd, colorize ColorFunc) {
	parts, err := shellwords.Parse(command.Script)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Couldn't parse command", err)
		os.Exit(1)
	}
	cmd := exec.Command(parts[0], parts[1:]...)
	cmd.Dir = command.Dir

	cmdOutReader, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating StdoutPipe", err)
		os.Exit(1)
	}
	outScanner := bufio.NewScanner(cmdOutReader)
	go func() {
		for outScanner.Scan() {
			fmt.Printf("%v - %v\n", colorize(command.Project), outScanner.Text())
		}
	}()

	cmdErrReader, err := cmd.StderrPipe()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error creating StdoutPipe", err)
		os.Exit(1)
	}
	errScanner := bufio.NewScanner(cmdErrReader)
	go func() {
		for errScanner.Scan() {
			fmt.Fprintf(os.Stderr, "%v - %v\n", colorize(command.Project), errScanner.Text())
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
