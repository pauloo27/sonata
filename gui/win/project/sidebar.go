package project

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/pauloo27/sonata/common/data"
	"github.com/pauloo27/sonata/gui/utils"
)

func newSidebar(project *data.Project, selectedProject chan *data.Request) *gtk.Box {
	container, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	utils.HandleErr(err)

	container.SetHExpand(true)

	list, err := gtk.ListBoxNew()
	utils.HandleErr(err)

	for _, request := range project.ListRequests() {
		row, err := gtk.ListBoxRowNew()
		utils.HandleErr(err)

		label, err := gtk.LabelNew(request.Name)
		utils.HandleErr(err)

		row.Add(label)
		list.Add(row)
	}

	list.Connect("row-selected", func() {
		selectedProject <- project.ListRequests()[list.GetSelectedRow().GetIndex()]
	})

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
