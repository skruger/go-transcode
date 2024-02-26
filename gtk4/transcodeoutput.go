package gtk4

import (
	_ "embed"
	"fmt"
	"github.com/diamondburned/gotk4/pkg/glib/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/labstack/gommon/log"
	"github.com/skruger/privatestudio/dao"
	"github.com/skruger/privatestudio/transcoder"
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
	t.addStatus(fmt.Sprintf("Starting transcode with profile: %s", profileName))

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

	t.addLine(fmt.Sprintf("transcode command: %s\n", cmd.String()))

	runner := transcoder.NewTranscodeRunner(cmd)

	err = runner.Start()
	if err != nil {
		t.addLine(fmt.Sprintf("error starting transcode: %s", err))
		return
	}

	go func() {
		for {
			output, finished := runner.ReceiveLine()
			if output != nil {
				if output.IsStatus() {
					t.addStatus(output.Data)
				} else {
					t.addLine(output.Data)
				}

			}
			if finished {
				break
			}
		}
	}()

	err = runner.Wait()
	if err != nil {
		t.addLine(fmt.Sprintf("Error waiting for transcode to finish: %s", err))
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
			fileInfo, statErr := os.Stat(output.FileName)
			if statErr != nil {
				t.addLine(fmt.Sprintf("Unable to stat output file: %s\n", statErr))
			} else {
				filesize := int(fileInfo.Size())
				if filesize > 0 {
					t.dbInstance.NewTranscodeOutput(output.FileName, filesize, t.sourceAsset, profileName,
						dao.Resolution{Width: output.Width, Height: output.Height})
				}

			}
		}
	}

	t.addLine("\n\nTranscode Complete\n")
	t.addStatus("Transcode Complete")

}

func (t *transcodeOutput) addLine(line string) {
	glib.IdleAdd(func() {
		buffer := t.textOutput.Buffer()
		end := buffer.EndIter()
		buffer.Insert(end, line)
	})
}

func (t *transcodeOutput) addStatus(line string) {
	glib.IdleAdd(func() {
		t.statusBar.Push(0, line)
	})
}

func (t *transcodeOutput) finish() {
	t.done = true
}
