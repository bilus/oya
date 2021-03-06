package oyafile

import (
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/bilus/oya/pkg/raw"
	"github.com/bilus/oya/pkg/semver"
	"github.com/bilus/oya/pkg/task"
	"github.com/bilus/oya/pkg/template"
	"github.com/bilus/oya/pkg/types"
	"github.com/pkg/errors"
)

// OyaCmdOverride is used in tests, to override the path to the current oya executable.
// It is used to invoke other tasks from a task body.
// When tests are run, the current process executable path points to the test runner
// so it has to be overridden (with 'go run oya.go', roughly speaking).
var OyaCmdOverride *string

type PackReference struct {
	ImportPath types.ImportPath
	Version    semver.Version
	// ReplacementPath is a path relative to the root directory, when the replacement for the pack can be found, based on the Replace: directive.
	ReplacementPath string
}

type PackReplacements map[types.ImportPath]string

type Oyafile struct {
	Dir      string
	Path     string
	RootDir  string
	Shell    string
	Imports  map[types.Alias]types.ImportPath
	Tasks    task.Table
	Values   template.Scope
	Project  string   // Project is set for root Oyafile.
	Ignore   []string // Ignore contains directory exclusion rules.
	Requires []PackReference
	// Replacements map packs to local paths relative to project root directory for development based on the Replace: directive.
	Replacements PackReplacements
	IsBuilt      bool

	relPath string

	OyaCmd string // OyaCmd contains the path to the current oya executable.
}

func New(oyafilePath string, rootDir string) (*Oyafile, error) {
	var oyaCmd string
	if OyaCmdOverride != nil {
		oyaCmd = *OyaCmdOverride
	} else {
		var err error
		oyaCmd, err = os.Executable()
		if err != nil {
			return nil, err
		}
	}

	relPath, err := filepath.Rel(rootDir, oyafilePath)
	log.Debug("Oyafile at ", oyafilePath)
	if err != nil {
		return nil, err
	}
	dir := path.Dir(oyafilePath)
	return &Oyafile{
		Dir:          filepath.Clean(dir),
		Path:         filepath.Clean(oyafilePath),
		RootDir:      filepath.Clean(rootDir),
		Shell:        "/bin/bash",
		Imports:      make(map[types.Alias]types.ImportPath),
		Tasks:        task.NewTable(),
		Values:       template.Scope{},
		relPath:      relPath,
		Replacements: make(PackReplacements),
		OyaCmd:       oyaCmd,
	}, nil
}

func Load(oyafilePath, rootDir string) (*Oyafile, bool, error) {
	raw, found, err := raw.Load(oyafilePath, rootDir)
	if err != nil || !found {
		return nil, found, err
	}
	oyafile, err := Parse(raw)
	if err != nil {
		return nil, false, wrapLoadErr(err, oyafilePath)
	}
	return oyafile, true, nil
}

func LoadFromDir(dirPath, rootDir string) (*Oyafile, bool, error) {
	raw, found, err := raw.LoadFromDir(dirPath, rootDir)
	if err != nil || !found {
		return nil, found, err
	}
	oyafile, err := Parse(raw)
	if err != nil {
		return nil, false, wrapLoadErr(err, raw.Path)
	}
	return oyafile, true, nil
}

func (oyafile Oyafile) RunTask(taskName task.Name, scope template.Scope, stdout, stderr io.Writer) (bool, error) {
	if !oyafile.IsBuilt {
		return false, errors.Errorf("Internal error: Oyafile has not been built")
	}
	task, ok := oyafile.Tasks.LookupTask(taskName)
	if !ok {
		return false, nil
	}
	tasks, err := oyafile.bindTasks(taskName, task, stdout, stderr)
	if err != nil {
		return true, err
	}
	scope["Tasks"] = tasks

	render, err := oyafile.bindRender(taskName, stdout, stderr)
	if err != nil {
		return true, err
	}
	scope["Render"] = render

	return true, task.Exec(oyafile.Dir, scope, stdout, stderr)
}

func (oyafile Oyafile) Equals(other Oyafile) bool {
	// TODO: Far from perfect, we should ensure relative vs absolute paths work.
	// The simplest thing is probably to ensure oyafile.Dir is always absolute.
	return filepath.Clean(oyafile.Dir) == filepath.Clean(other.Dir)
}

func wrapLoadErr(err error, oyafilePath string) error {
	return errors.Wrapf(err, "error loading Oyafile %v", oyafilePath)
}

func (o *Oyafile) Ignores() string {
	return strings.Join(o.Ignore, "\n")
}

func (o *Oyafile) RelPath() string {
	return o.relPath
}
