package gtk4

import (
	"fmt"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/skruger/privatestudio/dao"
)

// transcodeOutput

type outputListData struct {
	listView      *gtk.ListView
	itemRefs      *gtk.StringList
	itemDetails   map[string]dao.TranscodeOutput
	factory       *gtk.SignalListItemFactory
	itemSelection *gtk.MultiSelection
}

func newOutputList(view *gtk.ListView) *outputListData {
	stringList := gtk.NewStringList([]string{})
	stringSelection := gtk.NewMultiSelection(stringList)
	listData := &outputListData{
		listView:      view,
		itemRefs:      stringList,
		itemDetails:   make(map[string]dao.TranscodeOutput),
		factory:       gtk.NewSignalListItemFactory(),
		itemSelection: stringSelection,
	}
	view.SetModel(stringSelection)
	listData.factory.ConnectBind(listData.bind)
	listData.factory.ConnectSetup(listData.setup)

	view.SetFactory(&listData.factory.ListItemFactory)
	return listData
}

func (o *outputListData) setup(listItem *gtk.ListItem) {
	listItem.SetChild(gtk.NewLabel(""))
}

func (o *outputListData) bind(listItem *gtk.ListItem) {
	idx := listItem.Position()
	key := o.itemRefs.String(idx)
	output := o.itemDetails[key]
	labelInfo := fmt.Sprintf("%s - %s (%dx%d)", *output.ProfileName, output.Filename, output.Resolution.Width, output.Resolution.Height)
	label := listItem.Child().(*gtk.Label)
	label.SetLabel(labelInfo)
}

func (o *outputListData) updateOutputList(outputs []dao.TranscodeOutput) {
	for o.itemRefs.NItems() > 0 {
		key := o.itemRefs.String(0)
		delete(o.itemDetails, key)
		o.itemRefs.Remove(0)
	}

	for _, output := range outputs {
		o.itemDetails[output.Filename] = output
		o.itemRefs.Append(output.Filename)
	}
}

func (o *outputListData) getSelected() (outputs []dao.TranscodeOutput) {
	var i uint
	for i = 0; i < o.itemSelection.NItems(); i++ {
		if o.itemSelection.IsSelected(i) {
			key := o.itemRefs.String(i)
			item := o.itemDetails[key]
			outputs = append(outputs, item)
		}
	}
	return
}

func (o *outputListData) remove(outputFilename string) {
	var i uint
	for i = 0; i < o.itemRefs.NItems(); i++ {
		if o.itemRefs.String(i) == outputFilename {
			o.itemRefs.Remove(i)
			delete(o.itemDetails, outputFilename)
			return
		}
	}
}
