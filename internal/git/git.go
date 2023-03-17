package git

import (
	"context"
	"fmt"
	"os/exec"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/xperimental/git-multiswitch/internal/logger"
)

type Result struct {
	Err    error
	Path   string
	Name   string
	Branch string
}

func SwitchBranches(ctx context.Context, log logger.Logger, gitCmd, targetBranch string, repositoryPaths []string) []Result {
	results := []Result{}

	for _, path := range repositoryPaths {
		name := filepath.Base(path)
		log.Infof("Switching %q to %q...", name, targetBranch)
		err := switchRepository(ctx, log, gitCmd, targetBranch, path)

		results = append(results, Result{
			Err:    err,
			Path:   path,
			Name:   name,
			Branch: targetBranch,
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
