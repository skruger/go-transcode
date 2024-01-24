package gtk4

import (
	"database/sql"
	_ "embed"
	"fmt"
	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/labstack/gommon/log"
	"github.com/skruger/privatestudio/dao"
	"github.com/skruger/privatestudio/transcoder"
	"os"
	"strconv"
	"strings"
)

type ColumnType int

const (
	sourceFilenameColumn ColumnType = iota
	sourceNoteColumn
)

//go:embed privatestudio.ui
var psXML string

//go:embed sourcelistitem.ui
var listItemXML []byte

type windowState struct {
	initialized    bool
	builder        *gtk.Builder
	window         *gtk.ApplicationWindow
	app            *gtk.Application
	status         *gtk.Statusbar
	sourceList     *sourceListData
	statusMessages chan string
	daoInstance    *dao.DaoInstance
}

func newWindowState(db *sql.DB) *windowState {
	state := &windowState{
		app:         gtk.NewApplication("com.shaunkruger.privatestudio", gio.ApplicationFlagsNone),
		initialized: false,
		daoInstance: dao.NewDaoInstance(db),
	}
	state.app.ConnectActivate(func() {
		builder := gtk.NewBuilderFromString(psXML, len(psXML))
		window := builder.GetObject("psMainWindow").Cast().(*gtk.ApplicationWindow)
		window.SetApplication(state.app)
		window.SetDefaultSize(1024, 768)
		window.Show()
		status := builder.GetObject("statusBar").Cast().(*gtk.Statusbar)
		status.Push(0, "This is a status!")

		sourcesList := builder.GetObject("sourcesList").Cast().(*gtk.ListView)
		sourceListState := newSourceList(sourcesList)

		state.builder = builder
		state.window = window
		state.status = status
		state.sourceList = sourceListState

		sourceListState.itemSelection.ConnectSelectionChanged(state.sourceSelectionChanged)
		state.initialized = true

		openBtn := builder.GetObject("openFileBtn").Cast().(*gtk.Button)
		openBtn.ConnectClicked(state.openFile)
		state.loadFiles()
	})
	return state
}

func (w *windowState) sourceSelectionChanged(position, nItems uint) {
	w.status.Push(0, fmt.Sprintf("Selected: %d", position))

}

func (w *windowState) openFile() {
	//openDialog := &gtk.FileChooserDialog{}

	openDialog := gtk.NewFileChooserNative("Select video", &w.window.Window, gtk.FileChooserActionOpen, "Open", "Cancel")
	openDialog.SetSelectMultiple(false)
	openDialog.Show()
	openDialog.ConnectResponse(func(responseId int) {
		file := openDialog.File().Path()
		log.Infof("Got response event from openDialog: %d - %s", responseId, file)
		w.addFile(file)
	})
}

func (w *windowState) loadFiles() {
	assets, err := w.daoInstance.GetSourceAssets("ORDER BY filename")
	if err != nil {
		log.Errorf("loadFiles error: %s", err)
	}
	for _, asset := range assets {
		w.addSourceAsset(asset.Filename, asset)
	}
}

func (w *windowState) addFile(file string) {
	var sa *dao.SourceAsset
	sa, err := w.daoInstance.GetSourceAssetByFilename(file)
	if err != nil {
		sa, err = w.daoInstance.NewSourceAsset(file)
		if err != nil {
			w.status.Push(0, fmt.Sprintf("Unable to add file to DB: %s", err))
			return
		}
	}
	w.addSourceAsset(file, sa)
}

func (w *windowState) addSourceAsset(file string, sa *dao.SourceAsset) {

	w.status.Push(0, fmt.Sprintf("Opened file: %s (%d)", file, sa.Id))
	go func() {
		ffProbe, err := transcoder.GetVideoMetadata(file)
		if err != nil {
			w.status.Push(0, fmt.Sprintf("File %s: %s", file, err))
			return
		}
		fpsParts := strings.Split(ffProbe.Streams[0].RFrameRate, "/")
		if len(fpsParts) >= 2 {
			frames, err1 := strconv.ParseFloat(fpsParts[0], 32)
			seconds, err2 := strconv.ParseFloat(fpsParts[1], 32)
			if err1 != nil || err2 != nil {
				log.Errorf("unable to parse fps from '%s': frames=%s, seconds=%s", ffProbe.Streams[0].RFrameRate, err1, err2)
			}
			sa.Fps = float32(frames) / float32(seconds)
		}
		sa.Codec = ffProbe.Streams[0].CodecName
		sa.Resolution.Width = ffProbe.Streams[0].Width
		sa.Resolution.Height = ffProbe.Streams[0].Height
		duration64, err := strconv.ParseFloat(ffProbe.Streams[0].Duration, 32)
		if err != nil {
			log.Errorf("unable to parse duration: %s", err)
		}
		sa.DurationTime = float32(duration64)
		frames, err := strconv.Atoi(ffProbe.Streams[0].NbFrames)
		if err != nil {
			log.Errorf("unable to parse nb_frames: %s", err)
		}
		fileInfo, err := os.Stat(file)
		if err != nil {
			log.Errorf("unable to stat %s: %s", file, err)
		} else {
			sa.Filesize = int(fileInfo.Size())
		}
		sa.DurationFrames = frames
		err = sa.Save()
		if err != nil {
			log.Errorf("unable to save source asset: %s", err)
		}
		w.sourceList.add(sa)

	}()

}

func RunUI(db *sql.DB) {
	state := newWindowState(db)

	if code := state.app.Run(os.Args); code > 0 {
		os.Exit(code)
	}
}
