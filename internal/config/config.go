package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/pflag"
)

const (
	envPrefix = "GIT_MULTISWITCH_"

	envLogLevel   = envPrefix + "LOG_LEVEL"
	envGitCmd     = envPrefix + "GIT_COMMAND"
	envBasePath   = envPrefix + "BASE_PATH"
	envBranch     = envPrefix + "BRANCH"
	envDryRun     = envPrefix + "DRY_RUN"
	envShowOutput = envPrefix + "SHOW_GIT_OUTPUT"
	envEscapeRepo = envPrefix + "ESCAPE_REPO"
)

type Config struct {
	LogLevel   string
	GitCmd     string
	BasePath   string
	Branch     string
	DryRun     bool
	ShowOutput bool
	EscapeRepo bool
}

func (c Config) LogrusLevel() logrus.Level {
	lvl, _ := logrus.ParseLevel(c.LogLevel)
	return lvl
}

func Parse(cmd string, args []string, envFunc func(string) string) (*Config, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("can not get working directory: %w", err)
	}

	cfg := &Config{
		LogLevel: "info",
		GitCmd:   "git",
		BasePath: wd,
	}

	flags := pflag.NewFlagSet(cmd, pflag.ContinueOnError)
	flags.StringVar(&cfg.LogLevel, "log-level", cfg.LogLevel, "Minimum log level to display.")
	flags.StringVar(&cfg.GitCmd, "git-command", cfg.GitCmd, "git executable to use.")
	flags.StringVar(&cfg.BasePath, "base-path", cfg.BasePath, "Contains the path used as starting point for finding projects.")
	flags.StringVarP(&cfg.Branch, "branch", "b", cfg.Branch, "Name of branch to switch to.")
	flags.BoolVarP(&cfg.DryRun, "dry-run", "n", cfg.DryRun, "If true, only show what would be done without actually switching branches.")
	flags.BoolVar(&cfg.ShowOutput, "show-git-output", cfg.ShowOutput, "If true, shows git output while switching branches.")
	flags.BoolVar(&cfg.EscapeRepo, "escape-repo", cfg.EscapeRepo, "Escape to the parent repository if run inside a git repository.")
	if err := flags.Parse(args); err != nil {
		return nil, err
	}

	if logLevel := envFunc(envLogLevel); logLevel != "" {
		cfg.LogLevel = logLevel
	}

	if gitCmd := envFunc(envGitCmd); gitCmd != "" {
		cfg.GitCmd = gitCmd
	}

	if basePath := envFunc(envBasePath); basePath != "" {
		cfg.BasePath = basePath
	}

	if branch := envFunc(envBranch); branch != "" {
		cfg.Branch = branch
	}

	if dryRun := envFunc(envDryRun); dryRun != "" {
		cfg.DryRun = true
	}

	if showOutput := envFunc(envShowOutput); showOutput != "" {
		cfg.ShowOutput = true
	}

	if escapeRepo := envFunc(envEscapeRepo); escapeRepo != "" {
		cfg.EscapeRepo = true
	}

	if _, err := logrus.ParseLevel(cfg.LogLevel); err != nil {
		return nil, fmt.Errorf("can not parse log-level %q: %w", cfg.LogLevel, err)
	}

	if cfg.GitCmd == "" {
		return nil, errors.New("git-command can not be empty")
	}

	if cfg.BasePath == "" {
		return nil, errors.New("base-path can not be empty")
	}

	if cfg.Branch == "" {
		return nil, errors.New("target branch can not be empty")
	}

	return cfg, nil
}
