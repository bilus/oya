package pack

import (
	"github.com/bilus/oya/pkg/semver"
	"github.com/bilus/oya/pkg/types"
	"github.com/pkg/errors"
)

type Repo interface {
	Install(version semver.Version, installDir string) error
	IsInstalled(version semver.Version, installDir string) (bool, error)
	InstallPath(version semver.Version, installDir string) string
	ImportPath() types.ImportPath
}

// Pack represents a specific version of an Oya pack.
type Pack struct {
	repo    Repo
	version semver.Version
}

func New(repo Repo, version semver.Version) (Pack, error) {
	return Pack{
		repo:    repo,
		version: version,
	}, nil
}

func (p Pack) Install(installDir string) error {
	err := p.repo.Install(p.version, installDir)
	if err != nil {
		return errors.Wrapf(err, "error installing pack %v", p.ImportPath())
	}
	return nil
}

func (p Pack) IsInstalled(installDir string) (bool, error) {
	return p.repo.IsInstalled(p.version, installDir)
}

func (p Pack) Version() semver.Version {
	return p.version
}

func (p Pack) ImportPath() types.ImportPath {
	return p.repo.ImportPath()
}

func (p Pack) InstallPath(installDir string) string {
	return p.repo.InstallPath(p.version, installDir)
}
