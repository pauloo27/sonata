package main

import (
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
	"github.com/pauloo27/sonata/gui/utils"
)

var cssData = `
		.welcome-title {
			font-size: 40px;
		}
		.welcome-subtitle {
			font-size: 32px;
		}
`

func loadCSS() {
	cssProvider, err := gtk.CssProviderNew()
	utils.HandleErr(err)
	cssProvider.LoadFromData(cssData)
	screen, err := gdk.ScreenGetDefault()
	utils.HandleErr(err)
	gtk.AddProviderForScreen(screen, cssProvider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)
}
