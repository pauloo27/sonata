package welcome

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/pauloo27/sonata/gui/utils"
	"github.com/pauloo27/sonata/gui/win"
)

func init() {
	win.AddWindow("welcome", &win.SonataWindow{
		Start: Start,
	})
}

func Start(...interface{}) *gtk.Window {
	gtkWin := utils.Must(gtk.WindowNew(gtk.WINDOW_TOPLEVEL))

	gtkWin.SetTitle("Sonata")
	gtkWin.SetDefaultSize(800, 600)

	_ = gtkWin.Connect("destroy", win.HandleClose)

	gtkWin.Add(newContentContainer(gtkWin))

	gtkWin.ShowAll()
	return gtkWin
}
