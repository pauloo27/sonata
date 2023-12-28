package utils

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/gotk3/sourceview"
)

type Editor struct {
	*gtk.ScrolledWindow

	View   *sourceview.SourceView
	Buffer *sourceview.SourceBuffer
}

func NewEditor(initialText string, editable bool, lang string) *Editor {
	editorBuf := Must(
		sourceview.SourceBufferNewWithLanguage(
			Must(Must(sourceview.SourceLanguageManagerGetDefault()).GetLanguage(lang)),
		),
	)
	editorBuf.SetText(initialText)

	editorView := Must(sourceview.SourceViewNewWithBuffer(editorBuf))
	editorView.SetEditable(editable)
	editorView.SetHExpand(true)
	editorView.SetShowLineNumbers(true)
	editorView.SetMonospace(true)

	editor := &Editor{
		ScrolledWindow: Scrollable(editorView),
		View:           editorView,
		Buffer:         editorBuf,
	}

	return editor
}
