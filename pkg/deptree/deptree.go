package deptree

import (
	"github.com/bilus/oya/pkg/deptree/internal"
	"github.com/bilus/oya/pkg/mvs"
	"github.com/bilus/oya/pkg/oyafile"
	"github.com/bilus/oya/pkg/pack"
	"github.com/bilus/oya/pkg/types"
)

// DependencyTree defines a project's dependencies, allowing for loading them.
type DependencyTree struct {
	installDirs  []string
	dependencies []pack.Pack
	reqs         *internal.Reqs
}

// New returns a new dependency tree.
// BUG(bilus): It's called a 'tree' but it currently does not take into account inter-pack
// dependencies. This will likely change and then the name will fit like a glove. ;)
func New(rootDir string, installDirs []string, dependencies []pack.Pack) (*DependencyTree, error) {
	return &DependencyTree{
		installDirs:  installDirs,
		dependencies: dependencies,
		reqs:         internal.NewReqs(rootDir, installDirs),
	}, nil
}

// Explode takes the initial list of dependencies and builds the full list,
// taking into account packs' dependencies and using Minimal Version Selection.
func (dt *DependencyTree) Explode() error {
	list, err := mvs.List(dt.dependencies, dt.reqs)
	if err != nil {
		return err
	}
	dt.dependencies = list
	return nil
}

// Load loads an pack's Oyafile based on its import path.
// It supports two types of import paths:
// - referring to the project's Require: section (e.g. github.com/tooploox/oya-packs/docker), in this case it will load, the required version;
// - path relative to the project's root (e.g. /) -- does not support versioning, loads Oyafile directly from the path (<root dir>/<import path>).
func (dt *DependencyTree) Load(importPath types.ImportPath) (*oyafile.Oyafile, bool, error) {
	pack, found, err := dt.findRequiredPack(importPath)
	if err != nil {
		return nil, false, err
	}
	if found {
		return dt.reqs.LoadLocalOyafile(pack)
	}
	return nil, false, nil
}

// Find lookups pack by its import path.
func (dt *DependencyTree) Find(importPath types.ImportPath) (pack.Pack, bool, error) {
	for _, pack := range dt.dependencies {
		if pack.ImportPath() == importPath {
			return pack, true, nil
		}
	}
	return pack.Pack{}, false, nil
}

// ForEach iterates through the packs.
func (dt *DependencyTree) ForEach(f func(pack.Pack) error) error {
	for _, pack := range dt.dependencies {
		if err := f(pack); err != nil {
			return err
		}
	}
	return nil
}

func (dt *DependencyTree) findRequiredPack(importPath types.ImportPath) (pack.Pack, bool, error) {
	for _, pack := range dt.dependencies {
		if pack.ImportPath() == importPath {
			return pack, true, nil
		}
	}
	return pack.Pack{}, false, nil
}
