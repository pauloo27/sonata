package project

import (
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/pauloo27/sonata/common/client"
	"github.com/pauloo27/sonata/common/data"
	"github.com/pauloo27/sonata/gui/utils"
	"github.com/pauloo27/sonata/gui/win"
)

func init() {
	win.AddWindow("project", &win.SonataWindow{
		Start: Start,
	})
}

type ProjectStore struct {
	Project       *data.Project
	SavedRequest  *data.Request
	DraftRequest  *data.Request
	VarStore      *VariablesStore
	RequestCh     chan *data.Request
	ResponseCh    chan *client.Response
	ReloadSidebar func(selectedRequest *data.Request)
	Window        *gtk.Window
}

func Start(params ...interface{}) *gtk.Window {
	gtkWin := utils.Must(gtk.WindowNew(gtk.WINDOW_TOPLEVEL))
	gtkWin.SetTitle("Sonata")
	gtkWin.SetDefaultSize(800, 600)

	project := params[0].(*data.Project)

	container := utils.Must(gtk.PanedNew(gtk.ORIENTATION_HORIZONTAL))
	container.SetPosition(300)

	store := &ProjectStore{
		Window:    gtkWin,
		Project:   project,
		VarStore:  newVariablesStore(),
		RequestCh: make(chan *data.Request, 2),
	}

	_ = gtkWin.Connect("destroy", win.HandleClose)

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

	gtkWin.Add(container)
	gtkWin.ShowAll()

	return gtkWin
}
