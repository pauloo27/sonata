package utils

import "github.com/gotk3/gotk3/gtk"

func Scrollable[T gtk.IWidget](child T) *gtk.ScrolledWindow {
	bodyScroll, err := gtk.ScrolledWindowNew(nil, nil)
	HandleErr(err)

	bodyScroll.SetVExpand(true)
	bodyScroll.Add(child)

	return bodyScroll
}
