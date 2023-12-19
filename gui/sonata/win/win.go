package win

import (
	"fmt"

	"github.com/gotk3/gotk3/gtk"
)

func Start() {
	gtk.Init(nil)
	win, err := gtk.WindowNew(gtk.WINDOW_TOPLEVEL)
	handleErr(err)

	win.SetTitle("Sonata")
	win.SetDefaultSize(800, 600)

	_ = win.Connect("destroy", func() {
		fmt.Println("Closed")
		gtk.MainQuit()
	})
	handleErr(err)

	container, err := gtk.PanedNew(gtk.ORIENTATION_HORIZONTAL)
	handleErr(err)

	container.SetPosition(200)

	container.Add1(newSidebar())
	container.Add2(newContentContainer())

	win.Add(container)

	win.ShowAll()
	gtk.Main()
}

func handleErr(err error) {
	if err != nil {
		panic(err)
	}
}
