package project

import (
	"fmt"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/pauloo27/sonata/common/client"
	"github.com/pauloo27/sonata/common/data"
	"github.com/pauloo27/sonata/gui/utils"
)

type ProjectStore struct {
	Project       *data.Project
	SavedRequest  *data.Request
	DraftRequest  *data.Request
	VarStore      *VariablesStore
	RequestCh     chan *data.Request
	ResponseCh    chan *client.Response
	ReloadSidebar func()
}

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
	container.SetPosition(300)

	store := &ProjectStore{
		Project:   project,
		VarStore:  newVariablesStore(),
		RequestCh: make(chan *data.Request, 2),
	}

	contentContainer := newEmptyContentContainer()

	container.Add1(newSidebar(store))
	container.Add2(contentContainer)

	go func() {
		for {
			request := <-store.RequestCh

			glib.IdleAdd(func() {
				container.Remove(contentContainer)
				contentContainer.Destroy()
			})

			if request == nil {
				glib.IdleAdd(func() {
					contentContainer = newEmptyContentContainer()
					container.Add2(contentContainer)
					container.ShowAll()
				})
			} else {
				draftRequest := request.Clone()

				store.SavedRequest = request
				store.DraftRequest = draftRequest

				glib.IdleAdd(func() {
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
