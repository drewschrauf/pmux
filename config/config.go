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
	Script   string   `yaml:"script"`
	Requires []string `yaml:"requires"`
}

// Load : Load the config
func Load(workspace string) Config {
	dir, err := homedir.Expand("~/.pmux")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Unable to locate config directory")
		os.Exit(1)
	}

	cfgPath := path.Join(dir, fmt.Sprintf("%v.yml", workspace))
	data, err := ioutil.ReadFile(cfgPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to find config for workspace '%v'\n", workspace)
		os.Exit(1)
	}

	cfg := Config{}
	yaml.Unmarshal(data, &cfg)

	return cfg
}
