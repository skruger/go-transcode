package config

import "github.com/skruger/privatestudio/transcoder/config"

var DefaultConfig = ConfigFile{
	AssetHome:   "assets",
	DatabaseURL: "",
	EncodingProfiles: map[string]StandardEncodingProfile{
		"basic_qsv": {
			Outputs: config.Outputs{
				{
					Profile: config.Profile{
						Width:  1280,
						Height: 720,
						OutputOptions: map[string]string{
							"b:a":    "96k",
							"c:a":    "aac",
							"b:v":    "1800k",
							"vcodec": "h264_qsv",
							"f":      "mp4",
							"map":    "a",
						},
					},
					Filename: "output_720_1800k.mp4",
				},
				{
					Profile: config.Profile{
						Width:  854,
						Height: 480,
						OutputOptions: map[string]string{
							"b:a":    "96k",
							"c:a":    "aac",
							"b:v":    "1200k",
							"vcodec": "h264_qsv",
							"f":      "mp4",
							"map":    "a",
						},
					},
					Filename: "output_480_1200k.mp4",
				},
				{
					Profile: config.Profile{
						Width:  640,
						Height: 360,
						OutputOptions: map[string]string{
							"b:a":    "96k",
							"c:a":    "aac",
							"b:v":    "700k",
							"vcodec": "h264_qsv",
							"f":      "mp4",
							"map":    "a",
						},
					},
					Filename: "output_360_700k.mp4",
				},
			},
			GlobalArgs: []string{
				"-y",
			},
		},
	},
}
