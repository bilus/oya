package oyafile

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

const VendorDir = ".oya/vendor"

func (oyafile *Oyafile) resolveImports() error {
	for alias, path := range oyafile.Imports {
		log.Debugf("Importing pack %v as %v", path, alias)
		pack, err := oyafile.loadPack(path)
		if err != nil {
			return err
		}
		oyafile.Values[string(alias)] = pack.Values
		for key, task := range pack.Tasks {
			// TODO: Detect if task already set.
			log.Printf("Importing task %v.%v", alias, key)
			oyafile.Tasks[fmt.Sprintf("%v.%v", alias, key)] = task
		}
	}
	return nil
}

func (oyafile *Oyafile) loadPack(path ImportPath) (*Oyafile, error) {
	for _, importDir := range oyafile.importDirs() {
		fullPath := filepath.Join(importDir, string(path))
		if !isValidImportPath(fullPath) {
			continue
		}
		pack, found, err := LoadFromDir(fullPath, oyafile.RootDir)
		if err != nil {
			continue
		}
		if !found {
			continue
		}
		return pack, nil
	}

	return nil, errors.Errorf("missing pack %v", path)
}

func (oyafile *Oyafile) importDirs() []string {
	return []string{
		oyafile.RootDir,
		filepath.Join(oyafile.RootDir, VendorDir),
	}
}

func isValidImportPath(fullImportPath string) bool {
	f, err := os.Stat(fullImportPath)
	return err == nil && f.IsDir()
}

func AddImport(dirPath string, uri string) error {
	oyafilePath := filepath.Join(dirPath, DefaultName)
	info, _ := os.Stat(oyafilePath)
	file, err := ioutil.ReadFile(oyafilePath)
	if err != nil {
		return err
	}

	importStr := "Import:"
	uriStr := fmt.Sprintf("  oya: %s", uri)
	fileContent := string(file)
	fileArr := strings.Split(fileContent, "\n")
	var arr []string
	if strings.Contains(fileContent, "Import:") {
		for _, line := range fileArr {
			arr = append(arr, line)
			if strings.Contains(line, "Import") {
				arr = append(arr, uriStr)
			}
		}
	} else if strings.Contains(fileContent, "Project:") {
		for _, line := range fileArr {
			arr = append(arr, line)
			if strings.Contains(line, "Project") {
				arr = append(arr, importStr)
				arr = append(arr, uriStr)
			}
		}
	} else {
		arr = append(arr, importStr)
		arr = append(arr, uriStr)
		arr = append(fileArr, arr...)
	}

	ioutil.WriteFile(oyafilePath, []byte(strings.Join(arr, "\n")), info.Mode())
	return nil
}
