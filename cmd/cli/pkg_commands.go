package cli

import (
	"fmt"
	"github.com/faelmori/logz"
	p "github.com/faelmori/xtui/packages"
	"github.com/faelmori/xtui/types"
	"github.com/spf13/cobra"
	"os"
	"reflect"
	"strings"
)

// appsCmdsList retorna uma lista de comandos Cobra relacionados a aplicativos.
// Retorna um slice de ponteiros para comandos Cobra e um erro, se houver.
func PkgCmdsList() []*cobra.Command {
	cmdAdd := appsCmdAdd()
	cmdList := appsCmdList()
	cmdCheckDeps := checkDepsCmd()
	cmdAddAppsShell := appsCmdAddShell()
	cmdGenInstScript := appsCmdGenInstScript()

	return []*cobra.Command{
		cmdAdd,
		cmdList,
		cmdCheckDeps,
		cmdAddAppsShell,
		cmdGenInstScript,
	}
}

// appsCmdAdd cria um comando Cobra para instalar um aplicativo.
// Retorna um ponteiro para o comando Cobra configurado.
func appsCmdAdd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "install",
		Aliases: []string{"i", "ins", "add"},
		Annotations: GetDescriptions(
			[]string{
				"Install an application",
				"Install an application from a file or a repository and add it to the system",
			},
			false,
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			nameFlagValue, _ := cmd.Flags().GetStringArray("name")
			pathFlagValue, _ := cmd.Flags().GetString("path")
			yesFlagValue, _ := cmd.Flags().GetBool("yes")
			quietFlagValue, _ := cmd.Flags().GetBool("quiet")
			newArgs := []string{strings.Join(nameFlagValue, " "), pathFlagValue, fmt.Sprintf("%t", yesFlagValue), fmt.Sprintf("%t", quietFlagValue)}
			args = append(args, newArgs...)

			availableProperties := getAvailableProperties()
			if len(availableProperties) > 0 {
				adaptedArgs := adaptArgsToProperties(args, availableProperties)
				return p.InstallApps(adaptedArgs...)
			}

			return p.InstallApps(args...)
		},
	}

	cmd.Flags().StringArrayP("name", "n", []string{}, "App name")
	cmd.Flags().StringP("path", "p", "", "App path")
	cmd.Flags().BoolP("yes", "y", false, "Automatic yes to prompts")
	cmd.Flags().BoolP("quiet", "q", false, "Quiet mode")

	return cmd
}

// appsCmdGenInstScript cria um comando Cobra para gerar um script de dependências.
// Retorna um ponteiro para o comando Cobra configurado.
func appsCmdGenInstScript() *cobra.Command {
	genInstScriptCmd := &cobra.Command{
		Use:    "genDepsScript",
		Hidden: true,
		Annotations: GetDescriptions(
			[]string{
				"Gera script de dependências",
				"Gera um script para verificar e instalar dependências",
			},
			false,
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			depsList, depsListErr := getDepsList()
			if depsListErr != nil {
				logz.Error("Error getting dependencies list: "+depsListErr.Error(), nil)
				return depsListErr
			}
			return GenDepsScriptHandler(depsList, args...)
		},
	}

	return genInstScriptCmd
}

// appsCmdAddShell cria um comando Cobra para instalar dependências.
// Retorna um ponteiro para o comando Cobra configurado.
func appsCmdAddShell() *cobra.Command {
	addAppsShell := &cobra.Command{
		Use: "ins-deps",
		Annotations: GetDescriptions(
			[]string{
				"Install dependencies",
				"Install all dependencies provided on the system",
			},
			false,
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			return InstallDepsHandler(args...)
		},
	}

	return addAppsShell
}

