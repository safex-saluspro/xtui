package main

import (
	"os"
)

func main() {
	if err := RegX().Execute(); err != nil {
		os.Exit(1)
	}
}
