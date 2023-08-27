package main

import (
	"flag"
	"go-transcode/config"
	"go-transcode/transcode"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"time"
)

func main() {
	//outputPathPtr := flag.String("output", "./output", "Output Folder")
	tsOptionsPtr := flag.String("transcode-options", "", "JSON file containing outputs configuration")
	tsWriteOptionsPtr := flag.String("write-default-transcode-options", "", "File name to save default outputs configuration to")
	tsDryRunPtr := flag.Bool("dry-run", false, "Generate job definition and exit after printing the expected command line")
	flag.Parse()

	if *tsWriteOptionsPtr != "" {
		file, err := os.Create(*tsWriteOptionsPtr)
		if err != nil {
			log.Panicf("Can't write output file %s: %s", tsWriteOptionsPtr, err)
		}
		defer file.Close()
		defaultsBytes, _ := yaml.Marshal(config.DefaultOptions())
		file.Write(defaultsBytes)
		return
	}

	options := config.DefaultOptions()

	if *tsOptionsPtr != "" {
		data, err := os.ReadFile(*tsOptionsPtr)
		if err != nil {
			log.Panicf("unable to open config file %s: %s", *tsOptionsPtr, err)
		}

		fileOptions, err := config.LoadTranscodeOptions(data)
		if err != nil {
			log.Panicf("unable to parse transcode options file %s: %s", *tsOptionsPtr, err)
		}
		options = *fileOptions
	}

	inFileName := flag.Args()[0]
	log.Printf("transcoding %s", inFileName)

	ts := transcode.NewTranscodeSession(inFileName)

	tsStream, err := ts.BuildTranscodeStream(options)
	if err != nil {
		log.Panic(err)
	}

	cmd := tsStream.Compile()

	log.Print(cmd.Args)

	if *tsDryRunPtr {
		return
	}

	start := time.Now()

	if runerr := cmd.Run(); runerr != nil {
		log.Print(runerr)
	}

	duration := time.Now().Sub(start)
	log.Printf("Transcoding complete in %v", duration)
}
