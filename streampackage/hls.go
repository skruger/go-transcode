package streampackage

import (
	"fmt"
	"go-transcode/config"
	"go-transcode/stages"
	"os/exec"
)

type HlsPackage struct {
	inputFiles []stages.MediaOut
}

func NewHlsPackage(mediaInfo []stages.MediaOut) *HlsPackage {
	return &HlsPackage{
		inputFiles: mediaInfo,
	}
}

func (h *HlsPackage) BuildPackageCommand(packageSettings config.PackageHls) (*exec.Cmd, error) {
	args := []string{"-v"}
	if packageSettings.HlsVersion != "" {
		args = append(args, fmt.Sprintf("--hls-version=%s", packageSettings.HlsVersion))
	}
	if packageSettings.OutputDir != "" {
		args = append(args, fmt.Sprintf("--output-dir=%s", packageSettings.OutputDir))
	}
	if packageSettings.OutputSingleFile {
		args = append(args, "--output-single-file")
	}
	if packageSettings.AudioFormat != "" {
		args = append(args, fmt.Sprintf("--audio-format=%s", packageSettings.AudioFormat))
	}
	if packageSettings.SegmentDuration != "" {
		args = append(args, fmt.Sprintf("--segment-duration=%s", packageSettings.SegmentDuration))
	}
	if packageSettings.BaseUrl != "" {
		args = append(args, fmt.Sprintf("--base-url=%s", packageSettings.BaseUrl))
	}
	for _, media := range h.inputFiles {
		args = append(args, media.FileName)
	}
	cmd := exec.Command("mp4hls", args...)
	return cmd, nil
}
