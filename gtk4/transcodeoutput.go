package gtk4

import (
	"bufio"
	_ "embed"
	"fmt"
	"github.com/diamondburned/gotk4/pkg/glib/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/labstack/gommon/log"
	"github.com/skruger/privatestudio/dao"
	"github.com/skruger/privatestudio/transcoder/config"
	"github.com/skruger/privatestudio/transcoder/transcode"
	"os"
	"path"
)

//go:embed transcodeoutput.ui
var txOutputXML string

type transcodeOutput struct {
	window      *gtk.Window
	textOutput  *gtk.TextView
	progressBar *gtk.ProgressBar
	statusBar   *gtk.Statusbar
	builder     *gtk.Builder
	done        bool
	sourceAsset *dao.SourceAsset
	dbInstance  *dao.DaoInstance
	buffer      *gtk.TextBuffer
}

func newTranscodeOutputWindow(source *dao.SourceAsset, dbInstance *dao.DaoInstance) (*transcodeOutput, error) {
	builder := gtk.NewBuilderFromString(txOutputXML, len(txOutputXML))
	window := builder.GetObject("txOutputWindow").Cast().(*gtk.Window)
	textOutput := builder.GetObject("textOutput").Cast().(*gtk.TextView)
	progressBar := builder.GetObject("progressBar").Cast().(*gtk.ProgressBar)
	statusBar := builder.GetObject("statusBar").Cast().(*gtk.Statusbar)
	closeBtn := builder.GetObject("closeBtn").Cast().(*gtk.Button)

	output := &transcodeOutput{
		window:      window,
		textOutput:  textOutput,
		progressBar: progressBar,
		statusBar:   statusBar,
		builder:     builder,
		done:        false,
		sourceAsset: source,
		dbInstance:  dbInstance,
		buffer:      gtk.NewTextBuffer(gtk.NewTextTagTable()),
	}
	textOutput.SetEditable(false)
	output.buffer.SetText("")
	textOutput.SetBuffer(output.buffer)
	closeBtn.ConnectClicked(output.clickClose)
	output.window.ConnectCloseRequest(func() (ok bool) {
		if !output.done {
			statusBar.Push(0, "Unable to close window, transcode not complete")
		}
		return !output.done
	})

	return output, nil
}

func (t *transcodeOutput) clickClose() {
	if t.done {
		t.window.Close()

	} else {
		t.statusBar.Push(0, "Unable to close window, transcode not complete")
	}

}

func (t *transcodeOutput) startTranscode(profileName string, profile *config.TranscodeOptions) {
	defer t.finish()
	t.statusBar.Push(0, fmt.Sprintf("Starting transcode with profile: %s", profileName))

	t.addLine(fmt.Sprintf("Starting transcode of %s with profile %s\n", t.sourceAsset.Filename, profileName))
	//t.sourceAsset.Filename
	outputPath := path.Join("transcode_output", fmt.Sprintf("sourceAsset_%d", t.sourceAsset.Id), profileName)
	os.MkdirAll(outputPath, os.ModeDir|0700)
	session := transcode.NewTranscodeSession(t.sourceAsset.Filename, &outputPath)
	stream, err := session.BuildTranscodeStream(*profile)
	if err != nil {
		t.addLine(fmt.Sprintf("Unable to build transcode stream: %s\n", err))
	}
	cmd := stream.Compile()
	outPipe, err := cmd.StdoutPipe()
	if err != nil {
		t.addLine(fmt.Sprintf("error opening stdout pipe: %s\n", err))
	}
	errPipe, err := cmd.StderrPipe()
	if err != nil {
		t.addLine(fmt.Sprintf("error opening stdoerr pipe: %s\n", err))
	}

	t.addLine(cmd.String() + "\n")

	outDone := make(chan interface{})
	errDone := make(chan interface{})

	go func() {
		for {
			if cmd.Process != nil {
				break
			}
		}
		outReader := bufio.NewReader(outPipe)
		for {
			line, err := outReader.ReadString('\n')
			if err != nil {
				if err.Error() != "EOF" {
					t.addLine(fmt.Sprintf("outReader error: %s", err))
				}
				break
			}
			t.addLine(line)
		}
		close(outDone)
	}()

	go func() {
		for {
			if cmd.Process != nil {
				break
			}
		}
		errReader := bufio.NewReader(errPipe)
		for {
			line, err := errReader.ReadString('\n')
			if err != nil {
				if err.Error() != "EOF" {
					t.addLine(fmt.Sprintf("errReader error: %s\n", err))
				}
				break
			}
			t.addLine(line)
		}
		close(errDone)
	}()

	err = cmd.Start()
	if err != nil {
		t.addLine(fmt.Sprintf("Run error: %s\n", err))
		return
	}

	err = cmd.Wait()
	if err != nil {
		t.addLine(fmt.Sprintf("wait error: %s\n", err))
	}

	existingOutputs, err := t.dbInstance.GetTranscodeOutputs("where source=?", t.sourceAsset.Id)
	if err != nil {
		log.Errorf("unable to get existing transcode outputs for source asset %s: %s", t.sourceAsset.Filename, err)
		return
	}

	for _, output := range session.Outputs {
		outputFound := false
		for _, existingOutput := range existingOutputs {
			if existingOutput.Filename == output.FileName {
				outputFound = true
			}
		}
		if !outputFound {
			fileInfo, err := os.Stat(output.FileName)
			filesize := int(fileInfo.Size())
			if err == nil && filesize > 0 {
				t.dbInstance.NewTranscodeOutput(output.FileName, filesize, t.sourceAsset, profileName,
					dao.Resolution{Width: output.Width, Height: output.Height})
			}
		}
	}

	<-errDone
	<-outDone
	t.addLine("\n\nTranscode Complete\n")
	t.statusBar.Push(0, "Transcode Complete")

}

func (t *transcodeOutput) addLine(line string) {
	glib.IdleAdd(func() {
		buffer := t.textOutput.Buffer()
		end := buffer.EndIter()
		buffer.Insert(end, line)
		//scrollWindow := t.builder.GetObject("textScroll").Cast().(*gtk.ScrolledWindow)
		//scrollBar := scrollWindow.VScrollbar()
		//scrollBar.
	})
}

func (t *transcodeOutput) finish() {
	t.done = true
}
