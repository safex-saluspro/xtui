package main

import (
	"github.com/faelmori/xtui/cmd"
	"os"
)

func main() {
	x := cmd.RegX()
	x.Execute(os.Args)
}
