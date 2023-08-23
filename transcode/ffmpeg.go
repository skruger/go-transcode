package transcode

import (
	"fmt"
	ffmpeg_go "github.com/u2takey/ffmpeg-go"
	"os/exec"
)

type Profile struct {
	//ScaleOptions map[string]string
	Width         int               `json:"Width"`
	Height        int               `json:"Height"`
	OutputOptions map[string]string `json:"OutputOptions"`
}

type Output struct {
	Profile  Profile
	Filename string
}

type TranscodeSession struct {
	InputFile string
}

func NewTranscodeSession(file string) *TranscodeSession {
	return &TranscodeSession{
		InputFile: file,
	}
}

func (ts *TranscodeSession) BuildTranscodeStream(outputs []Output) (*exec.Cmd, error) {
	//inputStreams := make([]*ffmpeg_go.Stream, len(ts.InputFiles))
	//for num, file := range ts.InputFiles {
	//	inputStreams[num] = ffmpeg_go.Input(file)
	//}
	//
	//stream := ffmpeg_go.Concat(inputStreams).Split()
	input := ffmpeg_go.Input(ts.InputFile)
	stream := input.Split()

	outputStreams := make([]*ffmpeg_go.Stream, len(outputs))
	for num, output := range outputs {
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

	//audioKwargs := ffmpeg_go.KwArgs{
	//	"b:a":      "128k",
	//	"c:a":      "aac",
	//	"f":        "mp4",
	//	"movflags": "frag_keyframe+empty_moov",
	//}

	//outputStreams[len(outputs)] = input.
	//	Audio().Output("output_audio.mp4", audioKwargs)

	return ffmpeg_go.MergeOutputs(outputStreams...).ErrorToStdOut().Compile(), nil
}
