package main

import (
	"fmt"
	"github.com/seamounts/go-mod-tidy/cmd"
	"os"
)

func main() {
	if err := cmd.GoMod.Execute(); err != nil {
		fmt.Println(err)

		os.Exit(1)
	}
}
