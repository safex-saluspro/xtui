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

// New interfaces and structs for customization options, validation, and layout

type CustomizableField interface {
	FormField
	Label() string
	DefaultValue() string
	Group() string
	LayoutOptions() LayoutOptions
	Styles() Styles
}

type LayoutOptions struct {
	Horizontal bool
	Vertical   bool
}

type Styles struct {
	FieldStyle    string
	LabelStyle    string
	ErrorStyle    string
	TemplateStyle string
}

type CustomField struct {
	InputField
	Lbl  string
	DVal string
	Grp  string
	Lay  LayoutOptions
	Sty  Styles
}

func (f CustomField) Label() string                { return f.Lbl }
func (f CustomField) DefaultValue() string         { return f.DVal }
func (f CustomField) Group() string                { return f.Grp }
func (f CustomField) LayoutOptions() LayoutOptions { return f.Lay }
func (f CustomField) Styles() Styles               { return f.Sty }

type Exporter interface {
	ExportToCSV(filename string) error
	ExportToYAML(filename string) error
	ExportToJSON(filename string) error
	ExportToXML(filename string) error
	ExportToExcel(filename string) error
	ExportToPDF(filename string) error
	ExportToMarkdown(filename string) error
}

type DataExporter struct{}

func (e DataExporter) ExportToCSV(filename string) error {
	// Implementation for exporting to CSV
	return nil
}

func (e DataExporter) ExportToYAML(filename string) error {
	// Implementation for exporting to YAML
	return nil
}

func (e DataExporter) ExportToJSON(filename string) error {
	// Implementation for exporting to JSON
	return nil
}

func (e DataExporter) ExportToXML(filename string) error {
	// Implementation for exporting to XML
	return nil
}

func (e DataExporter) ExportToExcel(filename string) error {
	// Implementation for exporting to Excel
	return nil
}

func (e DataExporter) ExportToPDF(filename string) error {
	// Implementation for exporting to PDF
	return nil
}

func (e DataExporter) ExportToMarkdown(filename string) error {
	// Implementation for exporting to Markdown
	return nil
}
