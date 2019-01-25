package project

import (
	"io"
	"path/filepath"

	"github.com/bilus/oya/pkg/changeset"
	"github.com/bilus/oya/pkg/oyafile"
	"github.com/bilus/oya/pkg/pack"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// TODO: Duplicated in oyafile module.
const VendorDir = ".oya/vendor"

type Project struct {
	Root *oyafile.Oyafile
}

func Load(rootDir string) (Project, error) {
	root, ok, err := oyafile.LoadFromDir(rootDir, rootDir)
	if !ok {
		err = errors.Errorf("Missing Oyafile at %v", rootDir)
	}
	return Project{
		Root: root,
	}, err
}

func Detect(workDir string) (Project, error) {
	o, found, err := detectRoot(workDir)
	if err != nil {
		return Project{}, err
	}
	if !found {
		return Project{}, ErrNoProject{Path: workDir}
	}
	return Project{
		Root: o,
	}, nil
}

func (p Project) Run(workDir, taskName string, stdout, stderr io.Writer) error {
	log.Debugf("Task %q at %v", taskName, workDir)

	changes, err := p.changeset(workDir)
	if err != nil {
		return err
	}

	if len(changes) == 0 {
		return nil
	}

	foundAtLeastOneTask := false
	for _, o := range changes {
		found, err := o.RunTask(taskName, stdout, stderr)
		if err != nil {
			return errors.Wrapf(err, "error in %v", o.Path)
		}
		if found {
			foundAtLeastOneTask = found
		}
	}

	if !foundAtLeastOneTask {
		return ErrNoTask{
			Task: taskName,
		}
	}
	return nil
}

// Tasks returns tasks tables by its Oyafile path (relative to project root) for each Oyafile in the changeset.
// It returns only tasks for the current working directory and its subdirectories.
func (p Project) Tasks(workDir string, stdout, stderr io.Writer) (map[string]oyafile.TaskTable, error) {
	changes, err := p.changeset(workDir)
	if err != nil {
		return nil, err
	}

	tasksByDir := make(map[string]oyafile.TaskTable)
	for _, o := range changes {
		tasksByDir[o.RelPath()] = o.Tasks
	}

	return tasksByDir, nil
}

func (p Project) changeset(workDir string) ([]*oyafile.Oyafile, error) {
	oyafiles, err := listOyafiles(workDir)
	if err != nil {
		return nil, err
	}
	for _, o := range oyafiles {
		log.Println(o.Path)
	}
	if len(oyafiles) == 0 {
		return nil, ErrNoOyafiles{Path: workDir}
	}

	return changeset.Calculate(oyafiles)
}

func (p Project) Oyafile(oyafilePath string) (*oyafile.Oyafile, bool, error) {
	return oyafile.Load(oyafilePath, p.Root.RootDir)
}

func (p Project) Vendor(pack pack.Pack) error {
	return pack.Vendor(filepath.Join(p.Root.RootDir, VendorDir))
}

func isRoot(o *oyafile.Oyafile) bool {
	return len(o.Project) > 0
}

func detectRoot(startDir string) (*oyafile.Oyafile, bool, error) {
	path := startDir
	maxParts := 256
	for i := 0; i < maxParts; i++ {
		o, found, err := oyafile.LoadFromDir(path, path)
		if err != nil {
			return nil, false, err
		}
		if err == nil && found && isRoot(o) {
			return o, true, nil
		}
		if path == "/" {
			break
		}
		path = filepath.Dir(path)
	}

	return nil, false, nil
}