// appsCmdList cria um comando Cobra para listar aplicativos.
// Retorna um ponteiro para o comando Cobra configurado.
func appsCmdList() *cobra.Command {
	cmd := &cobra.Command{
		Use: "list",
		Annotations: GetDescriptions(
			[]string{
				"List system installed apps",
				"List all installed apps on the system in a interactive table with filters, export options and more",
			},
			false,
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			nameFlagValue, _ := cmd.Flags().GetStringArray("name")
			statusFlagValue, _ := cmd.Flags().GetString("status")
			methodFlagValue, _ := cmd.Flags().GetString("method")
			newArgs := []string{strings.Join(nameFlagValue, " "), statusFlagValue, methodFlagValue}
			args = append(args, newArgs...)

			availableProperties := getAvailableProperties()
			if len(availableProperties) > 0 {
				adaptedArgs := adaptArgsToProperties(args, availableProperties)
				return p.ShowInstalledAppsTable(adaptedArgs...)
			}

			return p.ShowInstalledAppsTable(args...)
		},
	}

	cmd.Flags().StringArrayP("name", "n", []string{}, "App name")
	cmd.Flags().StringP("status", "s", "", "App status")
	cmd.Flags().StringP("method", "m", "", "App method")

	return cmd
}

// checkDepsCmd cria um comando Cobra para verificar dependências.
// Retorna um ponteiro para o comando Cobra configurado.
func checkDepsCmd() *cobra.Command {
	depsCmd := &cobra.Command{
		Use: "checkDeps",
		Annotations: GetDescriptions(
			[]string{
				"Ensure dependencies",
				"Ensure that all dependencies are installed",
			},
			false,
		),
		RunE: func(cmd *cobra.Command, args []string) error {
			return CheckDepsHandler(args...)
		},
	}

	return depsCmd
}

func getDepsList() ([]string, error) {
	if len(os.Args) == 0 {
		return nil, fmt.Errorf("nenhuma dependência informada")
	}
	for i, dep := range os.Args {
		// Verifica se é um slice de strings
		if reflect.TypeOf(dep).String() == "[]string" {
			return os.Args[i+1:], nil
		}
	}
	return nil, fmt.Errorf("nenhuma dependência informada")
}

// CheckDepsHandler verifica as dependências a partir dos argumentos fornecidos.
// Recebe um slice de strings com os argumentos.
// Retorna um erro, se houver.
func CheckDepsHandler(args ...string) error {
	var validationFilePath, version string
	if len(args) == 0 {
		return fmt.Errorf("caminho do arquivo de validação não informado")
	}
	validationFilePath = args[0]
	if len(args) > 1 {
		version = args[1]
	} else {
		version = "latest"
	}
	p.CheckDeps(validationFilePath, version)
	return nil
}

// GenDepsScriptHandler gera um script para verificar e instalar dependências a partir dos argumentos fornecidos.
// Recebe um slice de strings com a lista de dependências e um slice de strings com os argumentos.
// Retorna um erro, se houver.
func GenDepsScriptHandler(depsList []string, args ...string) error {
	var scriptPath, validationFilePath, version string
	if len(args) < 4 {
		return fmt.Errorf("erro ao ler argumentos")
	}
	scriptPath = args[len(args)-3]
	validationFilePath = args[len(args)-2]
	version = args[len(args)-1]
	return p.GenDepsScript(depsList, scriptPath, validationFilePath, version)
}

// InstallDepsHandler instala dependências a partir dos argumentos fornecidos.
// Recebe um slice de strings com os argumentos.
// Retorna um erro, se houver.
func InstallDepsHandler(args ...string) error {
	var scriptPath string
	if len(args) == 0 {
		return fmt.Errorf("caminho do script de instalação não informado")
	}
	scriptPath = args[0]
	return p.InstallApps(scriptPath)
}

func getAvailableProperties() map[string]string {
	return map[string]string{
		"property1": "value1",
		"property2": "value2",
	}
}

func adaptArgsToProperties(args []string, properties map[string]string) []string {
	adaptedArgs := args
	for key, value := range properties {
		adaptedArgs = append(adaptedArgs, fmt.Sprintf("--%s=%s", key, value))
	}
	return adaptedArgs
}
