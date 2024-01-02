package env

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/pauloo27/sonata/common/data"
	"github.com/pauloo27/sonata/gui/utils"
	"github.com/pauloo27/sonata/gui/win"
)

func init() {
	win.AddWindow("env", &win.SonataWindow{
		Start: Start,
	})
}

type EnvStore struct {
	Window  *gtk.Window
	Project *data.Project
	EnvName string
	Envs    map[string]string
}

func Start(params ...any) *gtk.Window {
	gtkWin := utils.Must(gtk.WindowNew(gtk.WINDOW_TOPLEVEL))

	gtkWin.SetTitle("Sonata")
	gtkWin.SetDefaultSize(800, 600)

	store := &EnvStore{
		Window:  gtkWin,
		Project: params[0].(*data.Project),
		EnvName: params[1].(string),
		Envs:    params[2].(map[string]string),
	}

	gtkWin.Add(newContentContainer(store))

	gtkWin.ShowAll()
	return gtkWin
}
