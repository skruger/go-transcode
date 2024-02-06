package defaults

import (
	_ "embed"
	"github.com/skruger/privatestudio/transcoder/config"
	"gopkg.in/yaml.v3"
	"log"
)

//go:embed transcodeoptions/low-qsv.yaml
var lowQsvYaml []byte

//go:embed transcodeoptions/medium-qsv.yaml
var mediumQsvYaml []byte

//go:embed transcodeoptions/low-x264.yaml
var lowX264Yaml []byte

func loadTranscodeOptions(yamlData []byte) config.TranscodeOptions {
	var profile config.TranscodeOptions
	err := yaml.Unmarshal(yamlData, &profile)
	if err != nil {
		log.Panicf("unable to unmarshal transcoding options: %s", err)
	}
	return profile
}

func GetDefaultTranscodeOptions() map[string]config.TranscodeOptions {
	return map[string]config.TranscodeOptions{
		"low_qsv":    loadTranscodeOptions(lowQsvYaml),
		"medium_qsv": loadTranscodeOptions(mediumQsvYaml),
		"low_x264":   loadTranscodeOptions(lowX264Yaml),
	}
}
