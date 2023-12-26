package main

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/pauloo27/sonata/gui/win/welcome"
)

func main() {
	gtk.Init(nil)

	loadCSS()

	welcome.Start()
	gtk.Main()
}
