package project

import (
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
	"github.com/pauloo27/sonata/gui/utils"
)

func newResponseContainer(
	store *ProjectStore,
) *gtk.Box {
	container := utils.Must(gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0))

	notebook := utils.Must(gtk.NotebookNew())

	bodyContainer := utils.Must(gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0))

	headersContainer := utils.Must(gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0))

	notebook.AppendPage(bodyContainer, utils.Must(gtk.LabelNew(("Body"))))
	notebook.AppendPage(headersContainer, utils.Must(gtk.LabelNew("Headers")))
	notebook.SetVExpand(true)

	title := utils.Must(gtk.LabelNew("Response"))
	title.SetHAlign(gtk.ALIGN_CENTER)

	container.Add(title)
	container.Add(notebook)

	go func() {
		for {
			response := <-store.ResponseCh
			if response != nil {
				body := response.Body
				headers := response.Headers

				glib.IdleAdd(func() {
					notebook.SetCurrentPage(0)

					utils.ClearChildren(bodyContainer.Container)
					utils.ClearChildren(headersContainer.Container)

					bodyContainer.Add(utils.NewEditor(body, false, "json"))
					bodyContainer.ShowAll()

					values := make(map[string]string)
					for key, value := range headers {
						values[key] = value[0]
					}

					headerEditor := utils.NewKeyValueEditor(values)

					headerEditor.SetVExpand(true)
					headersContainer.SetVExpand(true)
					headersContainer.Add(headerEditor)
					headersContainer.ShowAll()
				})
			}
		}
	}()

	return container
}
