package project

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/pauloo27/sonata/common/data"
	"github.com/pauloo27/sonata/gui/utils"
)

func newSidebar(project *data.Project, selectedProject chan *data.Request) *gtk.Box {
	container := utils.Must(gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0))
	container.SetHExpand(true)

	list := utils.Must(gtk.ListBoxNew())

	for _, request := range project.ListRequests() {
		row := utils.Must(gtk.ListBoxRowNew())
		row.Add(
			utils.Must(gtk.LabelNew(request.Name)),
		)
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
	container := utils.Must(gtk.HeaderBarNew())
	container.SetTitle("Sonata")

	settingsBtn := utils.Must(gtk.ButtonNewFromIconName("open-menu-symbolic", gtk.ICON_SIZE_BUTTON))
	container.PackStart(settingsBtn)

	return container
}
