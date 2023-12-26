package project

import (
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/pauloo27/sonata/common/client"
	"github.com/pauloo27/sonata/common/data"
	"github.com/pauloo27/sonata/gui/utils"
)

func newContentContainer(request *data.Request) *gtk.Box {
	container, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 5)
	utils.HandleErr(err)

	container.SetMarginTop(5)
	container.SetMarginBottom(5)
	container.SetMarginStart(5)
	container.SetMarginEnd(5)

	paramsStore := newParametersStore()
	responseCh := make(chan *client.Response, 2)

	container.Add(newRequestURLContainer(request, paramsStore, responseCh))

	subContainer, err := gtk.PanedNew(gtk.ORIENTATION_VERTICAL)
	utils.HandleErr(err)

	subContainer.SetPosition(500)
	subContainer.Add1(newRequestStructureContainer(paramsStore))
	subContainer.Add2(newResponseContainer(paramsStore, responseCh))

	container.Add(subContainer)

	container.SetHExpand(true)

	return container
}

func newRequestURLContainer(
	request *data.Request, store *ParameterStore,
	responseCh chan *client.Response,
) *gtk.Box {
	container, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 5)
	utils.HandleErr(err)

	methods, err := gtk.ComboBoxTextNew()
	utils.HandleErr(err)

	requestMethodIdx := 0

	for i, method := range data.HTTPMethods {
		methods.AppendText(string(method))
		if method == request.Method {
			requestMethodIdx = i
		}
	}

	methods.SetActive(requestMethodIdx)

	entry, err := gtk.EntryNew()
	utils.HandleErr(err)

	entry.SetText(request.URL)

	entry.SetHExpand(true)

	sendBtn, err := gtk.ButtonNewWithLabel("Go")
	utils.HandleErr(err)

	sendBtn.Connect("clicked", func() {
		client := client.NewClient()

		params := make(map[string]string)

		for _, param := range store.List() {
			params[param.Key] = param.Value
		}

		res, err := client.Run(request, params)
		utils.HandleErr(err)
		responseCh <- res
	})

	container.Add(methods)
	container.Add(entry)
	container.Add(sendBtn)

	container.SetHExpand(true)

	return container
}

func newRequestStructureContainer(store *ParameterStore) *gtk.Notebook {
	container, err := gtk.NotebookNew()
	utils.HandleErr(err)

	bodyContainer, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	utils.HandleErr(err)

	headersContainer, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	utils.HandleErr(err)

	newLbl := func(label string) *gtk.Label {
		lbl, err := gtk.LabelNew(label)
		utils.HandleErr(err)

		return lbl
	}

	parameterContainer := newParametersContainer(store)

	container.AppendPage(parameterContainer, newLbl("Parameters"))
	container.AppendPage(bodyContainer, newLbl("Body"))
	container.AppendPage(headersContainer, newLbl("Headers"))

	container.SetVExpand(true)

	return container
}

func newResponseContainer(
	store *ParameterStore, responseCh chan *client.Response,
) *gtk.Box {
	container, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	utils.HandleErr(err)

	notebook, err := gtk.NotebookNew()
	utils.HandleErr(err)

	bodyContainer, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	utils.HandleErr(err)

	headersContainer, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	utils.HandleErr(err)

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

	title, err := gtk.LabelNew("Response")
	utils.HandleErr(err)

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

					bodyBuf, err := gtk.TextBufferNew(nil)
					utils.HandleErr(err)

					bodyBuf.SetText(response.Body)

					bodyView, err := gtk.TextViewNewWithBuffer(bodyBuf)
					utils.HandleErr(err)

					bodyView.SetEditable(false)
					bodyView.SetHExpand(true)

					bodyScroll, err := gtk.ScrolledWindowNew(nil, nil)
					utils.HandleErr(err)

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
