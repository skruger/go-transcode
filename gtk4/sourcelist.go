package gtk4

import (
	"fmt"
	"github.com/diamondburned/gotk4/pkg/gtk/v4"
	"github.com/skruger/privatestudio/dao"
)

type sourceListItem struct {
	asset *dao.SourceAsset
}

type sourceListData struct {
	listView      *gtk.ListView
	itemRefs      *gtk.StringList
	itemDetails   map[string]*sourceListItem
	factory       *gtk.SignalListItemFactory
	itemSelection *gtk.SingleSelection
}

func newSourceList(view *gtk.ListView, selectionChanged func(uint, uint)) *sourceListData {
	stringList := gtk.NewStringList([]string{})
	stringSelection := gtk.NewSingleSelection(stringList)
	listData := &sourceListData{
		listView:      view,
		itemRefs:      stringList,
		itemDetails:   make(map[string]*sourceListItem),
		factory:       gtk.NewSignalListItemFactory(),
		itemSelection: stringSelection,
	}
	stringSelection.ConnectSelectionChanged(selectionChanged)
	view.SetModel(stringSelection)
	listData.factory.ConnectBind(listData.bind)
	listData.factory.ConnectSetup(listData.setup)

	view.SetFactory(&listData.factory.ListItemFactory)
	return listData
}

func (s *sourceListData) setup(listItem *gtk.ListItem) {
	listItem.SetChild(gtk.NewLabel(""))
}

func (s *sourceListData) bind(listItem *gtk.ListItem) {
	idx := listItem.Position()
	key := s.itemRefs.String(idx)
	sa := s.itemDetails[key].asset
	labelInfo := fmt.Sprintf("%s (%s, %dx%d)", sa.Filename, sa.Codec, sa.Resolution.Width, sa.Resolution.Height)
	label := listItem.Child().(*gtk.Label)
	label.SetLabel(labelInfo)
}

func (s *sourceListData) add(asset *dao.SourceAsset) {
	s.itemDetails[asset.Filename] = &sourceListItem{asset: asset}
	for idx := 0; idx < int(s.itemRefs.NItems()); idx++ {
		if s.itemRefs.String(uint(idx)) == asset.Filename {
			return
		}
	}
	s.itemRefs.Append(asset.Filename)
}

func (s *sourceListData) getSelectedItem() *sourceListItem {
	selected := s.itemRefs.String(s.itemSelection.Selected())
	return s.itemDetails[selected]
}
