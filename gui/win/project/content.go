package project

import (
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/pauloo27/sonata/common/client"
	"github.com/pauloo27/sonata/common/data"
	"github.com/pauloo27/sonata/gui/utils"
)

func newContentContainer(request *data.Request) *gtk.Box {
	container := utils.Must(gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 5))
	container.SetMarginTop(5)
	container.SetMarginBottom(5)
	container.SetMarginStart(5)
	container.SetMarginEnd(5)

	varStore := newVariablesStore()
	responseCh := make(chan *client.Response, 2)

	container.Add(newRequestURLContainer(request, varStore, responseCh))

	subContainer := utils.Must(gtk.PanedNew(gtk.ORIENTATION_VERTICAL))
	subContainer.SetPosition(500)
	subContainer.Add1(newRequestStructureContainer(varStore))
	subContainer.Add2(newResponseContainer(varStore, responseCh))

	container.Add(subContainer)
	container.SetHExpand(true)

	return container
}

func newRequestURLContainer(
	request *data.Request, store *VariablesStore,
	responseCh chan *client.Response,
) *gtk.Box {
	container := utils.Must(gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 5))

	methods := utils.Must(gtk.ComboBoxTextNew())

	requestMethodIdx := 0

	for i, method := range data.HTTPMethods {
		methods.AppendText(string(method))
		if method == request.Method {
			requestMethodIdx = i
		}
	}

	methods.SetActive(requestMethodIdx)

	entry := utils.Must(gtk.EntryNew())
	entry.SetText(request.URL)
	entry.SetHExpand(true)

	sendBtn := utils.Must(gtk.ButtonNewWithLabel("Go"))
	sendBtn.Connect("clicked", func() {
		client := client.NewClient()

		variables := make(map[string]string)

		for _, variable := range store.List() {
			variables[variable.Key] = variable.Value
		}

		res, err := client.Run(request, variables)
		// FIXME: proper error handling
		utils.HandleErr(err)
		responseCh <- res
	})

	container.Add(methods)
	container.Add(entry)
	container.Add(sendBtn)

	container.SetHExpand(true)

	return container
}

func newRequestStructureContainer(store *VariablesStore) *gtk.Notebook {
	container := utils.Must(gtk.NotebookNew())
	bodyContainer := utils.Must(gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0))
	headersContainer := utils.Must(gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0))

	parameterContainer := newParametersContainer(store)

	container.AppendPage(parameterContainer, utils.Must(gtk.LabelNew("Parameters")))
	container.AppendPage(bodyContainer, utils.Must(gtk.LabelNew("Body")))
	container.AppendPage(headersContainer, utils.Must(gtk.LabelNew("Headers")))
	container.SetVExpand(true)

	return container
}

func newResponseContainer(
	store *VariablesStore, responseCh chan *client.Response,
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
			response := <-responseCh
			if response != nil {
				glib.IdleAdd(func() {
					notebook.SetCurrentPage(0)

					bodyContainer.GetChildren().Foreach(func(item interface{}) {
						item.(*gtk.Widget).Destroy()
					})
					headersContainer.GetChildren().Foreach(func(item interface{}) {
						item.(*gtk.Widget).Destroy()
					})

					bodyBuf := utils.Must(gtk.TextBufferNew(nil))
					bodyBuf.SetText(response.Body)

					bodyView := utils.Must(gtk.TextViewNewWithBuffer(bodyBuf))
					bodyView.SetEditable(false)
					bodyView.SetHExpand(true)

					bodyScroll := utils.Must(gtk.ScrolledWindowNew(nil, nil))

					bodyScroll.SetVExpand(true)
					bodyScroll.Add(bodyView)

					bodyContainer.Add(bodyScroll)

					bodyContainer.ShowAll()
					headersContainer.ShowAll()
				})
			}
		}
	}()

	return container
}
