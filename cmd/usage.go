package main

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/spf13/cobra"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func colorYellow(s string) string {
	return color.New(color.FgYellow).SprintFunc()(s)
}
func colorGreen(s string) string {
	return color.New(color.FgGreen).SprintFunc()(s)
}
func colorBlue(s string) string {
	return color.New(color.FgBlue).SprintFunc()(s)
}
func colorRed(s string) string {
	return color.New(color.FgRed).SprintFunc()(s)
}
func colorHelp(s string) string {
	return color.New(color.FgCyan).SprintFunc()(s)
}
func hasServiceCommands(cmds []*cobra.Command) bool {
	for _, cmd := range cmds {
		if cmd.Annotations["service"] == "true" {
			return true
		}
	}
	return false
}
func hasModuleCommands(cmds []*cobra.Command) bool {
	for _, cmd := range cmds {
		if cmd.Annotations["service"] != "true" {
			return true
		}
	}
	return false
}
func setUsageDefinition(cmd *cobra.Command) {
	cobra.AddTemplateFunc("colorYellow", colorYellow)
	cobra.AddTemplateFunc("colorGreen", colorGreen)
	cobra.AddTemplateFunc("colorRed", colorRed)
	cobra.AddTemplateFunc("colorBlue", colorBlue)
	cobra.AddTemplateFunc("colorHelp", colorHelp)
	cobra.AddTemplateFunc("hasServiceCommands", hasServiceCommands)
	cobra.AddTemplateFunc("hasModuleCommands", hasModuleCommands)

	// Altera o template de uso do cobra
	cmd.SetUsageTemplate(cliUsageTemplate)
}

var cliUsageTemplate = `{{- if index .Annotations "banner" }}{{colorBlue (index .Annotations "banner")}}{{end}}{{- if (index .Annotations "description") }}
{{index .Annotations "description"}}
{{- end }}

{{colorYellow "Usage:"}}{{if .Runnable}}
  {{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [command] [args]{{end}}{{if gt (len .Aliases) 0}}

{{colorYellow "Aliases:"}}
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

{{colorYellow "Example:"}}
  {{.Example}}{{end}}{{if .HasAvailableSubCommands}}
{{colorYellow "Available Commands:"}}{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{colorGreen (rpad .Name .NamePadding) }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

{{colorYellow "Flags:"}}
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces | colorHelp}}{{end}}{{if .HasAvailableInheritedFlags}}

{{colorYellow "Global Options:"}}
  {{.InheritedFlags.FlagUsages | trimTrailingWhitespaces | colorHelp}}{{end}}{{if .HasHelpSubCommands}}

{{colorYellow "Additional help topics:"}}
{{range .Commands}}{{if .IsHelpCommand}}
  {{colorGreen (rpad .CommandPath .CommandPathPadding) }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasSubCommands}}

{{colorYellow (printf "Use \"%s [command] --help\" for more information about a command." .CommandPath)}}{{end}}
`

const AppName = "xtui"

func installCheck() {
	usrEnvs := os.Environ()
	envPath := os.Getenv("PATH")
	usrEnvs = append(usrEnvs, fmt.Sprintf("PATH=%s", envPath))
	appBinPath, appBinPathErr := exec.LookPath(AppName)
	if appBinPathErr != nil {
		fmt.Printf("Error: %v\n", appBinPathErr)
		return
	}
	appBinPath = strings.Replace(appBinPath, AppName, "", 1)

}

func cacheFlagAskForInstall(action string) (bool, error) {
	cacheDir, cacheDirErr := os.UserCacheDir()
	if cacheDirErr != nil {
		cacheDir = os.TempDir()
		if cacheDir == "" {
			cacheDir = "/tmp"
		}
	}
	userName, userNameErr := os.UserHomeDir()
	if userNameErr != nil {
		userName = "user"
	}
	cacheFilePath := filepath.Join(cacheDir, AppName, userName)

	// Essa lógica será para gravar a opção do usuário de instalar o módulo ou não.
	// A opção será gravada em um arquivo na pasta de cache do usuário.
	// O arquivo só será lido/buscado se o módulo não estiver instalado.
	// Se o arquivo existir, não haverá pergunta ao usuário.
	// O parâmetro action será para definir se será criado, lido ou removido o arquivo.
	switch action {
	case "create":
		// Cria o arquivo de cache, não perguntando nada pois a pergunta será feita em outro mtodo, anterior a este.
		// O arquivo será criado sem conteúdo, somente para sinalizar que a pergunta foi feita reduzindo a leiura em
		// busca de conteúdo, somente a verificação da existência do arquivo já é o próprio valor esperacdo.
		if mkdirErr := os.MkdirAll(cacheFilePath, 0755); mkdirErr != nil {
			return false, mkdirErr
		}
	case "check":
		// Verifica se o arquivo de cache existe
		if _, statErr := os.Stat(cacheFilePath); statErr != nil {
			return false, statErr
		} else {
			return true, nil
		}
	case "remove":
		// Remove o arquivo de cache
		if removeErr := os.RemoveAll(cacheFilePath); removeErr != nil {
			return false, removeErr
		}
	}

	return false, nil
}
