package transcode

import (
	"fmt"
	"github.com/skruger/privatestudio/transcoder/config"
	"github.com/skruger/privatestudio/transcoder/stages"
	ffmpeg_go "github.com/u2takey/ffmpeg-go"
)

type TranscodeSession struct {
	InputFile string
	Outputs   []stages.MediaOut
}

func NewTranscodeSession(file string) *TranscodeSession {
	return &TranscodeSession{
		InputFile: file,
	}
}

func (ts *TranscodeSession) BuildTranscodeStream(options config.TranscodeOptions) (*ffmpeg_go.Stream, error) {
	input := ffmpeg_go.Input(ts.InputFile)
	stream := input.Split()

	ts.Outputs = make([]stages.MediaOut, len(options.Outputs))
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
		ts.Outputs[num] = stages.MediaOut{
			MediaType: stages.OutputVideo,
			FileName:  output.Filename,
		}
	}

	outputStream := ffmpeg_go.MergeOutputs(outputStreams...)

	if len(options.GlobalArgs) > 0 {
		outputStream = outputStream.GlobalArgs(options.GlobalArgs...)
	}

	return outputStream.ErrorToStdOut(), nil
}
