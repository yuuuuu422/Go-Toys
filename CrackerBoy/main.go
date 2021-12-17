package main

import (
	"CrackerBoy/common"
	"github.com/urfave/cli"
	"os"
)

func main() {
	app := cli.NewApp()
	app.Name="CrackerBoy"
	app.Author="Theoyu"
	app.Usage="A tool for cracking weak passwords."
	app.Commands=common.Commands
	app.Flags=common.Commands[0].Flags
	app.Run(os.Args)

}

