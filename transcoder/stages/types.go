package stages

const (
	OutputAudio = iota
	OutputVideo = iota
	OutputHls   = iota
)

type MediaOut struct {
	MediaType int
	FileName  string
	Width     int
	Height    int
}
