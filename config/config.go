package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	yaml "gopkg.in/yaml.v2"

	homedir "github.com/mitchellh/go-homedir"
)

// Config : Workspace config
type Config struct {
	Projects map[string]Project `yaml:"projects"`
}

// Project : Project config
type Project struct {
	Dir      string             `yaml:"dir"`
	Commands map[string]Command `yaml:"commands"`
}

// Command : Command config
type Command struct {
	Script  string   `yaml:"script"`
	Depends []string `yaml:"depends"`
}

// Load : Load the config
func Load() Config {
	dir, err := homedir.Expand("~/.pmux")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Unable to locate config directory")
		os.Exit(1)
	}

	project := os.Getenv("PMUX_WORKSPACE")
	if project == "" {
		project = "default"
	}

	cfgPath := path.Join(dir, fmt.Sprintf("%v.yml", project))
	data, err := ioutil.ReadFile(cfgPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to find config for workspace '%v'\n", project)
		os.Exit(1)
	}

	cfg := Config{}
	yaml.Unmarshal(data, &cfg)

	return cfg
}
