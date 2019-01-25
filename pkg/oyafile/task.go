package oyafile

import (
	"io"
	"strings"

	"github.com/bilus/oya/pkg/template"
)

type Task interface {
	Exec(workDir string, stdout, stderr io.Writer) error
	GetName() string
	IsBuiltIn() bool
}

type ScriptedTask struct {
	Name string
	Script
	Shell string
	Scope *template.Scope
}

func (t ScriptedTask) Exec(workDir string, stdout, stderr io.Writer) error {
	return t.Script.Exec(workDir, *t.Scope, stdout, stderr, t.Shell)
}

func (t ScriptedTask) GetName() string {
	return t.Name
}

func (t ScriptedTask) IsBuiltIn() bool {
	firstChar := t.Name[0:1]
	return firstChar == strings.ToUpper(firstChar)
}

type BuiltinTask struct {
	Name   string
	OnExec func(stdout, stderr io.Writer) error
}

func (t BuiltinTask) Exec(workDir string, stdout, stderr io.Writer) error {
	return t.OnExec(stdout, stderr)
}

func (t BuiltinTask) GetName() string {
	return t.Name
}

func (t BuiltinTask) IsBuiltIn() bool {
	return true
}
