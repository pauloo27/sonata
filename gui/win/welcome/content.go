package welcome

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/pauloo27/sonata/gui/utils"
	"github.com/pauloo27/sonata/gui/win/project"
)

func newContentContainer(win *gtk.Window) *gtk.Box {
	contentContainer := utils.Must(gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0))
	contentContainer.SetVExpand(true)
	contentContainer.SetHAlign(gtk.ALIGN_CENTER)
	contentContainer.SetVAlign(gtk.ALIGN_CENTER)

	welcomeLbl := utils.Must(gtk.LabelNew("Welcome to"))
	utils.AddCSSClass(welcomeLbl.Widget, "welcome-subtitle")

	sonataLbl := utils.Must(gtk.LabelNew("💤 SONATA 💤"))
	utils.AddCSSClass(sonataLbl.Widget, "welcome-title")

	createNewBtn := utils.Must(gtk.ButtonNewWithLabel("Create new project"))

	openBtn := utils.Must(gtk.ButtonNewWithLabel("Open project"))
	openBtn.Connect("clicked", func() {
		swappingWindow = true
		project.Start("/home/paulo/dev/sonata")
		win.Close()
	})

	infoLabel := utils.Must(gtk.LabelNew("(eventually it will list recent projects here)"))

	contentContainer.PackStart(welcomeLbl, false, false, 5)
	contentContainer.PackStart(sonataLbl, false, false, 5)
	contentContainer.PackStart(createNewBtn, false, false, 5)
	contentContainer.PackStart(openBtn, false, false, 5)
	contentContainer.PackStart(infoLabel, false, false, 0)

	return contentContainer
}
