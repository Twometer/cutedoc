package main

import (
	"cutedoc/core"
	"cutedoc/diagnostics"
	"cutedoc/manifest"
	"github.com/alecthomas/kong"
	"net/http"
	"strconv"
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
	err := runGenerator()
	if err != nil {
		return err
	}

	siteManifest, err := manifest.ParseSiteManifest(core.SiteManifestName)
	if err != nil {
		return err
	}

	server := http.FileServer(http.Dir(siteManifest.OutputPath))
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
	diagnostics.HandleError(err)
}
