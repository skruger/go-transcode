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

type HlsEncryption struct {
	Mode              string `mapstructure:"Mode"`
	Key               string `mapstructure:"Key"`
	IVMode            string `mapstructure:"IVMode"`
	KeyUri            string `mapstructure:"KeyUri"`
	KeyFormat         string `mapstructure:"KeyFormat"`
	KeyFormatVersions string `mapstructure:"KeyFormatVersions"`
}

type PackageHls struct {
	HlsVersion       string         `mapstructure:"HlsVersion"`
	OutputDir        string         `mapstrucutre:"OutputDir"`
	OutputSingleFile bool           `mapstructure:"OutputSingleFile"`
	AudioFormat      string         `mapstructure:"AudioFormat"`
	SegmentDuration  string         `mapstructure:"SegmentDuration"`
	BaseUrl          string         `mapstructure:"BaseUrl"`
	Encryption       *HlsEncryption `mapstructure:"Encryption"`
}

type TranscodeOptions struct {
	Outputs          Outputs     `mapstructure:"Outputs"`
	GlobalArgs       []string    `mapstructure:"GlobalArgs"`
	OverwriteOutputs bool        `mapstructure:"OverwriteOutputs"`
	PackageHls       *PackageHls `mapstructure:"PackageHls"`
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
