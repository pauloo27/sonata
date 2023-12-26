package utils

import "github.com/gotk3/gotk3/gtk"

func AddCSSClass(widget gtk.Widget, cssClass string) {
	context, err := widget.GetStyleContext()
	HandleErr(err)
	context.AddClass(cssClass)
}
