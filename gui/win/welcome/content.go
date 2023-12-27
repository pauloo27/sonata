package welcome

import (
	"path"

	"github.com/gotk3/gotk3/gtk"
	"github.com/pauloo27/sonata/common/data"
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
	createNewBtn.Connect("clicked", func() {
		selectedPath := utils.ChooseFolder(win, "Select project folder")
		if selectedPath == "" {
			return
		}

		name := path.Base(selectedPath)

		p, err := data.NewProject(selectedPath, name)
		if err != nil {
			utils.ShowErrorDialog(win, "Failed to create project")
			return
		}

		if err = p.Save(); err != nil {
			utils.ShowErrorDialog(win, "Failed to save project")
			return
		}

		swappingWindow = true
		if ok := project.Start(selectedPath); ok {
			win.Destroy()
		}
		swappingWindow = false
	})

	openBtn := utils.Must(gtk.ButtonNewWithLabel("Open project"))
	openBtn.Connect("clicked", func() {
		selectedPath := utils.ChooseFolder(win, "Open project")
		if selectedPath == "" {
			return
		}

		swappingWindow = true
		if ok := project.Start(selectedPath); ok {
			win.Destroy()
		}
		swappingWindow = false
	})

	infoLabel := utils.Must(gtk.LabelNew("(eventually it will list recent projects here)"))

	contentContainer.PackStart(welcomeLbl, false, false, 5)
	contentContainer.PackStart(sonataLbl, false, false, 5)
	contentContainer.PackStart(createNewBtn, false, false, 5)
	contentContainer.PackStart(openBtn, false, false, 5)
	contentContainer.PackStart(infoLabel, false, false, 0)

	return contentContainer
}
