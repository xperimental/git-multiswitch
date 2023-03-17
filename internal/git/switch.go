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
	Err    error
	Config SwitchConfig
}

func SwitchBranches(ctx context.Context, log logger.Logger, cfg *config.Config, switchConfigs []SwitchConfig) []Result {
	results := []Result{}

	for _, repo := range switchConfigs {
		log.Infof("Switching %q to %q...", repo.Name, repo.TargetBranch)
		err := switchRepository(ctx, log, cfg.GitCmd, repo.TargetBranch, repo.Path)

		results = append(results, Result{
			Err:    err,
			Config: repo,
		})
	}

	return results
}

func switchRepository(ctx context.Context, log logger.Logger, gitCmd, branch, path string) error {
	cmd := exec.CommandContext(ctx, gitCmd, "checkout", branch)
	cmd.Stdout = log.WriterLevel(logrus.DebugLevel)
	cmd.Stderr = log.WriterLevel(logrus.InfoLevel)
	cmd.Dir = path

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("can not run %q: %w", gitCmd, err)
	}

	return nil
}
