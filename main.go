package main

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/spf13/pflag"
	"github.com/xperimental/git-multiswitch/internal/config"
	"github.com/xperimental/git-multiswitch/internal/git"
	"github.com/xperimental/git-multiswitch/internal/logger"
	"github.com/xperimental/git-multiswitch/internal/paths"
)

var (
	log = logger.New()
)

func main() {
	cfg, err := config.Parse(os.Args[0], os.Args[1:])
	switch {
	case err == pflag.ErrHelp:
		return
	case err != nil:
		log.Fatalf("Error in configuration: %s", err)
	default:
	}
	log.SetLevel(cfg.LogrusLevel())

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	repos, err := paths.ListRepositories(ctx, cfg.BasePath)
	if err != nil {
		log.Fatalf("Error listing repositories: %s", err)
	}
	log.Debugf("Found repositories: %s", repos)

	if len(repos) == 0 {
		log.Fatal("No repositories found!")
	}

	switchConfigs, skipped, err := git.FindBranches(ctx, log, cfg, repos)
	if err != nil {
		log.Fatalf("Error while finding branches: %s", err)
	}

	if len(skipped) > 0 {
		fmt.Println("Skipped repositories:")
		for _, s := range skipped {
			fmt.Printf(" - %s\n", s)
		}
	}

	if cfg.DryRun {
		fmt.Println("Would switch the following repositories:")
		for _, repo := range switchConfigs {
			fmt.Printf(" - %s (at %q) to %q\n", repo.Name, repo.Path, repo.TargetBranch)
		}

		return
	}

	results := git.SwitchBranches(ctx, log, cfg, switchConfigs)

	sort.Slice(results, func(i, j int) bool {
		return strings.Compare(results[i].Config.Name, results[j].Config.Name) < 0
	})

	errors := 0
	fmt.Println("Results:")
	for _, r := range results {
		if r.Err != nil {
			fmt.Printf(" - %s had error while switching: %s\n", r.Config.Name, r.Err)
			if !cfg.ShowOutput {
				fmt.Printf("   git output:\n%s\n", r.Output)
			}

			errors++
			continue
		}

		fmt.Printf(" - %s switched to %q\n", r.Config.Name, r.Config.TargetBranch)
	}

	os.Exit(errors)
}
