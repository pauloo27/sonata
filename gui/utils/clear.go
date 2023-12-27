package utils

import "github.com/gotk3/gotk3/gtk"

func ClearChildren(parent gtk.Container) {
	parent.GetChildren().Foreach(func(item interface{}) {
		item.(*gtk.Widget).Destroy()
	})
}
