package project

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/pauloo27/sonata/gui/utils"
)

func newContentContainer() *gtk.Box {
	container, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 5)
	utils.HandleErr(err)

	container.SetMarginTop(5)
	container.SetMarginBottom(5)
	container.SetMarginStart(5)
	container.SetMarginEnd(5)

	container.Add(newRequestURLContainer())

	subContainer, err := gtk.PanedNew(gtk.ORIENTATION_VERTICAL)
	utils.HandleErr(err)

	subContainer.SetPosition(500)
	subContainer.Add1(newRequestStructureContainer())
	subContainer.Add2(newResponseContainer())

	container.Add(subContainer)

	container.SetHExpand(true)

	return container
}

func newRequestURLContainer() *gtk.Box {
	container, err := gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 5)
	utils.HandleErr(err)

	methods, err := gtk.ComboBoxTextNew()
	utils.HandleErr(err)

	methods.AppendText("GET")
	methods.AppendText("POST")
	methods.AppendText("PUT")
	methods.AppendText("DELETE")
	methods.AppendText("PATCH")
	methods.AppendText("HEAD")
	methods.AppendText("OPTIONS")

	entry, err := gtk.EntryNew()
	utils.HandleErr(err)

	entry.SetHExpand(true)

	sendBtn, err := gtk.ButtonNewWithLabel("Go")
	utils.HandleErr(err)

	container.Add(methods)
	container.Add(entry)
	container.Add(sendBtn)

	container.SetHExpand(true)

	return container
}

func newRequestStructureContainer() *gtk.Notebook {
	container, err := gtk.NotebookNew()
	utils.HandleErr(err)

	overviewContainer, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
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

	overviewContainer.Add(newLbl("Overview"))
	bodyContainer.Add(newLbl("Body"))
	headersContainer.Add(newLbl("Headers"))

	container.AppendPage(overviewContainer, newLbl("Overview"))
	container.AppendPage(bodyContainer, newLbl("Body"))
	container.AppendPage(headersContainer, newLbl("Headers"))

	container.SetVExpand(true)

	return container
}

func newResponseContainer() *gtk.Box {
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

	return container
}
