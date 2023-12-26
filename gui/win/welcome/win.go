package welcome

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/pauloo27/sonata/gui/utils"
)

var (
	swappingWindow bool
)

func Start() {
	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	utils.HandleErr(err)

	win.SetTitle("Sonata")
	win.SetDefaultSize(800, 600)

	_ = win.Connect("destroy", func() {
		if swappingWindow {
			return
		}
		gtk.MainQuit()
	})
	utils.HandleErr(err)

	win.Add(newContentContainer(win))

	win.ShowAll()
}
