package project

import (
	"fmt"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/pauloo27/sonata/common/client"
	"github.com/pauloo27/sonata/common/data"
	"github.com/pauloo27/sonata/gui/utils"
)

func Start(path string) bool {
	win := utils.Must(gtk.WindowNew(gtk.WINDOW_TOPLEVEL))
	win.SetTitle("Sonata")
	win.SetDefaultSize(800, 600)

	project, err := data.LoadProject(path)
	if err != nil {
		utils.ShowErrorDialog(win, fmt.Sprintf("Failed to load project: %s", path))
		return false
	}

	_ = win.Connect("destroy", func() {
		fmt.Println("Closed")
		gtk.MainQuit()
	})

	container := utils.Must(gtk.PanedNew(gtk.ORIENTATION_HORIZONTAL))
	container.SetPosition(200)

	store := &ContentStore{
		Project:    project,
		VarStore:   newVariablesStore(),
		ResponseCh: make(chan *client.Response, 2),
		RequestCh:  make(chan *data.Request, 2),
	}

	contentContainer := newEmptyContentContainer()

	container.Add1(newSidebar(store))
	container.Add2(contentContainer)

	go func() {
		for {
			request := <-store.RequestCh
			draftRequest := request.Clone()

			store.SavedRequest = request
			store.DraftRequest = draftRequest

			if request != nil {
				glib.IdleAdd(func() {
					container.Remove(contentContainer)
					contentContainer.Destroy()
					contentContainer = newContentContainer(store)
					container.Add2(contentContainer)
					container.ShowAll()
				})
			}
		}
	}()

	win.Add(container)
	win.ShowAll()
	return true
}
