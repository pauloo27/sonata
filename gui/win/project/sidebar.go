package project

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/pauloo27/sonata/gui/utils"
)

var (
	sampleList = []string{
		"Get user",
		"Get users",
		"Create user",
		"Update user",
		"Delete user",
	}
)

func newSidebar() *gtk.Box {
	container, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	utils.HandleErr(err)

	container.SetHExpand(true)

	list, err := gtk.ListBoxNew()
	utils.HandleErr(err)

	for _, item := range sampleList {
		row, err := gtk.ListBoxRowNew()
		utils.HandleErr(err)

		label, err := gtk.LabelNew(item)
		utils.HandleErr(err)

		row.Add(label)
		list.Add(row)
	}

	container.Add(newSidebarHeader())
	container.Add(list)

	return container
}

func newSidebarHeader() *gtk.HeaderBar {
	container, err := gtk.HeaderBarNew()
	utils.HandleErr(err)

	container.SetTitle("Sonata")

	settingsBtn, err := gtk.ButtonNewFromIconName("open-menu-symbolic", gtk.ICON_SIZE_BUTTON)
	utils.HandleErr(err)

	container.PackStart(settingsBtn)

	return container
}
