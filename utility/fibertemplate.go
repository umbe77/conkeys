package utility

import (
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"sync"
)

type Engine struct {
	Dir      string
	Template *template.Template
	mutex    sync.RWMutex
}

func NewEngine(dir string) *Engine {
	tmpl := template.New(dir)
	return &Engine{
		Dir:      dir,
		Template: tmpl,
	}
}

func (e *Engine) Load() error {
	e.mutex.Lock()
	defer e.mutex.Unlock()
	return filepath.Walk(e.Dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info == nil || info.IsDir() {
			return nil
		}

		name, err := filepath.Rel(e.Dir, path)
		if err != nil {
			return err
		}
		name = name[:len(name)-len(filepath.Ext(name))]

		buf, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		_, err = e.Template.New(name).Parse(string(buf))
		if err != nil {
			return err
		}

		return err
	})
}

func (e *Engine) Render(out io.Writer, template string, binding interface{}, lay ...string) error {
	tmpl := e.Template.Lookup(template)
	if tmpl == nil {
		return fmt.Errorf("Cannot find template: '%s'", template)
	}
	return tmpl.Execute(out, binding)
}
