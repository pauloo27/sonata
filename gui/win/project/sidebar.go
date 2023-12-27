package project

import (
	"github.com/gotk3/gotk3/gtk"
	"github.com/pauloo27/sonata/common/data"
	"github.com/pauloo27/sonata/gui/utils"
)

func newSidebar(store *ContentStore) *gtk.Box {
	container := utils.Must(gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0))
	container.SetHExpand(true)

	list := utils.Must(gtk.ListBoxNew())

	for _, request := range store.Project.ListRequests() {
		row := utils.Must(gtk.ListBoxRowNew())
		row.Add(
			utils.Must(gtk.LabelNew(request.Name)),
		)
		list.Add(row)
	}

	list.Connect("row-selected", func() {
		store.RequestCh <- store.Project.
			ListRequests()[list.GetSelectedRow().GetIndex()]
	})

	container.Add(newSidebarHeader(store))
	container.Add(list)

	return container
}

func newSidebarHeader(store *ContentStore) *gtk.HeaderBar {
	container := utils.Must(gtk.HeaderBarNew())
	container.SetTitle("Sonata")

	settingsBtn := utils.Must(
		gtk.MenuButtonNew(),
	)
	settingsBtn.SetImage(
		utils.Must(
			gtk.ImageNewFromIconName("open-menu-symbolic", gtk.ICON_SIZE_BUTTON),
		),
	)

	settingsPopover := utils.Must(gtk.PopoverNew(settingsBtn))

	settingsBtn.SetPopover(settingsPopover)

	settingsPopover.SetPosition(gtk.POS_BOTTOM)
	settingsPopover.SetRelativeTo(settingsBtn)

	// TODO: use MenuModel?
	settingsContainer := utils.Must(gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 5))
	settingsPopover.Add(settingsContainer)

	newRequestBtn := utils.Must(gtk.ButtonNewWithLabel("New Request"))
	aboutBtn := utils.Must(gtk.ButtonNewWithLabel("About"))

	newRequestBtn.Connect("clicked", func() {
		handleNewRequest(store)
	})

	settingsContainer.Add(newRequestBtn)
	settingsContainer.Add(aboutBtn)
	settingsContainer.ShowAll()

	container.PackEnd(settingsBtn)

	return container
}

func handleNewRequest(store *ContentStore) {
	dialog := utils.Must(gtk.DialogNew())
	dialog.SetTitle("New Request")

	dialog.AddButton("Cancel", gtk.RESPONSE_CANCEL)
	dialog.AddButton("Create", gtk.RESPONSE_OK)

	dialog.SetDefaultResponse(gtk.RESPONSE_OK)

	content := utils.Must(gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 5))
	content.SetMarginTop(5)
	content.SetMarginBottom(5)
	content.SetMarginStart(5)
	content.SetMarginEnd(5)

	nameEntry := utils.Must(gtk.EntryNew())
	nameEntry.SetPlaceholderText("Name")
	nameEntry.SetHExpand(true)

	urlEntry := utils.Must(gtk.EntryNew())
	urlEntry.SetPlaceholderText("URL")
	urlEntry.SetHExpand(true)

	methodsCombo := utils.Must(gtk.ComboBoxTextNew())
	for _, method := range data.HTTPMethods {
		methodsCombo.AppendText(string(method))
	}
	methodsCombo.SetActive(0)

	bodyTypeCombo := utils.Must(gtk.ComboBoxTextNew())
	for _, bType := range data.BodyTypes {
		bodyTypeCombo.AppendText(string(bType))
	}
	bodyTypeCombo.SetActive(0)

	content.Add(nameEntry)
	content.Add(urlEntry)
	content.Add(methodsCombo)
	content.Add(bodyTypeCombo)

	utils.Must(dialog.GetContentArea()).Add(content)
	dialog.ShowAll()

	if dialog.Run() == gtk.RESPONSE_OK {
		req := store.Project.NewRequest(
			utils.Must(nameEntry.GetText()),
			"",
			data.HTTPMethod(methodsCombo.GetActiveText()),
			utils.Must(urlEntry.GetText()),
			data.BodyType(bodyTypeCombo.GetActiveText()),
			"",
		)
		err := req.Save()
		if err != nil {
			utils.ShowErrorDialog(nil, "Failed to save request")
			return
		}

		// TODO: reload sidebar instead of begging for a restart
		utils.ShowInfoDialog(nil, "PLEASE, RESTART ME")
	}
}
