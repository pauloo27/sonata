package project

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/gotk3/sourceview"
	"github.com/pauloo27/sonata/common/data"
	"github.com/pauloo27/sonata/gui/utils"
)

func newBodyContainer(store *ProjectStore) *gtk.Box {
	container := utils.Must(gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0))

	bodyTypeLbl := utils.Must(gtk.LabelNew("Body type:"))

	selectedBodyTypeIdx := 0

	bodyTypeEntry := utils.Must(gtk.ComboBoxTextNew())
	for i, bodyType := range data.BodyTypes {
		bodyTypeEntry.AppendText(string(bodyType))
		if bodyType == store.DraftRequest.BodyType {
			selectedBodyTypeIdx = i
		}
	}

	var editor *utils.Editor

	handleEditorChange := func() {
		store.DraftRequest.Body = utils.Must(
			editor.Buffer.GetText(editor.Buffer.GetStartIter(), editor.Buffer.GetEndIter(), true),
		)
	}

	if store.DraftRequest.BodyType != data.BodyTypeNone {
		editor = utils.NewEditor(
			store.DraftRequest.Body,
			true,
			data.BodyTypeExtensions[store.DraftRequest.BodyType],
		)
		editor.Buffer.Connect("changed", handleEditorChange)
	}

	bodyTypeEntry.SetActive(selectedBodyTypeIdx)
	bodyTypeEntry.Connect("changed", func() {
		newBodyType := data.BodyType(bodyTypeEntry.GetActiveText())
		store.DraftRequest.BodyType = newBodyType

		if newBodyType == data.BodyTypeNone {
			editor.Destroy()
			editor = nil
			return
		}
		if editor == nil {
			editor = utils.NewEditor(
				store.DraftRequest.Body,
				true,
				data.BodyTypeExtensions[store.DraftRequest.BodyType],
			)
			container.Add(editor)
			editor.Buffer.Connect("changed", handleEditorChange)
			editor.ShowAll()
		}

		ext := data.BodyTypeExtensions[store.DraftRequest.BodyType]
		lang, _ := utils.Must(sourceview.SourceLanguageManagerGetDefault()).
			GetLanguage(ext)

		editor.Buffer.SetLanguage(lang)
	})

	bodyTypeContainer := utils.Must(gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 5))
	bodyTypeContainer.Add(bodyTypeLbl)
	bodyTypeContainer.Add(bodyTypeEntry)

	container.Add(bodyTypeContainer)
	if editor != nil {
		container.Add(editor)
	}

	return container
}
