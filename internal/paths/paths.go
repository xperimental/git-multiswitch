package paths

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/xperimental/git-multiswitch/internal/config"
)

func ListRepositories(ctx context.Context, cfg *config.Config) ([]string, error) {
	basePath := cfg.BasePath
	if cfg.EscapeRepo {
		// If run from inside a git repository use the parent directory of the repository as a
		// base path instead.
		repoPath, insideRepo, err := isPathInsideGitRepository(ctx, cfg, basePath)
		if err != nil {
			return nil, fmt.Errorf("can not determine if inside a git repository: %w", err)
		}

		if insideRepo {
			basePath = filepath.Dir(repoPath)
		}
	}

	var paths []string
	err := filepath.WalkDir(basePath, func(path string, d fs.DirEntry, err error) error {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		if !d.IsDir() {
			return nil
		}

		hasGit, err := hasGitDir(path)
		if err != nil {
			return err
		}

		if hasGit {
			paths = append(paths, path)
		}

		if filepath.Base(path) == ".git" {
			return filepath.SkipDir
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return paths, nil
}

func hasGitDir(path string) (bool, error) {
	gitDir := filepath.Join(path, ".git")
	info, err := os.Stat(gitDir)
	switch {
	case os.IsNotExist(err):
		return false, nil
	case err != nil:
		return false, err
	default:
		return info.IsDir(), nil
	}
}

func isPathInsideGitRepository(ctx context.Context, cfg *config.Config, path string) (string, bool, error) {
	hasGit, err := hasGitDir(path)
	if err != nil {
		return "", false, err
	}

	if hasGit {
		return path, true, nil
	}

	parent := filepath.Dir(path)
	if parent == path {
		return "", false, nil
	}
	
	return isPathInsideGitRepository(ctx, cfg, parent)
}
