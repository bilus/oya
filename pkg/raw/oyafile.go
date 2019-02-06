package raw

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"os"
	"path"

	yaml "gopkg.in/yaml.v2"
)

const DefaultName = "Oyafile"

type Oyafile struct {
	Path    string
	RootDir string
	file    []byte
}

// DecodedOyafile is an Oyafile that has been loaded from YAML
// but hasn't been parsed yet.
type DecodedOyafile map[string]interface{}

func Load(oyafilePath, rootDir string) (*Oyafile, bool, error) {
	raw, err := New(oyafilePath, rootDir)
	if err != nil {
		return nil, false, nil
	}
	return raw, true, nil
}

func LoadFromDir(dirPath, rootDir string) (*Oyafile, bool, error) {
	oyafilePath := fullPath(dirPath, "")
	fi, err := os.Stat(oyafilePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, false, nil
		}
		return nil, false, err
	}
	if fi.IsDir() {
		return nil, false, nil
	}
	return Load(oyafilePath, rootDir)
}

func New(oyafilePath, rootDir string) (*Oyafile, error) {
	file, err := ioutil.ReadFile(oyafilePath)
	if err != nil {
		return nil, err
	}

	return &Oyafile{
		RootDir: rootDir,
		Path:    oyafilePath,
		file:    file,
	}, nil
}

func (raw *Oyafile) Decode() (DecodedOyafile, error) {
	// YAML parser does not handle files without at least one node.
	empty, err := isEmptyYAML(raw.Path)
	if err != nil {
		return nil, err
	}
	if empty {
		return make(DecodedOyafile), nil
	}
	reader := bytes.NewReader(raw.file)
	decoder := yaml.NewDecoder(reader)
	var of DecodedOyafile
	err = decoder.Decode(&of)
	if err != nil {
		return nil, err
	}
	return of, nil
}

func (raw *Oyafile) HasKey(key string) (bool, error) {
	of, err := raw.Decode()
	if err != nil {
		return false, err
	}
	_, ok := of[key]
	return ok, nil
}

// isEmptyYAML returns true if the Oyafile contains only blank characters or YAML comments.
func isEmptyYAML(oyafilePath string) (bool, error) {
	file, err := os.Open(oyafilePath)
	if err != nil {
		return false, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if isNode(scanner.Text()) {
			return false, nil
		}
	}

	return true, scanner.Err()
}

func isNode(line string) bool {
	for _, c := range line {
		switch c {
		case '#':
			return false
		case ' ', '\t', '\n', '\f', '\r':
			continue
		default:
			return true
		}
	}
	return false
}

func fullPath(projectDir, name string) string {
	if len(name) == 0 {
		name = DefaultName
	}
	return path.Join(projectDir, name)
}
