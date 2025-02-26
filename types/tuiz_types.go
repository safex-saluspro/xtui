package types

const DOWN = "down"
const TAB = "tab"
const SHIFTTAB = "shift+tab"
const ENTER = "enter"
const UP = "up"
const CTRLR = "ctrl+r"
const CTRLC = "ctrl+c"
const ESC = "esc"
const PASSWORD = "password"

type TableDataHandler interface {
	GetHeaders() []string
	GetRows() [][]string
}

type Field interface {
	Type() string
	Value() string
	Placeholder() string
	ValidationRules() []ValidationRule
}

type ValidationRule interface {
	Validate(value string) error
}

type FormConfig struct {
	Title  string
	Fields []Field
}

type TuizFieldz interface {
	InputType() string
	Inputs() []interface{}
}
type TuizInputz interface {
	Placeholder() string
	Type() string
	Value() string
	Required() bool
	MinLength() int
	MaxLength() int
	ErrorMessage() string
	Validator() func(string) error
}
type TuizConfig interface {
	Title() string
	Fields() interface{}
}

type TuizInput struct {
	Ph  string
	Tp  string
	Val string
	Req bool
	Min int
	Max int
	Err string
	Vld func(string) error
}

func (f TuizInput) Placeholder() string           { return f.Ph }
func (f TuizInput) Type() string                  { return f.Tp }
func (f TuizInput) Value() string                 { return f.Val }
func (f TuizInput) Required() bool                { return f.Req }
func (f TuizInput) MinLength() int                { return f.Min }
func (f TuizInput) MaxLength() int                { return f.Max }
func (f TuizInput) ErrorMessage() string          { return f.Err }
func (f TuizInput) Validator() func(string) error { return f.Vld }

type TuizFields struct {
	Tt  string
	Fds []TuizInputz
}

func (f TuizFields) InputType() string {
	return f.Tt
}
func (f TuizFields) Inputs() []interface{} {
	var inputs []interface{}
	for _, field := range f.Fds {
		inputs = append(inputs, field)
	}
	return inputs
}

type TuizConfigz struct {
	Tt  string
	Fds TuizFieldz
}

func (c TuizConfigz) Title() string { return c.Tt }
func (c TuizConfigz) Fields() interface{} {
	return c.Fds
}
