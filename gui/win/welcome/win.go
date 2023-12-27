package welcome

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/pauloo27/sonata/gui/utils"
)

var (
	swappingWindow bool
)

func Start() {
	win := utils.Must(gtk.WindowNew(gtk.WINDOW_TOPLEVEL))

	win.SetTitle("Sonata")
	win.SetDefaultSize(800, 600)

	_ = win.Connect("destroy", func() {
		if swappingWindow {
			return
		}
		gtk.MainQuit()
	})

	win.Add(newContentContainer(win))

	win.ShowAll()
}
