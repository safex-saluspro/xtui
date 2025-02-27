package components

import (
	"encoding/csv"
	"fmt"
	"github.com/atotto/clipboard"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/faelmori/logz"
	. "github.com/faelmori/xtui/types"
	"os"
	"sort"
	"strconv"
	"strings"
)

type KuizTableModel struct {
	kTb          *table.Table
	headers      []string
	rows         [][]string
	filter       string
	filteredRows [][]string
	sortColumn   int
	sortAsc      bool
	page         int
	pageSize     int
	search       string
	selectedRow  int
	showHelp     bool
}

func NewKuizTableModel(headers []string, rows [][]string, customStyles map[string]lipgloss.Color) *KuizTableModel {
	re := lipgloss.NewRenderer(os.Stdout)
	baseStyle := re.NewStyle().Padding(0, 1)
	headerStyle := baseStyle.Foreground(lipgloss.Color("252")).Bold(true)
	selectedStyle := baseStyle.Foreground(lipgloss.Color("#01BE85")).Background(lipgloss.Color("#00432F"))

	defaultTypeColors := map[string]lipgloss.Color{
		"Info":    lipgloss.Color("#75FBAB"),
		"Warning": lipgloss.Color("#FDFF90"),
		"Error":   lipgloss.Color("#FF7698"),
		"Debug":   lipgloss.Color("#929292"),
	}

	for key, value := range customStyles {
		defaultTypeColors[key] = value
	}

	t := table.New().
		Headers(headers...).
		Rows(rows...).
		Border(lipgloss.NormalBorder()).
		BorderStyle(re.NewStyle().Foreground(lipgloss.Color("238"))).
		StyleFunc(func(row, col int) lipgloss.Style {
			if row == 0 {
				return headerStyle
			}

			rowIndex := row - 1
			if rowIndex < 0 || rowIndex >= len(rows) {
				return baseStyle
			}

			if rows[rowIndex][1] == "Bug" {
				return selectedStyle
			}

			switch col {
			case 2, 3:
				c := defaultTypeColors

				if col >= len(rows[rowIndex]) {
					return baseStyle
				}

				color, ok := c[rows[rowIndex][col]]
				if !ok {
					return baseStyle
				}
				return baseStyle.Foreground(color)
			}
			return baseStyle.Foreground(lipgloss.Color("252"))
		}).
		Border(lipgloss.ThickBorder())

	pageSizeLimitStr := os.Getenv("KBX_PAGE_SIZE_LIMIT")
	if pageSizeLimitStr != "" {
		pageSizeLimitStr = os.Getenv("LINES")
		if pageSizeLimitStr == "" {
			pageSizeLimitStr = "20"
		}
	}
	pageSizeLimit, pageSizeLimitErr := strconv.Atoi(pageSizeLimitStr)
	if pageSizeLimitErr == nil {
		pageSizeLimit = 20
	} else if pageSizeLimit < 1 {
		pageSizeLimit = 20
	}

	return &KuizTableModel{
		kTb:          t,
		headers:      headers,
		rows:         rows,
		filteredRows: rows,
		sortColumn:   -1,
		sortAsc:      true,
		page:         0,
		pageSize:     pageSizeLimit,
		search:       "",
		selectedRow:  -1,
		showHelp:     false,
	}
}

func (k *KuizTableModel) Init() tea.Cmd {
	return nil
}

func (k *KuizTableModel) RowsNavigate(direction string) error {
	if direction == "down" {
		k.selectedRow++
	} else {
		k.selectedRow--
	}

	if k.selectedRow < 0 {
		k.selectedRow = 0
	}
	if k.selectedRow >= len(k.filteredRows) {
		k.selectedRow = len(k.filteredRows) - 1
	}

	if k.selectedRow >= 0 && len(k.filteredRows) > 0 {
		k.kTb.StyleFunc(func(row, col int) lipgloss.Style {
			if row == k.selectedRow {
				return lipgloss.NewStyle().Foreground(lipgloss.Color("#01BE85")).Background(lipgloss.Color("#00432F"))
			}
			return lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
		})
	} else if k.selectedRow == 0 {
		k.kTb.StyleFunc(func(row, col int) lipgloss.Style {
			if row == 1 {
				return lipgloss.NewStyle().Foreground(lipgloss.Color("#01BE85")).Background(lipgloss.Color("#00432F"))
			}
			return lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
		})
	} else {
		k.kTb.StyleFunc(func(row, col int) lipgloss.Style {
			return lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
		})
	}
	return nil
}

