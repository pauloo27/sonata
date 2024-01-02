package main

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/pauloo27/sonata/gui/win"

	_ "github.com/pauloo27/sonata/gui/win/env"
	_ "github.com/pauloo27/sonata/gui/win/project"
	_ "github.com/pauloo27/sonata/gui/win/welcome"
)

func main() {
	gtk.Init(nil)

	loadCSS()

	win.Replace("welcome")
	gtk.Main()
}
