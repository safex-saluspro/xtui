package types

const (
	DOWN     = "down"
	TAB      = "tab"
	SHIFTTAB = "shift+tab"
	ENTER    = "enter"
	UP       = "up"
	CTRLR    = "ctrl+r"
	CTRLC    = "ctrl+c"
	ESC      = "esc"
	PASSWORD = "password"
)

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

type FormField interface {
	Placeholder() string
	Type() string
	Value() string
	Required() bool
	MinLength() int
	MaxLength() int
	ErrorMessage() string
	Validator() func(string) error
}

type InputField struct {
	Ph  string
	Tp  string
	Val string
	Req bool
	Min int
	Max int
	Err string
	Vld func(string) error
}

func (f InputField) Placeholder() string           { return f.Ph }
func (f InputField) Type() string                  { return f.Tp }
func (f InputField) Value() string                 { return f.Val }
func (f InputField) Required() bool                { return f.Req }
func (f InputField) MinLength() int                { return f.Min }
func (f InputField) MaxLength() int                { return f.Max }
func (f InputField) ErrorMessage() string          { return f.Err }
func (f InputField) Validator() func(string) error { return f.Vld }

// Collection of form fields
type FormFields struct {
	Title  string
	Fields []FormField
}

func (f FormFields) InputType() string {
	return f.Title
}
func (f FormFields) Inputs() []FormField {
	return f.Fields
}

type Config struct {
	Title  string
	Fields FormFields
}

func (c Config) GetTitle() string {
	return c.Title
}
func (c Config) GetFields() FormFields {
	return c.Fields
}
