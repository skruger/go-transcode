package main

import (
	"github.com/skruger/privatestudio/web"
	config2 "github.com/skruger/privatestudio/web/config"
	"github.com/skruger/privatestudio/web/datastore/localfile"
	"log"
	"os"
)

func main() {

	cfg, err := config2.LoadWebConfig(os.Args)
	if err != nil {
		log.Panicf("unable to load configuration: %s", err)
	}

	sourcePath := cfg.GetSourcePath()
	transcodePath := cfg.GetTranscodePath()
	manifestPath := cfg.GetManifestPath()
	_ = os.MkdirAll(sourcePath, 0755)
	_ = os.MkdirAll(transcodePath, 0755)
	_ = os.MkdirAll(manifestPath, 0755)
	assetStorage, _ := localfile.NewLocalFileAssetStorage(sourcePath, transcodePath, manifestPath)
	config := web.NewServerConfig(assetStorage)

	server, err := web.NewServer(config)
	if err != nil {
		log.Panicf("unable to get echo server: %s", err)
	}
	err = server.Start(":1313")
	if err != nil {
		log.Panicf("unable to listen: %s", err)
	}
}
