package main

import (
	"github.com/neucn/ipgw/pkg/cmd"
	"github.com/neucn/ipgw/pkg/console"
	"os"
)

func main() {
	if err := cmd.App.Run(os.Args); err != nil {
		console.FatalL(err.Error())
	}
}
