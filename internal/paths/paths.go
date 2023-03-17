package paths

import (
	"context"
	"io/fs"
	"os"
	"path/filepath"
)

func ListRepositories(ctx context.Context, basePath string) ([]string, error) {
	var paths []string
	err := filepath.WalkDir(basePath, func(path string, d fs.DirEntry, err error) error {
		if ctx.Err() != nil {
			return ctx.Err()
		}

		if !d.IsDir() {
			return nil
		}

		gitDir := filepath.Join(path, ".git")
		info, err := os.Stat(gitDir)
		switch {
		case os.IsNotExist(err):
			return nil
		case err != nil:
			return err
		default:
		}

		if info.IsDir() {
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
