// Package git will manage all requirements to clone repository
package git

import (
	"os"

	git "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/icrowley/fake"
	"github.com/rs/zerolog/log"
)

type Repository struct {
	Repository string // HTTP(s) git repository
	Username   string // Username to use to fetch repository if required
	Password   string // Password to use to fetch repository if required
	Ref        string // Ref in which branch e.g test or refs/head/test
}

// Clone permit to clone git repository
func (c *Repository) Clone() (z string, statusCode int, err error) {
	var (
		targetRef string
		result    *git.Repository
	)

	if c.Ref != "" {
		if c.Ref == "master" || c.Ref == "refs/heads/master" {
			targetRef = ""
		} else {
			targetRef = c.Ref
		}
	} else {
		targetRef = ""
	}
	log.Debug().Msgf("Branch or tag %s", targetRef)

	z, err = os.MkdirTemp(os.TempDir(), fake.CharactersN(10))
	if err != nil {
		return z, 500, err
	}

	if targetRef != "" {
		if c.Username != "" {
			result, err = git.PlainClone(z, false, &git.CloneOptions{
				URL: c.Repository,
				Auth: &http.BasicAuth{
					Username: c.Username,
					Password: c.Password,
				},
				ReferenceName: plumbing.ReferenceName(targetRef),
				SingleBranch:  true,
				Depth:         1,
			})
		} else {
			result, err = git.PlainClone(z, false, &git.CloneOptions{
				URL: c.Repository,
			})
		}
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
			Branch: plumbing.ReferenceName(targetRef),
		})
		if err != nil {
			log.Debug().Msgf("Checkout %s", err.Error())
			return z, 400, err
		}
	} else {
		if c.Username != "" {
			_, err = git.PlainClone(z, false, &git.CloneOptions{
				URL: c.Repository,
				Auth: &http.BasicAuth{
					Username: c.Username,
					Password: c.Password,
				},
				ReferenceName: plumbing.ReferenceName(targetRef),
				SingleBranch:  true,
				Depth:         1,
			})
		} else {
			_, err = git.PlainClone(z, false, &git.CloneOptions{
				URL:           c.Repository,
				ReferenceName: plumbing.ReferenceName(targetRef),
				SingleBranch:  true,
				Depth:         1,
			})
		}
		if err != nil {
			return z, 400, err
		}
	}
	return
}
