package main

import (
	"flag"
	"go-transcode/transcode"
	"log"
	"time"
)

func main() {
	//outputPathPtr := flag.String("output", "./output", "Output Folder")
	flag.Parse()

	ts := transcode.NewTranscodeSession(flag.Args()[0])

	outputs := []transcode.Output{
		{
			Profile: transcode.Profile{
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
			Profile: transcode.Profile{
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
			Profile: transcode.Profile{
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
	}
	cmd, err := ts.BuildTranscodeStream(outputs)
	if err != nil {
		log.Panic(err)
	}

	log.Print(cmd.Args)

	start := time.Now()

	if runerr := cmd.Run(); runerr != nil {
		log.Print(runerr)
	}

	duration := time.Now().Sub(start)
	log.Printf("Transcoding complete in %v", duration)
}
