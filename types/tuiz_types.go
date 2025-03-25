package types

type TableDataHandler interface {
	GetHeaders() []string
	GetRows() [][]string
}
