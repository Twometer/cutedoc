package main

import (
	"cutedoc/manifest"
	"cutedoc/utils"
	"github.com/alecthomas/kong"
	"net/http"
	"strconv"
)

const themeManifestName = "theme.ini"
const sourceManifestName = "cutedoc.ini"

type GenerateCommand struct {
}

func (cmd *GenerateCommand) Run() error {
	return runGenerator()
}

type ServeCommand struct {
	Port uint16 `arg:"" optional:"" help:"On which port to start the server" default:"9080"`
}

func (cmd *ServeCommand) Run() error {
	err := runGenerator()
	if err != nil {
		return err
	}

	sourceManifest, err := manifest.ParseSourceManifest(sourceManifestName)
	if err != nil {
		return err
	}

	server := http.FileServer(http.Dir(sourceManifest.OutputPath))
	err = http.ListenAndServe(":"+strconv.Itoa(int(cmd.Port)), server)
	if err != nil {
		return err
	}

	return nil
}

var Cli struct {
	Generate GenerateCommand `cmd:"" help:"Generate the documentation"`
	Serve    ServeCommand    `cmd:"" help:"Start the development server"`
}

func main() {
	ctx := kong.Parse(&Cli)
	err := ctx.Run()
	utils.HandleError(err)
}
