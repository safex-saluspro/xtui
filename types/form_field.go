package types

type FormFields struct {
	Title  string
	Fields []FormInputObject[any]
}

func (f FormFields) InputType() string {
	return f.Title
}
func (f FormFields) Inputs() []FormInputObject[any] {
	return f.Fields
}
