package gtk4

import (
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/skruger/privatestudio/transcoder/config"
)

type transcodeOptionItem struct {
	profile *config.TranscodeOptions
	name    string
}

type transcodeOptionData struct {
	dropDown    *gtk.DropDown
	itemRefs    *gtk.StringList
	itemDetails map[string]*transcodeOptionItem
}

func newTranscodeOptionList(dropDown *gtk.DropDown) *transcodeOptionData {
	stringList := gtk.NewStringList([]string{})
	//stringSelection := gtk.NewSingleSelection(stringList)
	dropDown.SetModel(stringList)
	listData := &transcodeOptionData{
		dropDown:    dropDown,
		itemRefs:    stringList,
		itemDetails: make(map[string]*transcodeOptionItem),
	}
	return listData
}

func (t *transcodeOptionData) add(name string, profile *config.TranscodeOptions) {
	t.itemDetails[name] = &transcodeOptionItem{
		profile: profile,
		name:    name,
	}
	t.itemRefs.Append(name)
}

func (t *transcodeOptionData) getSelected() *transcodeOptionItem {
	profileName := t.itemRefs.String(t.dropDown.Selected())
	return t.itemDetails[profileName]
}
