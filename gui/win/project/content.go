package project

import (
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/pauloo27/sonata/common/client"
	"github.com/pauloo27/sonata/common/data"
	"github.com/pauloo27/sonata/gui/utils"
)

type ContentStore struct {
	Project      *data.Project
	SavedRequest *data.Request
	DraftRequest *data.Request
	VarStore     *VariablesStore
	ResponseCh   chan *client.Response
}

func newContentContainer(request *data.Request) *gtk.Box {
	container := utils.Must(gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 5))
	container.SetMarginTop(5)
	container.SetMarginBottom(5)
	container.SetMarginStart(5)
	container.SetMarginEnd(5)

	draftRequest := request.Clone()

	store := &ContentStore{
		SavedRequest: request,
		DraftRequest: draftRequest,
		VarStore:     newVariablesStore(),
		ResponseCh:   make(chan *client.Response, 2),
	}

	container.Add(newRequestURLContainer(store))

	subContainer := utils.Must(gtk.PanedNew(gtk.ORIENTATION_VERTICAL))
	subContainer.SetPosition(500)
	subContainer.Add1(newRequestStructureContainer(store))
	subContainer.Add2(newResponseContainer(store))

	container.Add(subContainer)
	container.SetHExpand(true)

	return container
}

func newRequestURLContainer(
	store *ContentStore,
) *gtk.Box {
	container := utils.Must(gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 5))

	methods := utils.Must(gtk.ComboBoxTextNew())

	requestMethodIdx := 0

	for i, method := range data.HTTPMethods {
		methods.AppendText(string(method))
		if method == store.DraftRequest.Method {
			requestMethodIdx = i
		}
	}

	methods.SetActive(requestMethodIdx)

	methods.Connect("changed", func() {
		store.DraftRequest.Method = data.HTTPMethod(methods.GetActiveText())
	})

	entry := utils.Must(gtk.EntryNew())
	entry.SetText(store.DraftRequest.URL)
	entry.SetHExpand(true)

	entry.Connect("changed", func() {
		store.DraftRequest.URL = utils.Must(entry.GetText())
	})

	sendBtn := utils.Must(gtk.ButtonNewWithLabel("Go"))
	sendBtn.Connect("clicked", func() {
		sendBtn.SetSensitive(false)

		go func() {
			client := client.NewClient()
			variables := make(map[string]string)

			for _, variable := range store.VarStore.List() {
				variables[variable.Key] = variable.Value
			}

			res, err := client.Run(store.DraftRequest, variables)
			// FIXME: proper error handling
			utils.HandleErr(err)

			glib.IdleAdd(func() {
				sendBtn.SetSensitive(true)
			})

			store.ResponseCh <- res
		}()
	})

	saveBtn := utils.Must(gtk.ButtonNewWithLabel("Save"))
	saveBtn.Connect("clicked", func() {
		*store.SavedRequest = *store.DraftRequest
		err := store.SavedRequest.Save()
		utils.HandleErr(err)
	})

	container.Add(methods)
	container.Add(entry)
	container.Add(saveBtn)
	container.Add(sendBtn)

	container.SetHExpand(true)

	return container
}

func newRequestStructureContainer(store *ContentStore) *gtk.Notebook {
	container := utils.Must(gtk.NotebookNew())
	headersContainer := utils.Must(gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0))

	bodyContainer := newBodyContainer(store)

	variablesContainer := newVariablesContainer(store)

	container.AppendPage(variablesContainer, utils.Must(gtk.LabelNew("Variables")))
	container.AppendPage(bodyContainer, utils.Must(gtk.LabelNew("Body")))
	container.AppendPage(headersContainer, utils.Must(gtk.LabelNew("Headers")))
	container.SetVExpand(true)

	return container
}

func newBodyContainer(store *ContentStore) *gtk.Box {
	container := utils.Must(gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0))
	editor := utils.NewEditor(store.DraftRequest.Body, true)

	editor.Buffer.Connect("changed", func() {
		store.DraftRequest.Body = utils.Must(
			editor.Buffer.GetText(editor.Buffer.GetStartIter(), editor.Buffer.GetEndIter(), true),
		)
	})

	container.Add(editor)

	return container
}

func newResponseContainer(
	store *ContentStore,
) *gtk.Box {
	container := utils.Must(gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0))

	notebook := utils.Must(gtk.NotebookNew())

	bodyContainer := utils.Must(gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0))

	headersContainer := utils.Must(gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0))

	newLbl := func(label string) *gtk.Label {
		lbl, err := gtk.LabelNew(label)
		utils.HandleErr(err)

		return lbl
	}

	bodyContainer.Add(newLbl("Body"))
	headersContainer.Add(newLbl("Headers"))

	notebook.AppendPage(bodyContainer, newLbl("Body"))
	notebook.AppendPage(headersContainer, newLbl("Headers"))
	notebook.SetVExpand(true)

	title := utils.Must(gtk.LabelNew("Response"))
	title.SetHAlign(gtk.ALIGN_CENTER)

	container.Add(title)
	container.Add(notebook)

	go func() {
		for {
			response := <-store.ResponseCh
			if response != nil {
				glib.IdleAdd(func() {
					notebook.SetCurrentPage(0)

					utils.ClearChildren(bodyContainer.Container)

					bodyContainer.Add(utils.NewEditor(response.Body, false))

					bodyContainer.ShowAll()
				})
			}
		}
	}()

	return container
}

func newEmptyContentContainer() *gtk.Box {
	container := utils.Must(gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0))
	container.SetHExpand(true)
	container.SetVExpand(true)
	container.SetHAlign(gtk.ALIGN_CENTER)
	container.SetVAlign(gtk.ALIGN_CENTER)

	title := utils.Must(gtk.LabelNew("Select a request to start"))
	title.SetHAlign(gtk.ALIGN_CENTER)
	title.SetVAlign(gtk.ALIGN_CENTER)

	utils.AddCSSClass(title.Widget, "welcome-subtitle")

	container.Add(title)

	return container
}
