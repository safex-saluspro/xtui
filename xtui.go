package xtui

import (
	c "github.com/faelmori/xtui/components"
	t "github.com/faelmori/xtui/types"
)

type Config struct{ t.Config }
type FormFields = t.FormFields
type FormField = t.FormField
type InputField = *t.InputField

func LogViewer(args ...string) error {
	return t.LogViewer(args...)
}
func ShowForm(form Config) (map[string]string, error) {
	return c.ShowForm(form.Config)
}

func NewConfig(title string, fields FormFields) Config {
	return Config{Config: t.Config{Title: title, Fields: fields}}
}
func NewInputField(placeholder string, typ string, value string, required bool, minValue int, maxValue int, err string, validation func(string) error) *FormField {
	return &InputField{
		t.InputField{
			Ph:  placeholder,
			Tp:  typ,
			Val: value,
			Req: required,
			Min: minValue,
			Max: maxValue,
			Err: err,
			Vld: validation,
		},
	}
}
func NewFormFields(title string, fields []FormField) FormFields {
	ffs := make([]t.FormField, len(fields))
	for i, f := range fields {
		ffs[i] = f.FormField
	}
	return FormFields{
		t.FormFields{
			Title:  title,
			Fields: ffs,
		},
	}
}
func NewFormModel(config t.Config) (map[string]string, error) { return c.ShowForm(config) }
