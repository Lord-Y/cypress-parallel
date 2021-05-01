// Package hooks will manage all hooks requirements
package hooks

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/Lord-Y/cypress-parallel-api/tools"
	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/icrowley/fake"
	"github.com/rs/zerolog/log"
)

func (p *projects) plainClone(branch string, specs string) (z []string, statusCode int, err error) {
	var (
		targetSpecs  string
		targetBranch string
	)
	if specs != "" {
		targetSpecs = specs
	} else {
		targetSpecs = p.Specs
	}

	if branch != "" {
		if branch == "master" {
			targetBranch = ""
		} else {
			targetBranch = branch
		}
	} else {
		if p.Branch == "master" {
			targetBranch = ""
		} else {
			targetBranch = p.Branch
		}
	}
	log.Debug().Msgf("branch %s specs %s", targetBranch, targetSpecs)

	directory, err := os.MkdirTemp(os.TempDir(), fake.CharactersN(10))
	if err != nil {
		return z, 500, err
	}
	defer os.RemoveAll(directory)

	if targetBranch != "" {
		result, err := git.PlainClone(directory, false, &git.CloneOptions{
			URL: p.Repository,
		})
		if err != nil {
			log.Debug().Msgf("Cloning repo error %s", err.Error())
			return z, 400, err
		}
		w, err := result.Worktree()
		if err != nil {
			log.Debug().Msgf("Tree error %s", err.Error())
			return z, 500, err
		}
		err = w.Checkout(&git.CheckoutOptions{
			Branch: plumbing.ReferenceName(targetBranch),
		})
		if err != nil {
			log.Debug().Msgf("Checkout %s", err.Error())
			return z, 400, err
		}
	} else {
		_, err = git.PlainClone(directory, false, &git.CloneOptions{
			URL: p.Repository,
		})
		if err != nil {
			log.Debug().Msgf("cloning repo with no new branch setted %s", err.Error())
			return z, 400, err
		}
	}

	if strings.HasSuffix(targetSpecs, ".spec.js") || strings.HasSuffix(targetSpecs, ".ts") {
		err := tools.CheckIsFile(fmt.Sprintf("%s/%s", directory, targetSpecs))
		if err != nil {
			return z, 400, err
		}
		z = append(z, targetSpecs)
	} else {
		err := filepath.Walk(directory, func(file string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.Mode().IsRegular() && (strings.HasSuffix(file, ".spec.js") || strings.HasSuffix(file, ".ts")) {
				z = append(z, strings.ReplaceAll(file, fmt.Sprintf("%s/", directory), ""))
			}
			return nil
		})
		if err != nil {
			return z, 400, err
		}
	}
	return
}
