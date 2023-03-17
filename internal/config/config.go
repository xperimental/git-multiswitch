package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/pflag"
)

type Config struct {
	GitCmd   string
	BasePath string
	Branch   string
}

func Parse(cmd string, args []string) (*Config, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("can not get working directory: %w", err)
	}

	cfg := &Config{
		GitCmd:   "git",
		BasePath: wd,
	}

	flags := pflag.NewFlagSet(cmd, pflag.ContinueOnError)
	flags.StringVar(&cfg.GitCmd, "git-command", cfg.GitCmd, "git executable to use.")
	flags.StringVar(&cfg.BasePath, "base-path", cfg.BasePath, "Contains the path used as starting point for finding projects.")
	flags.StringVarP(&cfg.Branch, "branch", "b", cfg.Branch, "Name of branch to switch to.")
	if err := flags.Parse(args); err != nil {
		return nil, err
	}

	if cfg.BasePath == "" {
		return nil, errors.New("base-path can not be empty")
	}

	if cfg.Branch == "" {
		return nil, errors.New("target branch can not be empty")
	}

	return cfg, nil
}
