package utils

import (
	"github.com/gotk3/gotk3/gtk"
)

func ShowEntryDialog(parent *gtk.Window, title, label string) string {
	dialog := gtk.MessageDialogNew(
		parent,
		gtk.DIALOG_MODAL,
		gtk.MESSAGE_QUESTION,
		gtk.BUTTONS_OK,
		title,
	)

	entry := Must(gtk.EntryNew())
	entryLbl := Must(gtk.LabelNew(label))

	rootContainer := Must(dialog.GetContentArea())

	entryContainer := Must(gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 5))

	entryContainer.PackStart(entryLbl, false, false, 3)
	entryContainer.PackStart(entry, false, false, 3)

	rootContainer.Add(entryContainer)
	rootContainer.ShowAll()

	entry.SetPlaceholderText("Env name")

	res := dialog.Run()
	value := Must(entry.GetText())

	dialog.Destroy()

	if res != gtk.RESPONSE_OK {
		return ""
	}

	return value
}
