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

					grid := utils.Must(gtk.GridNew())
					grid.SetColumnHomogeneous(true)

					counter := 0

					for key, values := range headers {
						keyLbl := utils.Must(gtk.LabelNew(key))
						valuesLbl := utils.Must(gtk.LabelNew(values[0]))

						keyLbl.SetSelectable(true)
						valuesLbl.SetSelectable(true)

						grid.Attach(utils.Scrollable(keyLbl), 0, counter, 1, 1)
						grid.Attach(utils.Scrollable(valuesLbl), 1, counter, 1, 1)

						counter++
					}

					headersContainer.Add(utils.Scrollable(grid))
					headersContainer.SetVExpand(true)
					headersContainer.ShowAll()
				})
			}
		}
	}()

	return container
}
