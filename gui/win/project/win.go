package project

import (
	"fmt"
	"os"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/joho/godotenv"
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

	contentWrapperContainer := newContentWrapperContainer(store)

	container.Add1(newSidebar(store))
	container.Add2(contentWrapperContainer)

	contentContainer := newEmptyContentContainer()
	contentWrapperContainer.Add(contentContainer)

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
					contentWrapperContainer.Add(contentContainer)
					container.ShowAll()
				})
			}
		}
	}()

	gtkWin.Add(container)
	gtkWin.ShowAll()

	return gtkWin
}

func newContentWrapperContainer(store *ProjectStore) *gtk.Box {
	container := utils.Must(gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 5))
	container.SetMarginTop(5)

	topBar := utils.Must(gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 5))

	envCombo := utils.Must(gtk.ComboBoxTextNew())
	envCombo.AppendText("None")
	envCombo.SetActive(0)

	envs := utils.Must(store.Project.ListEnvironments())

	for _, env := range envs {
		envCombo.AppendText(env)
	}

	envCombo.Connect("changed", func() {
		name := envCombo.GetActiveText()
		variables, err := parseEnv(fmt.Sprintf("%s/%s", store.Project.RootDir, name))
		if err != nil {
			utils.ShowErrorDialog(store.Window, "Cannot load environment")
		}

		client.GetEnv = func(name string) string {
			return variables[name]
		}
	})

	topBar.PackEnd(
		envCombo,
		false, false, 5,
	)
	topBar.PackEnd(
		utils.Must(gtk.LabelNew("Environment:")),
		false, false, 5,
	)

	container.Add(topBar)

	return container
}

func parseEnv(path string) (map[string]string, error) {
	file, err := os.OpenFile(path, os.O_RDONLY, 0644)
	if err != nil {
		fmt.Println(path, err)
		return nil, err
	}
	return godotenv.Parse(file)
}
