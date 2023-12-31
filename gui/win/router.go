package win

import (
	"fmt"
	"sync"

	"github.com/gotk3/gotk3/gtk"
)

type SonataWindow struct {
	Start func(...any) *gtk.Window
}

var (
	routes           = make(map[string]*SonataWindow)
	lastWindow       *gtk.Window
	lock             = sync.Mutex{}
	replacingWindows = false
)

func HandleClose() {
	if replacingWindows {
		return
	}

	gtk.MainQuit()
}

func AddWindow(name string, win *SonataWindow) {
	routes[name] = win
}

func Replace(name string, params ...any) {
	lock.Lock()
	defer lock.Unlock()

	replacingWindows = true

	win, found := routes[name]
	if !found {
		panic(fmt.Sprintf("Window %s not found", name))
	}

	if lastWindow != nil {
		lastWindow.Destroy()
	}

	lastWindow = win.Start(params...)
	replacingWindows = false
}

func ShowWindow(name string, params ...any) {
	win, found := routes[name]
	if !found {
		panic(fmt.Sprintf("Window %s not found", name))
	}

	win.Start(params...)
}
