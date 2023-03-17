package git

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/sirupsen/logrus"
	"github.com/xperimental/git-multiswitch/internal/config"
	"github.com/xperimental/git-multiswitch/internal/logger"
)

type Result struct {
	Config SwitchConfig
	Output string
	Err    error
}

func SwitchBranches(ctx context.Context, log logger.Logger, cfg *config.Config, switchConfigs []SwitchConfig) []Result {
	results := []Result{}

	for _, repo := range switchConfigs {
		output, err := switchRepository(ctx, log, cfg, repo)

		results = append(results, Result{
			Config: repo,
			Output: output,
			Err:    err,
		})
	}

	return results
}

func switchRepository(ctx context.Context, log logger.Logger, cfg *config.Config, repo SwitchConfig) (string, error) {
	cmd := exec.CommandContext(ctx, cfg.GitCmd, "checkout", repo.TargetBranch)
	cmd.Dir = repo.Path

	if cfg.ShowOutput {
		log.Infof("Switching %q to %q...", repo.Name, repo.TargetBranch)
		cmd.Stdout = log.WriterLevel(logrus.InfoLevel)
		cmd.Stderr = log.WriterLevel(logrus.WarnLevel)
		if err := cmd.Run(); err != nil {
			return "", fmt.Errorf("can not run %q: %w", cfg.GitCmd, err)
		}

		return "", nil
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return string(output), fmt.Errorf("can not run %q: %w", cfg.GitCmd, err)
	}

	return string(output), nil
}
