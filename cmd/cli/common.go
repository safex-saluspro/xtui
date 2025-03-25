package cli

import (
	"os"
	"strings"
)

func GetDescriptions(descriptionArg []string, _ bool) map[string]string {
	var description, banner string
	if descriptionArg != nil {
		if strings.Contains(strings.Join(os.Args[0:], ""), "-h") {
			description = descriptionArg[0]
		} else {
			description = descriptionArg[1]
		}
	} else {
		description = ""
	}

	banner = `
 ___    ___  _________   ___  ___   ___     
|\  \  /  /||\___   ___\|\  \|\  \ |\  \    
\ \  \/  / /\|___ \  \_|\ \  \\\  \\ \  \   
 \ \    / /      \ \  \  \ \  \\\  \\ \  \  
  /     \/        \ \  \  \ \  \\\  \\ \  \ 
 /  /\   \         \ \__\  \ \_______\\ \__\
/__/ /\ __\         \|__|   \|_______| \|__|
|__|/ \|__|                                 
`
	return map[string]string{"banner": banner, "description": description}
}