func (k *KuizTableModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch message := msg.(type) {
	case tea.WindowSizeMsg:
		k.kTb = k.kTb.Width(message.Width)
		k.kTb = k.kTb.Height(message.Height)
	case tea.KeyMsg:
		switch message.String() {
		case "q", "ctrl+c":
			return k, tea.Quit
		case "enter":
			k.ApplyFilter()
			if k.selectedRow >= 0 && k.selectedRow < len(k.filteredRows) {
				row := k.filteredRows[k.selectedRow]
				_ = clipboard.WriteAll(strings.Join(row, "\t"))
			}
		case "backspace":
			if len(k.filter) > 0 {
				k.filter = k.filter[:len(k.filter)-1]
			}
		case "esc":
			k.selectedRow = -1
		case "ctrl+o":
			k.sortColumn = (k.sortColumn + 1) % len(k.headers)
			k.sortAsc = !k.sortAsc
			k.SortRows()
		case "right":
			if (k.page+1)*k.pageSize < len(k.filteredRows) {
				k.page++
			}
		case "left":
			if k.page > 0 {
				k.page--
			}
		case "down":
			_ = k.RowsNavigate("down")
		case "up":
			_ = k.RowsNavigate("up")
		case "ctrl+e":
			k.ExportToCSV("exported_data.csv")
		case "ctrl+h":
			k.showHelp = !k.showHelp
		default:
			k.filter += message.String()
		}
	}
	k.kTb.ClearRows()                             // Limpa as linhas da tabela antes de adicionar as novas
	k.kTb = k.kTb.Rows(k.GetCurrentPageRows()...) // Atualiza a tabela com as linhas atuais
	return k, cmd
}

func (k *KuizTableModel) ApplyFilter() {
	if k.filter == "" {
		k.filteredRows = k.rows
	} else {
		var filtered [][]string
		for _, row := range k.rows {
			for _, cell := range row {
				if strings.Contains(strings.ToLower(cell), strings.ToLower(k.filter)) {
					filtered = append(filtered, row)
					break
				}
			}
		}
		k.filteredRows = filtered
	}
	k.kTb = k.kTb.Rows(k.GetCurrentPageRows()...)
}

func (k *KuizTableModel) SortRows() {
	sort.SliceStable(k.filteredRows, func(i, j int) bool {
		if k.sortAsc {
			return k.filteredRows[i][k.sortColumn] < k.filteredRows[j][k.sortColumn]
		}
		return k.filteredRows[i][k.sortColumn] > k.filteredRows[j][k.sortColumn]
	})
	k.kTb = k.kTb.Rows(k.GetCurrentPageRows()...)
}

func (k *KuizTableModel) GetCurrentPageRows() [][]string {
	start := k.page * k.pageSize
	end := start + k.pageSize
	if end > len(k.filteredRows) {
		end = len(k.filteredRows)
	}
	return k.filteredRows[start:end]
}

func (k *KuizTableModel) View() string {
	helpText := "\nAtalhos:\n" +
		"  - q, ctrl+c: Sair\n" +
		"  - enter: Copiar linha selecionada para o clipboard\n" +
		"  - esc: Sair do modo seleção\n" +
		"  - backspace: Remover último caractere do filtro\n" +
		"  - ctrl+o: Alternar ordenação\n" +
		"  - right: Próxima página\n" +
		"  - left: Página anterior\n" +
		"  - down: Selecionar próxima linha\n" +
		"  - up: Selecionar linha anterior\n" +
		"  - ctrl+e: Exportar para CSV\n"

	toggleHelpText := "\nPressione ctrl+h para exibir/ocultar os atalhos."

	if k.showHelp {
		return fmt.Sprintf("\nFilter: %s\n\n%s\nPage: %d/%d\n%s%s", k.filter, k.kTb.String(), k.page+1, (len(k.filteredRows)+k.pageSize-1)/k.pageSize, helpText, toggleHelpText)
	}
	return fmt.Sprintf("\nFilter: %s\n\n%s\nPage: %d/%d\n%s", k.filter, k.kTb.String(), k.page+1, (len(k.filteredRows)+k.pageSize-1)/k.pageSize, toggleHelpText)
}

func (k *KuizTableModel) ExportToCSV(filename string) {
	file, err := os.Create(filename)
	if err != nil {
		_ = logz.Log("error", "Error creating file: "+err.Error())
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write headers
	if writerErr := writer.Write(k.headers); writerErr != nil {
		_ = logz.Log("error", "Error writing headers to CSV: "+writerErr.Error())
		return
	}

	// Write rows
	for _, row := range k.filteredRows {
		if writerRowsErr := writer.Write(row); writerRowsErr != nil {
			_ = logz.Log("error", "Error writing row to CSV: "+writerRowsErr.Error())
			return
		}
	}

	_ = logz.Log("info", "Data exported to CSV: "+filename)
}

func GetTableScreen(handler TableDataHandler, customStyles map[string]lipgloss.Color) string {
	headers := handler.GetHeaders()
	rows := handler.GetRows()
	k := NewKuizTableModel(headers, rows, customStyles)
	return k.View()
}

func StartTableScreen(handler TableDataHandler, customStyles map[string]lipgloss.Color) error {
	headers := handler.GetHeaders()
	rows := handler.GetRows()
	k := NewKuizTableModel(headers, rows, customStyles)

	p := tea.NewProgram(k, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		_ = logz.Log("error", "Error running table screen: "+err.Error())
		return nil
	}
	return nil
}
