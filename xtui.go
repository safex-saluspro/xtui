package xtui

import (
	c "github.com/faelmori/xtui/components"
	t "github.com/faelmori/xtui/types"
	x "github.com/faelmori/xtui/wrappers"
)

func NewAppDepsModel(apps []string, path string, yes bool, quiet bool) x.AppDepsModel {
	return x.NewAppDepsModel(apps, path, yes, quiet)
}

func LogViewer(args ...string) error {
	return x.LogViewer(args...)
}

func ShowForm(form t.Config) (map[string]string, error) {
	return c.ShowForm(form)
}
