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

func NewEditor(initialText string, editable bool, extention string) *Editor {
	lang, _ := Must(sourceview.SourceLanguageManagerGetDefault()).GetLanguage(extention)

	var editorBuf *sourceview.SourceBuffer
	if lang == nil {
		// for some reason, if i dont set the lang at the start it will crash
		// hacky solution -- works tho
		editorBuf = Must(
			sourceview.SourceBufferNewWithLanguage(
				Must(Must(sourceview.SourceLanguageManagerGetDefault()).GetLanguage("json")),
			),
		)
		editorBuf.SetLanguage(nil)
	} else {
		editorBuf = Must(sourceview.SourceBufferNewWithLanguage(lang))
	}

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
