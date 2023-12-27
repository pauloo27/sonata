package utils

import "github.com/gotk3/gotk3/gtk"

type Editor struct {
	*gtk.ScrolledWindow

	TextView *gtk.TextView
	Buffer   *gtk.TextBuffer
}

func NewEditor(initialText string, editable bool) *Editor {
	bodyBuf := Must(gtk.TextBufferNew(nil))
	bodyBuf.SetText(initialText)

	bodyView := Must(gtk.TextViewNewWithBuffer(bodyBuf))
	bodyView.SetEditable(editable)
	bodyView.SetHExpand(true)

	editor := &Editor{
		ScrolledWindow: Scrollable(bodyView),
		TextView:       bodyView,
		Buffer:         bodyBuf,
	}

	return editor
}
