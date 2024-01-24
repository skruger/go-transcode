package transcoder

import (
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
	"time"
)

//
//type VideoMetadata struct {
//	Filename       string
//	Codec          string
//	Filesize       int
//	DurationMs     int
//	DurationFrames int
//	Fps            float32
//	Width          int
//	Height         int
//}

type FfProbe struct {
	Streams []struct {
		Index            int    `json:"index"`
		CodecName        string `json:"codec_name"`
		CodecLongName    string `json:"codec_long_name"`
		Profile          string `json:"profile"`
		CodecType        string `json:"codec_type"`
		CodecTagString   string `json:"codec_tag_string"`
		CodecTag         string `json:"codec_tag"`
		Width            int    `json:"width"`
		Height           int    `json:"height"`
		CodedWidth       int    `json:"coded_width"`
		CodedHeight      int    `json:"coded_height"`
		ClosedCaptions   int    `json:"closed_captions"`
		HasBFrames       int    `json:"has_b_frames"`
		PixFmt           string `json:"pix_fmt"`
		Level            int    `json:"level"`
		ColorRange       string `json:"color_range"`
		ColorSpace       string `json:"color_space"`
		ColorTransfer    string `json:"color_transfer"`
		ColorPrimaries   string `json:"color_primaries"`
		ChromaLocation   string `json:"chroma_location"`
		Refs             int    `json:"refs"`
		IsAvc            string `json:"is_avc"`
		NalLengthSize    string `json:"nal_length_size"`
		RFrameRate       string `json:"r_frame_rate"`
		AvgFrameRate     string `json:"avg_frame_rate"`
		TimeBase         string `json:"time_base"`
		StartPts         int    `json:"start_pts"`
		StartTime        string `json:"start_time"`
		DurationTs       int    `json:"duration_ts"`
		Duration         string `json:"duration"`
		BitRate          string `json:"bit_rate"`
		BitsPerRawSample string `json:"bits_per_raw_sample"`
		NbFrames         string `json:"nb_frames"`
		Disposition      struct {
			Default         int `json:"default"`
			Dub             int `json:"dub"`
			Original        int `json:"original"`
			Comment         int `json:"comment"`
			Lyrics          int `json:"lyrics"`
			Karaoke         int `json:"karaoke"`
			Forced          int `json:"forced"`
			HearingImpaired int `json:"hearing_impaired"`
			VisualImpaired  int `json:"visual_impaired"`
			CleanEffects    int `json:"clean_effects"`
			AttachedPic     int `json:"attached_pic"`
			TimedThumbnails int `json:"timed_thumbnails"`
		} `json:"disposition"`
		Tags struct {
			CreationTime time.Time `json:"creation_time"`
			Language     string    `json:"language"`
			HandlerName  string    `json:"handler_name"`
			VendorId     string    `json:"vendor_id"`
		} `json:"tags"`
	} `json:"streams"`
}

func GetVideoMetadata(filename string) (*FfProbe, error) {
	args := []string{
		"-show_streams",
		"-select_streams", "v",
		"-print_format", "json",
		filename,
	}
	cmd := exec.Command("ffprobe", args...)
	if errors.Is(cmd.Err, exec.ErrDot) {
		cmd.Err = nil
	}
	data, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("ffprobe failed in output: %s", err)
	}
	ffMetadata := &FfProbe{}
	err = json.Unmarshal(data, ffMetadata)
	if err != nil {
		return nil, fmt.Errorf("unable to parse ffprobe output: %s", err)
	}
	return ffMetadata, nil
}
