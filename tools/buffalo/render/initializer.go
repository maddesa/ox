package render

import (
	"bytes"
	"context"
	"errors"
	"html/template"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/spf13/pflag"
)

var (
	ErrIncompleteArgs = errors.New("incomplete args")
)

// Initializer
type Initializer struct{}

func (i Initializer) Name() string {
	return "model/initializer"
}

func (i *Initializer) Initialize(ctx context.Context) error {
	m := ctx.Value("module")
	if m == nil {
		return ErrIncompleteArgs
	}

	f := ctx.Value("folder")
	if f == nil {
		return ErrIncompleteArgs
	}

	folder := filepath.Join(f.(string), "app", "render")
	err := os.MkdirAll(folder, 0777)
	if err != nil {
		return err
	}

	tmpl, err := template.New("render.go").Parse(renderGo)
	if err != nil {
		return err
	}

	var data = struct {
		Module string
	}{
		Module: m.(string),
	}

	sbf := bytes.NewBuffer([]byte{})
	err = tmpl.Execute(sbf, data)
	if err != nil {
		return err
	}

	path := filepath.Join(folder, "render.go")
	err = ioutil.WriteFile(path, sbf.Bytes(), 0777)
	if err != nil {
		return err
	}

	return nil
}

func (i *Initializer) ParseFlags([]string) {}
func (i *Initializer) Flags() *pflag.FlagSet {
	return pflag.NewFlagSet("buffalo-models-initializer", pflag.ContinueOnError)
}
