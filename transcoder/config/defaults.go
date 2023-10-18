package config

func DefaultOptions() TranscodeOptions {
	return TranscodeOptions{

		Outputs: Outputs{
			{
				Profile: Profile{
					Width:  1280,
					Height: 720,
					OutputOptions: map[string]string{
						"vcodec": "h264_qsv",
						"b:v":    "1800k",
						"f":      "mp4",
						//"movflags": "frag_keyframe+empty_moov",
						"map": "a",
						"c:a": "aac",
						"b:a": "128k",
					},
				},
				Filename: "output_720_1800k.mp4",
			},
			{
				Profile: Profile{
					Width:  854,
					Height: 480,
					OutputOptions: map[string]string{
						"vcodec": "h264_qsv",
						"b:v":    "1200k",
						"f":      "mp4",
						//"movflags": "frag_keyframe+empty_moov",
						"map": "a",
						"c:a": "aac",
						"b:a": "96k",
					},
				},
				Filename: "output_480_1200k.mp4",
			},
			{
				Profile: Profile{
					Width:  640,
					Height: 360,
					OutputOptions: map[string]string{
						"vcodec": "h264_qsv",
						"b:v":    "700k",
						"acodec": "aac",
						"f":      "mp4",
						//"movflags": "frag_keyframe+empty_moov",
						"map": "a",
						"c:a": "aac",
						"b:a": "96k",
					},
				},
				Filename: "output_360_700k.mp4",
			},
		},
	}

}
