package main

import (
	"cutedoc/core"
	"cutedoc/diagnostics"
	"cutedoc/manifest"
	"github.com/rjeczalik/notify"
	"net/http"
	"os"
	"path/filepath"
	_ "path/filepath"
	"strconv"
)

func watchRecursive(path string, channel chan notify.EventInfo) error {
	infoChannel := make(chan notify.EventInfo)

	// Watch main directory
	err := notify.Watch(filepath.Join(path, "..."), infoChannel, notify.All)
	if err != nil {
		return err
	}

	// Preprocess events
	go func() {
		for {
			info := <-infoChannel
			stat, err := os.Stat(info.Path())
			if err != nil || !stat.IsDir() {
				channel <- info
			}
		}
	}()

	return nil
}

func runServer(port int) error {
	// Parse manifest
	siteManifest, err := manifest.ParseSiteManifest(core.SiteManifestName)
	if err != nil {
		return err
	}

	// Start watcher
	channel := make(chan notify.EventInfo)
	err = watchRecursive(siteManifest.InputPath, channel)
	if err != nil {
		return err
	}

	// Start watch-event handler
	go func() {
		for {
			<-channel

			// File system has changed, generate new version
			err := runGenerator()
			if err != nil {
				diagnostics.PrintError(err, "failed to regenerate")
			}
		}
	}()

	// Generate initial version
	err = runGenerator()
	if err != nil {
		return err
	}

	// Start server
	server := http.FileServer(http.Dir(siteManifest.OutputPath))
	err = http.ListenAndServe(":"+strconv.Itoa(port), server)
	if err != nil {
		return err
	}

	return nil
}
