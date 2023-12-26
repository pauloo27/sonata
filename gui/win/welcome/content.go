package welcome

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/pauloo27/sonata/gui/utils"
	"github.com/pauloo27/sonata/gui/win/project"
)

func newContentContainer(win *gtk.Window) *gtk.Box {
	contentContainer, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	utils.HandleErr(err)

	contentContainer.SetVExpand(true)
	contentContainer.SetHAlign(gtk.ALIGN_CENTER)
	contentContainer.SetVAlign(gtk.ALIGN_CENTER)

	welcomeLbl, err := gtk.LabelNew("Welcome to")
	utils.HandleErr(err)

	utils.AddCSSClass(welcomeLbl.Widget, "welcome-subtitle")

	sonataLbl, err := gtk.LabelNew("💤 SONATA 💤")
	utils.HandleErr(err)

	utils.AddCSSClass(sonataLbl.Widget, "welcome-title")

	createNewBtn, err := gtk.ButtonNewWithLabel("Create new project")
	utils.HandleErr(err)

	openBtn, err := gtk.ButtonNewWithLabel("Open project")
	utils.HandleErr(err)

	openBtn.Connect("clicked", func() {
		swappingWindow = true
		project.Start("/home/paulo/dev/sonata")
		win.Close()
	})

	infoLabel, err := gtk.LabelNew("(eventually it will list recent projects here)")
	utils.HandleErr(err)

	contentContainer.PackStart(welcomeLbl, false, false, 5)
	contentContainer.PackStart(sonataLbl, false, false, 5)
	contentContainer.PackStart(createNewBtn, false, false, 5)
	contentContainer.PackStart(openBtn, false, false, 5)
	contentContainer.PackStart(infoLabel, false, false, 0)

	return contentContainer
}
