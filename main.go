package main

import (
	"github.com/faelmori/xtui/cmd"
	"os"
)

func main() {
	if err := cmd.RegX().Execute(os.Args); err != nil {
		os.Exit(1)
	}
}
