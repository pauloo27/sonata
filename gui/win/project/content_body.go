package project

import (
	"encoding/json"

	"github.com/gotk3/gotk3/gtk"
	"github.com/gotk3/sourceview"
	"github.com/pauloo27/sonata/common/data"
	"github.com/pauloo27/sonata/gui/utils"
)

type BodyEditor interface {
	Widget(*ProjectStore) gtk.Widget
	Update(*ProjectStore)
	Destroy()
}

type FormBodyEditor struct {
	editor *utils.KeyValueEditor
}

func (fbe *FormBodyEditor) Widget(store *ProjectStore) gtk.Widget {
	initialValues := make(map[string]string)

	_ = json.Unmarshal([]byte(store.DraftRequest.Body), &initialValues)

	fbe.editor = utils.NewKeyValueEditor(initialValues)

	fbe.editor.SetVExpand(true)
	fbe.editor.OnAfterUpdate(func() {
		body, err := json.Marshal(fbe.editor.Store.ToMap())
		utils.HandleErr(err)

		store.DraftRequest.Body = string(body)
	})

	return fbe.editor.Widget
}

func (fbe *FormBodyEditor) Update(store *ProjectStore) {
	// noop, at least, for now
}

func (fbe *FormBodyEditor) Destroy() {
	if fbe.editor == nil {
		return
	}
	fbe.editor.Destroy()
	fbe.editor = nil
}

type TextBodyEditor struct {
	editor *utils.Editor
}

func (tbe *TextBodyEditor) Destroy() {
	if tbe.editor == nil {
		return
	}
	tbe.editor.Destroy()
	tbe.editor = nil
}

func (tbe *TextBodyEditor) Update(store *ProjectStore) {
	ext := data.BodyTypeExtensions[store.DraftRequest.BodyType]
	lang, _ := utils.Must(sourceview.SourceLanguageManagerGetDefault()).
		GetLanguage(ext)

	tbe.editor.Buffer.SetLanguage(lang)
}

func (tbe *TextBodyEditor) Widget(store *ProjectStore) gtk.Widget {
	ext := data.BodyTypeExtensions[store.DraftRequest.BodyType]

	tbe.editor = utils.NewEditor(
		store.DraftRequest.Body,
		true,
		ext,
	)

	tbe.editor.Buffer.Connect("changed", func() {
		store.DraftRequest.Body = utils.Must(
			tbe.editor.Buffer.GetText(
				tbe.editor.Buffer.GetStartIter(),
				tbe.editor.Buffer.GetEndIter(),
				true,
			),
		)
	})

	return tbe.editor.Widget
}

var (
	textBodyEditor = &TextBodyEditor{}
	formBodyEditor = &FormBodyEditor{}

	editorByBodyType = map[data.BodyType]BodyEditor{
		data.BodyTypeJSON: textBodyEditor,
		data.BodyTypeText: textBodyEditor,
		data.BodyTypeForm: formBodyEditor,
	}
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

	editor := editorByBodyType[store.DraftRequest.BodyType]

	bodyTypeEntry.SetActive(selectedBodyTypeIdx)
	bodyTypeEntry.Connect("changed", func() {
		newBodyType := data.BodyType(bodyTypeEntry.GetActiveText())
		store.DraftRequest.BodyType = newBodyType

		newEditor := editorByBodyType[newBodyType]

		if newEditor == editor {
			if editor != nil {
				editor.Update(store)
			}
			return
		}

		if editor != nil {
			editor.Destroy()
		}

		if newEditor != nil {
			editorWidget := newEditor.Widget(store)
			container.Add(&editorWidget)
			container.ShowAll()
		}
		editor = newEditor
	})

	bodyTypeContainer := utils.Must(gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 5))
	bodyTypeContainer.Add(bodyTypeLbl)
	bodyTypeContainer.Add(bodyTypeEntry)

	container.Add(bodyTypeContainer)
	if editor != nil {
		editorWidget := editor.Widget(store)
		container.Add(&editorWidget)
	}

	return container
}
