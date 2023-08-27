package transcode

import (
	"fmt"
	ffmpeg_go "github.com/u2takey/ffmpeg-go"
	"go-transcode/config"
)

type TranscodeSession struct {
	InputFile string
}

func NewTranscodeSession(file string) *TranscodeSession {
	return &TranscodeSession{
		InputFile: file,
	}
}

func (ts *TranscodeSession) BuildTranscodeStream(options config.TranscodeOptions) (*ffmpeg_go.Stream, error) {
	input := ffmpeg_go.Input(ts.InputFile)
	stream := input.Split()

	outputStreams := make([]*ffmpeg_go.Stream, len(options.Outputs))
	for num, output := range options.Outputs {
		filtered := stream.Get(fmt.Sprintf("%d", num)).Filter(
			"scale",
			ffmpeg_go.Args{fmt.Sprintf("%d:%d", output.Profile.Width, output.Profile.Height)},
		)
		kwargs := ffmpeg_go.KwArgs{}
		for key, val := range output.Profile.OutputOptions {
			kwargs[key] = val
		}
		outputStreams[num] = filtered.Output(
			output.Filename,
			kwargs,
		)
	}

	outputStream := ffmpeg_go.MergeOutputs(outputStreams...)

	if len(options.GlobalArgs) > 0 {
		outputStream = outputStream.GlobalArgs(options.GlobalArgs...)
	}

	return outputStream.ErrorToStdOut(), nil
}
