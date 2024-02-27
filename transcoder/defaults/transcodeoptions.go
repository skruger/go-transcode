package defaults

import (
	"embed"
	_ "embed"
	"fmt"
	log2 "github.com/labstack/gommon/log"
	"github.com/skruger/privatestudio/transcoder/config"
	"gopkg.in/yaml.v3"
	"log"
	"path"
	"strings"
)

//go:embed transcodeoptions/*.yaml
var fs embed.FS

func loadTranscodeOptions(yamlData []byte) config.TranscodeOptions {
	var profile config.TranscodeOptions
	err := yaml.Unmarshal(yamlData, &profile)
	if err != nil {
		log.Panicf("unable to unmarshal transcoding options: %s", err)
	}
	return profile
}

func GetDefaultTranscodeOptions() map[string]config.TranscodeOptions {
	options := map[string]config.TranscodeOptions{}
	entries, err := fs.ReadDir("transcodeoptions")
	if err != nil {
		log2.Errorf("Unable to load default transcode options: %s", err)
	}
	for _, option := range entries {
		filename := option.Name()
		if path.Ext(filename) == ".yaml" {
			log2.Infof("Found options file: %s", filename)
			name := strings.TrimSuffix(filename, path.Ext(filename))
			yamlBytes, readErr := fs.ReadFile(fmt.Sprintf("transcodeoptions/%s", filename))
			if readErr != nil {
				log2.Infof("Unable to load data for default profile '%s': %s", name, err)
			} else {
				options[name] = loadTranscodeOptions(yamlBytes)
			}
		}
	}
	return options
}
