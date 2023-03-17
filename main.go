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

	results := git.SwitchBranches(ctx, log, cfg, switchConfigs)

	sort.Slice(results, func(i, j int) bool {
		return strings.Compare(results[i].Config.Name, results[j].Config.Name) < 0
	})

	errors := 0
	for _, r := range results {
		if r.Err != nil {
			fmt.Printf("%s: Error while switching: %s\n", r.Config.Name, r.Err)
			errors++

			continue
		}
	}

	os.Exit(errors)
}
