package utils

import "github.com/gotk3/gotk3/gtk"

func ClearChildren(parent gtk.Container) {
	parent.GetChildren().Foreach(func(item any) {
		parent.Remove(item.(*gtk.Widget))
		item.(*gtk.Widget).Destroy()
	})
}
