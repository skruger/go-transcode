package config

import (
	"gopkg.in/yaml.v3"
)

type Profile struct {
	//ScaleOptions map[string]string
	Width         int               `mapstructure:"Width"`
	Height        int               `mapstructure:"Height"`
	OutputOptions map[string]string `mapstructure:"OutputOptions"`
}

type Output struct {
	Profile  Profile
	Filename string
}

type Outputs []Output

type TranscodeOptions struct {
	Outputs    Outputs  `mapstructure:"Outputs"`
	GlobalArgs []string `mapstructure:"GlobalArgs"`
}

type TranscodeRequest struct {
	TranscodeOptions TranscodeOptions `mapstructure:"TranscodeOptions"`
	InputFile        string           `mapstructure:"InputFile"`
}

func LoadTranscodeOptions(data []byte) (*TranscodeOptions, error) {
	options := &TranscodeOptions{}

	if err := yaml.Unmarshal(data, options); err != nil {
		return nil, err
	}
	return options, nil
}
