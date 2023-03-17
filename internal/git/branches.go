package git

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"
	"github.com/xperimental/git-multiswitch/internal/config"
	"github.com/xperimental/git-multiswitch/internal/logger"
)

type SwitchConfig struct {
	Path         string
	Name         string
	TargetBranch string
}

func FindBranches(ctx context.Context, log logger.Logger, cfg *config.Config, repositoryPaths []string) ([]SwitchConfig, []string, error) {
	var (
		result  []SwitchConfig
		skipped []string
	)

	for _, path := range repositoryPaths {
		hasBranch, err := hasTargetBranch(ctx, log, cfg, path)
		if err != nil {
			return nil, nil, fmt.Errorf("can not list branches of %q: %w", path, err)
		}

		if !hasBranch {
			skipped = append(skipped, path)
			continue
		}

		result = append(result, SwitchConfig{
			Path:         path,
			Name:         filepath.Base(path),
			TargetBranch: cfg.Branch,
		})
	}

	return result, skipped, nil
}

func hasTargetBranch(ctx context.Context, log logger.Logger, cfg *config.Config, repoPath string) (bool, error) {
	stdOut := &bytes.Buffer{}

	cmd := exec.CommandContext(ctx, cfg.GitCmd, "branch", "--list", cfg.Branch)
	cmd.Stdout = stdOut
	cmd.Stderr = log.WriterLevel(logrus.WarnLevel)
	cmd.Dir = repoPath

	if err := cmd.Run(); err != nil {
		return false, fmt.Errorf("can not run %q: %w", cfg.GitCmd, err)
	}

	return strings.Contains(stdOut.String(), cfg.Branch), nil
}
