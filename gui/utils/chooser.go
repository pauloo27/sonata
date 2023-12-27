package utils

import (
	"github.com/gotk3/gotk3/gtk"
)

func ChooseFolder(parent *gtk.Window) string {

	dialog := Must(gtk.FileChooserDialogNewWith2Buttons(
		"Open project",
		parent,
		gtk.FILE_CHOOSER_ACTION_OPEN,
		"Open",
		gtk.RESPONSE_ACCEPT,
		"Cancel",
		gtk.RESPONSE_CANCEL,
	))
	dialog.SetAction(gtk.FILE_CHOOSER_ACTION_SELECT_FOLDER)

	response := dialog.Run()
	defer dialog.Destroy()

	if response == gtk.RESPONSE_ACCEPT {
		return dialog.GetFilename()
	}

	return ""
}
