package project

import (
	"fmt"

	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/pauloo27/sonata/common/client"
	"github.com/pauloo27/sonata/common/data"
	"github.com/pauloo27/sonata/gui/utils"
)

func newContentContainer(store *ProjectStore) *gtk.Box {
	container := utils.Must(gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 5))
	container.SetMarginTop(5)
	container.SetMarginBottom(5)
	container.SetMarginStart(5)
	container.SetMarginEnd(5)

	store.ResponseCh = make(chan *client.Response, 2)

	container.Add(newRequestNameContainer(store))
	container.Add(newRequestURLContainer(store))

	subContainer := utils.Must(gtk.PanedNew(gtk.ORIENTATION_VERTICAL))
	subContainer.SetPosition(500)
	subContainer.Add1(newRequestStructureContainer(store))
	subContainer.Add2(newResponseContainer(store))

	container.Add(subContainer)
	container.SetHExpand(true)

	return container
}

func newRequestNameContainer(store *ProjectStore) *gtk.Box {
	nameEntry := utils.Must(gtk.EntryNew())
	nameEntry.SetText(store.DraftRequest.Name)
	nameEntry.Connect("changed", func() {
		store.DraftRequest.Name = utils.Must(nameEntry.GetText())
	})

	saveBtn := utils.Must(gtk.ButtonNewWithLabel("Save"))
	saveBtn.Connect("clicked", func() {
		store.DraftRequest.Headers = store.HeaderStore.ToMap()
		store.DraftRequest.Variables = store.VarStore.ToMap()

		if store.DraftRequest.Name != store.SavedRequest.Name {
			err := store.DraftRequest.Rename(store.DraftRequest.Name)
			if err != nil {
				utils.ShowErrorDialog(store.Window, "Failed to rename request")
				return
			}

			if err := store.Project.ReloadRequests(); err != nil {
				utils.ShowErrorDialog(store.Window, "Failed to reload requests")
				return
			}

			store.ReloadSidebar(store.DraftRequest)
			store.RequestCh <- store.DraftRequest
		} else {
			*store.SavedRequest = *store.DraftRequest
			err := store.SavedRequest.Save()
			if err != nil {
				utils.ShowErrorDialog(store.Window, "Failed to save request")
				return
			}
		}
	})

	container := utils.Must(gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 5))
	container.SetHExpand(true)

	container.PackStart(nameEntry, true, true, 0)
	container.PackEnd(saveBtn, false, false, 0)

	container.SetMarginBottom(5)

	return container
}

func newRequestURLContainer(
	store *ProjectStore,
) *gtk.Box {
	container := utils.Must(gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 5))

	methodsCombo := utils.Must(gtk.ComboBoxTextNew())

	requestMethodIdx := 0

	for i, method := range data.HTTPMethods {
		methodsCombo.AppendText(string(method))
		if method == store.DraftRequest.Method {
			requestMethodIdx = i
		}
	}

	methodsCombo.SetActive(requestMethodIdx)

	methodsCombo.Connect("changed", func() {
		store.DraftRequest.Method = data.HTTPMethod(methodsCombo.GetActiveText())
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
			variables := store.VarStore.ToMap()

			store.DraftRequest.Headers = store.HeaderStore.ToMap()

			fmt.Println("Running request...")
			res, err := client.Run(store.DraftRequest, variables)
			fmt.Println("Request finished", res == nil)
			if err == nil {
				store.ResponseCh <- res
			} else {
				handleRequestError(store, err)
			}

			glib.IdleAdd(func() {
				sendBtn.SetSensitive(true)
			})

		}()
	})

	container.Add(methodsCombo)
	container.Add(entry)
	container.Add(sendBtn)

	container.SetHExpand(true)

	return container
}

func newRequestStructureContainer(store *ProjectStore) *gtk.Notebook {
	container := utils.Must(gtk.NotebookNew())
	headersContainer := newHeadersContainer(store)

	bodyContainer := newBodyContainer(store)

	variablesContainer := newVariablesContainer(store)

	container.AppendPage(variablesContainer, utils.Must(gtk.LabelNew("Variables")))
	container.AppendPage(bodyContainer, utils.Must(gtk.LabelNew("Body")))
	container.AppendPage(headersContainer, utils.Must(gtk.LabelNew("Headers")))
	container.SetVExpand(true)

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

func handleRequestError(store *ProjectStore, err error) {
	utils.ShowErrorDialog(store.Window, err.Error())
}

func newHeadersContainer(store *ProjectStore) *utils.KeyValueEditor {
	kvEditor := utils.NewKeyValueEditor(store.DraftRequest.Headers)
	store.HeaderStore = kvEditor.Store

	return kvEditor
}
