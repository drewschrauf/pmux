package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	yaml "gopkg.in/yaml.v2"

	homedir "github.com/mitchellh/go-homedir"
)

// Config Workspace config
type Config struct {
	// Projects Projects available in the workspace
	Projects map[string]Project `yaml:projects`
}

type Project struct {
	Dir      string             `yaml:dir`
	Commands map[string]Command `yaml:commands`
}

type Command struct {
	Script  string   `yaml:script`
	Depends []string `yaml:depends`
}

// Load Load the default config
func Load() Config {
	dir, err := homedir.Expand("~/.pmux")
	if err != nil {
		fmt.Fprintln(os.Stderr, "Unable to locate config directory")
		os.Exit(1)
	}

	cfgPath := path.Join(dir, "default.yml")
	data, err := ioutil.ReadFile(cfgPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Unable to load config")
		os.Exit(1)
	}

	cfg := Config{}
	yaml.Unmarshal(data, &cfg)

	return cfg
}
