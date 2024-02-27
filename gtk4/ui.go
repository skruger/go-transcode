package gtk4

import (
	"database/sql"
	_ "embed"
	"fmt"
	"github.com/diamondburned/gotk4/pkg/gio/v2"
	"github.com/diamondburned/gotk4/pkg/glib/v2"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/labstack/gommon/log"
	"github.com/skruger/privatestudio/dao"
	"github.com/skruger/privatestudio/transcoder"
	"github.com/skruger/privatestudio/transcoder/defaults"
	"os"
	"strconv"
	"strings"
	"time"
)

type ColumnType int

const (
	sourceFilenameColumn ColumnType = iota
	sourceNoteColumn
)

//go:embed privatestudio.ui
var psXML string

type windowState struct {
	initialized    bool
	builder        *gtk.Builder
	window         *gtk.ApplicationWindow
	app            *gtk.Application
	status         *gtk.Statusbar
	sourceList     *sourceListData
	outputList     *outputListData
	profileList    *transcodeOptionData
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
		sourceListState := newSourceList(sourcesList, state.sourceSelected)

		outputList := builder.GetObject("outputList").Cast().(*gtk.ListView)
		outputListState := newOutputList(outputList)

		state.builder = builder
		state.window = window
		state.status = status
		state.sourceList = sourceListState
		state.outputList = outputListState

		openBtn := builder.GetObject("openFileBtn").Cast().(*gtk.Button)
		openBtn.ConnectClicked(state.openFile)

		txProfiles := builder.GetObject("transcodeProfileDropDown").Cast().(*gtk.DropDown)
		state.profileList = newTranscodeOptionList(txProfiles)
		txBtn := builder.GetObject("transcodeBtn").Cast().(*gtk.Button)
		txBtn.ConnectClicked(state.transcodeFile)

		packageBtn := builder.GetObject("outputPackageBtn").Cast().(*gtk.Button)
		packageBtn.ConnectClicked(state.packageOutput)

		deleteOutputBtn := builder.GetObject("outputDeleteBtn").Cast().(*gtk.Button)
		deleteOutputBtn.ConnectClicked(state.deleteOutput)

		state.loadProfiles()
		state.loadFiles()

		state.initialized = true
		go func() {
			time.Sleep(100 * time.Millisecond)
			glib.IdleAdd(func() {
				state.sourceSelected(0, state.sourceList.itemRefs.NItems())
			})
		}()

	})
	return state
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

func (w *windowState) loadProfiles() {
	profileMap := defaults.GetDefaultTranscodeOptions()
	for name, profile := range profileMap {
		profile := profile
		w.profileList.add(name, &profile)
	}
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
			w.safePushStatus(fmt.Sprintf("Unable to add file to DB: %s", err))
			return
		}
	}
	w.addSourceAsset(file, sa)
}

func (w *windowState) addSourceAsset(file string, sa *dao.SourceAsset) {

	w.safePushStatus(fmt.Sprintf("Opened file: %s (%d)", file, sa.Id))
	go func() {
		ffProbe, err := transcoder.GetVideoMetadata(file)
		if err != nil {
			w.safePushStatus(fmt.Sprintf("File %s: %s", file, err))
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
		glib.IdleAdd(func() {
			w.sourceList.add(sa)
		})

	}()
}

func (w *windowState) sourceSelected(position, nItems uint) {
	if nItems == 0 {
		return
	}

	glib.IdleAdd(func() {
		selected := w.sourceList.itemRefs.String(w.sourceList.itemSelection.Selected())
		item := w.sourceList.itemDetails[selected]

		labelFilename := w.builder.GetObject("labelFilename").Cast().(*gtk.Label)
		labelFilename.SetLabel(item.asset.Filename)
		labelCodec := w.builder.GetObject("labelCodec").Cast().(*gtk.Label)
		labelCodec.SetLabel(item.asset.Codec)
		labelDuration := w.builder.GetObject("labelDuration").Cast().(*gtk.Label)
		labelDuration.SetLabel(fmt.Sprintf("%.02f seconds", item.asset.DurationTime))
		resolution := fmt.Sprintf("%dx%d", item.asset.Resolution.Width, item.asset.Resolution.Height)
		labelResolution := w.builder.GetObject("labelResolution").Cast().(*gtk.Label)
		labelResolution.SetLabel(resolution)

		labelFilesize := w.builder.GetObject("labelSize").Cast().(*gtk.Label)
		fileInfo, err := os.Stat(item.asset.Filename)
		if err != nil {
			labelFilesize.SetLabel("")
		} else {
			labelFilesize.SetLabel(fmt.Sprintf("%d", fileInfo.Size()))
		}
		outputs, err := w.daoInstance.GetTranscodeOutputs("WHERE source = ? ORDER BY profile_name, filename", item.asset.Id)
		if err != nil {
			w.status.Push(0, fmt.Sprintf("Unable to get transcode output list: %s", err))
		} else {
			w.outputList.updateOutputList(outputs)
		}
	})
}

func (w *windowState) safePushStatus(line string) {
	glib.IdleAdd(func() {
		w.status.Push(0, line)
	})

}

func (w *windowState) transcodeFile() {
	profile := w.profileList.getSelected()
	source := w.sourceList.getSelectedItem()

	msg := fmt.Sprintf("Transcode %s with profile %s", source.asset.Filename, profile.name)
	w.status.Push(0, msg)

	to, err := newTranscodeOutputWindow(source.asset, w.daoInstance)
	if err != nil {
		w.status.Push(0, fmt.Sprintf("Unable to create transcode output window: %s", err))
	} else {
		to.window.Show()
		go func() {
			time.Sleep(500 * time.Millisecond)
			to.startTranscode(profile.name, profile.profile)
		}()
	}
}

func (w *windowState) packageOutput() {
	selection := w.outputList.getSelected()
	files := []string{}
	for _, item := range selection {
		files = append(files, item.Filename)
	}
	w.safePushStatus(fmt.Sprintf("Package files: %s", strings.Join(files, ", ")))
}

func (w *windowState) deleteOutput() {
	selection := w.outputList.getSelected()
	for _, item := range selection {
		w.outputList.remove(item.Filename)
		_ = os.Remove(item.Filename)
		err := w.daoInstance.DeleteTranscodeOutput(item.Id)
		if err != nil {
			log.Errorf("Unable to delete transcode outputs: %s", err)
		}
	}
}

func RunUI(db *sql.DB) {
	state := newWindowState(db)

	if code := state.app.Run(os.Args); code > 0 {
		os.Exit(code)
	}
}
