package project

import (
	"fmt"

	"github.com/gotk3/gotk3/gtk"
	"github.com/pauloo27/sonata/common/data"
	"github.com/pauloo27/sonata/gui/utils"
)

func Start(path string) {
	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	utils.HandleErr(err)

	win.SetTitle("Sonata")
	win.SetDefaultSize(800, 600)

	_ = win.Connect("destroy", func() {
		fmt.Println("Closed")
		gtk.MainQuit()
	})
	utils.HandleErr(err)

	project, err := data.LoadProject(path)
	utils.HandleErr(err)

	container, err := gtk.PanedNew(gtk.ORIENTATION_HORIZONTAL)
	utils.HandleErr(err)

	container.SetPosition(200)

	container.Add1(newSidebar(project))
	container.Add2(newContentContainer())

	win.Add(container)

	win.ShowAll()
}
