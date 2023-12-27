package project

import (
	"fmt"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/pauloo27/sonata/common/data"
	"github.com/pauloo27/sonata/gui/utils"
)

func Start(path string) {
	win := utils.Must(gtk.WindowNew(gtk.WINDOW_TOPLEVEL))
	win.SetTitle("Sonata")
	win.SetDefaultSize(800, 600)

	_ = win.Connect("destroy", func() {
		fmt.Println("Closed")
		gtk.MainQuit()
	})

	// FIXME: proper error handler
	project, err := data.LoadProject(path)
	utils.HandleErr(err)

	container := utils.Must(gtk.PanedNew(gtk.ORIENTATION_HORIZONTAL))
	container.SetPosition(200)

	selectedRequest := make(chan *data.Request)

	contentContainer := newEmptyContentContainer()

	container.Add1(newSidebar(project, selectedRequest))
	container.Add2(contentContainer)

	go func() {
		for {
			request := <-selectedRequest
			if request != nil {
				glib.IdleAdd(func() {
					container.Remove(contentContainer)
					contentContainer.Destroy()
					contentContainer = newContentContainer(request)
					container.Add2(contentContainer)
					container.ShowAll()
				})
			}
		}
	}()

	win.Add(container)
	win.ShowAll()
}
