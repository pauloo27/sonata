package utils

import "github.com/gotk3/gotk3/gtk"

func HandleErr(err error) {
	if err != nil {
		panic(err)
	}
}

func Must[T any](val T, err error) T {
	HandleErr(err)
	return val
}

func ShowErrorDialog(parent *gtk.Window, message string) {
	dialog := gtk.MessageDialogNew(
		parent,
		gtk.DIALOG_MODAL,
		gtk.MESSAGE_ERROR,
		gtk.BUTTONS_OK,
		message,
	)

	_ = dialog.Run()
	dialog.Destroy()
}
