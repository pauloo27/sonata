package project

import (
	"fmt"

	"github.com/gotk3/gotk3/gtk"
	"github.com/pauloo27/sonata/common/data"
	"github.com/pauloo27/sonata/gui/utils"
)

func newSidebar(store *ProjectStore) *gtk.Box {
	container := utils.Must(gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0))
	container.SetHExpand(true)

	requestsContainer := utils.Must(gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0))
	requestsContainer.SetMarginTop(5)
	requestsContainer.SetMarginBottom(5)
	requestsContainer.SetMarginStart(5)
	requestsContainer.SetMarginEnd(5)

	container.Add(newSidebarHeader(store))
	container.Add(utils.Scrollable(requestsContainer))

	store.ReloadSidebar = func() {
		utils.ClearChildren(requestsContainer.Container)
		appendRequests(store, requestsContainer)
		requestsContainer.ShowAll()
	}

	appendRequests(store, requestsContainer)

	return container
}

func newSidebarHeader(store *ProjectStore) *gtk.HeaderBar {
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
	settingsContainer := utils.Must(gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 10))
	settingsContainer.SetMarginTop(5)
	settingsContainer.SetMarginBottom(5)
	settingsContainer.SetMarginStart(5)
	settingsContainer.SetMarginEnd(5)

	settingsPopover.Add(settingsContainer)

	newRequestBtn := utils.Must(gtk.ButtonNewWithLabel("New Request"))
	aboutBtn := utils.Must(gtk.ButtonNewWithLabel("About"))

	newRequestBtn.Connect("clicked", func() {
		go func() {
			var requestName string

			for i := 0; i < 100; i++ {
				nameToCheck := fmt.Sprintf("New Request %d", i)
				_, found := store.Project.GetRequest(nameToCheck)
				if !found {
					requestName = nameToCheck
					break
				}
			}

			req := store.Project.NewRequest(
				requestName,
				"",
				data.HTTPMethodGet,
				"",
				data.BodyTypeJSON,
				"",
			)
			err := req.Save()
			if err != nil {
				utils.ShowErrorDialog(nil, "Failed to create request")
				return
			}

			if err := store.Project.ReloadRequests(); err != nil {
				utils.ShowErrorDialog(nil, "Failed to reload requests")
				return
			}

			store.ReloadSidebar()
		}()
	})

	settingsContainer.Add(newRequestBtn)
	settingsContainer.Add(aboutBtn)
	settingsContainer.ShowAll()

	container.PackEnd(settingsBtn)

	return container
}

var (
	lastSelectedRequest *gtk.Button
)

func appendRequests(store *ProjectStore, container *gtk.Box) {
	for _, req := range store.Project.ListRequests() {
		requestContainer := utils.Must(gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0))
		requestContainer.SetHExpand(true)
		requestContainer.SetMarginTop(5)
		requestContainer.SetMarginBottom(5)
		requestContainer.SetMarginStart(5)

		selectRequestBtn := utils.Must(gtk.ButtonNewWithLabel(req.Name))

		reqCopy := req

		selectRequestBtn.Connect("clicked", func() {
			if lastSelectedRequest != nil {
				lastSelectedRequest.SetSensitive(true)
			}
			selectRequestBtn.SetSensitive(false)
			lastSelectedRequest = selectRequestBtn
			store.RequestCh <- reqCopy
		})

		requestContainer.PackStart(selectRequestBtn, true, true, 5)

		deleteBtn := utils.Must(
			gtk.ButtonNewFromIconName("user-trash-symbolic", gtk.ICON_SIZE_BUTTON),
		)
		deleteBtn.Connect("clicked", func() {
			if lastSelectedRequest == selectRequestBtn {
				store.RequestCh <- nil
				lastSelectedRequest = nil
			}
			if err := reqCopy.Delete(); err != nil {
				utils.ShowErrorDialog(nil, "Failed to delete request")
				return
			}

			if err := store.Project.ReloadRequests(); err != nil {
				utils.ShowErrorDialog(nil, "Failed to reload requests")
				return
			}

			store.ReloadSidebar()
		})

		requestContainer.PackEnd(deleteBtn, false, false, 5)

		container.PackStart(requestContainer, false, false, 5)
	}
}
