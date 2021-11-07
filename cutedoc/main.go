package main

import (
	"cutedoc/diagnostics"
	"github.com/alecthomas/kong"
)

type GenerateCommand struct {
}

func (cmd *GenerateCommand) Run() error {
	return runGenerator()
}

type ServeCommand struct {
	Port uint16 `arg:"" optional:"" help:"On which port to start the server" default:"9080"`
}

func (cmd *ServeCommand) Run() error {
	return runServer(int(cmd.Port))
}

var Cli struct {
	Generate GenerateCommand `cmd:"" help:"Generate the documentation"`
	Serve    ServeCommand    `cmd:"" help:"Start the development server"`
}

func main() {
	ctx := kong.Parse(&Cli)
	err := ctx.Run()
	diagnostics.HandleError(err)
}
