package project

import (
	"fmt"
	"os"
	"path"

	"github.com/gotk3/gotk3/gtk"
	"github.com/pauloo27/sonata/common/data"
	"github.com/pauloo27/sonata/gui/utils"
	"github.com/pauloo27/sonata/gui/win"
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

	store.ReloadSidebar = func(selectedRequest *data.Request) {
		utils.ClearChildren(requestsContainer.Container)
		appendRequests(store, requestsContainer, selectedRequest)
		requestsContainer.ShowAll()
	}

	appendRequests(store, requestsContainer, nil)

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
	newEnvBtn := utils.Must(gtk.ButtonNewWithLabel("New Environment"))
	closeProjectBtn := utils.Must(gtk.ButtonNewWithLabel("Close Project"))
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
				nil,
			)
			err := req.Save()
			if err != nil {
				utils.ShowErrorDialog(store.Window, "Failed to create request")
				return
			}

			if err := store.Project.ReloadRequests(); err != nil {
				utils.ShowErrorDialog(store.Window, "Failed to reload requests")
				return
			}

			store.ReloadSidebar(nil)
		}()
	})

	newEnvBtn.Connect("clicked", func() {
		name := utils.ShowEntryDialog(store.Window, "New Environment", "Name:")
		name = name + ".env"

		_, err := os.Create(path.Join(store.Project.RootDir, name))
		if err != nil {
			utils.ShowErrorDialog(store.Window, "Failed to create environment")
		}
	})

	closeProjectBtn.Connect("clicked", func() {
		win.Replace("welcome")
	})

	aboutBtn.Connect("clicked", func() {
		dialog := utils.Must(gtk.AboutDialogNew())
		dialog.SetProgramName("Sonata")
		dialog.SetTitle("About Sonata")
		dialog.SetComments("A stupid file-based REST client made with Go and GTK")
		dialog.SetVersion(data.CurrentVersion)
		dialog.SetLogoIconName("face-yawn-symbolic")
		dialog.SetCopyright("2023 The Sonata Sleepy Dev Team")
		dialog.SetTranslatorCredits("No one")

		_ = dialog.Run()
	})

	settingsContainer.Add(newRequestBtn)
	settingsContainer.Add(newEnvBtn)
	settingsContainer.Add(closeProjectBtn)
	settingsContainer.Add(aboutBtn)
	settingsContainer.ShowAll()

	container.PackEnd(settingsBtn)

	return container
}

var (
	lastSelectedRequest *gtk.Button
)

func appendRequests(
	store *ProjectStore, container *gtk.Box, selectedRequest *data.Request,
) {
	for _, req := range store.Project.ListRequests() {
		requestContainer := utils.Must(gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 0))
		requestContainer.SetHExpand(true)
		requestContainer.SetMarginTop(5)
		requestContainer.SetMarginBottom(5)
		requestContainer.SetMarginStart(5)

		selectRequestBtn := utils.Must(gtk.ButtonNewWithLabel(req.Name))
		if selectedRequest != nil && selectedRequest.Name == req.Name {
			selectRequestBtn.SetSensitive(false)
			lastSelectedRequest = selectRequestBtn
		}

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
				utils.ShowErrorDialog(store.Window, "Failed to delete request")
				return
			}

			if err := store.Project.ReloadRequests(); err != nil {
				utils.ShowErrorDialog(store.Window, "Failed to reload requests")
				return
			}

			store.ReloadSidebar(nil)
		})

		requestContainer.PackEnd(deleteBtn, false, false, 5)

		container.PackStart(requestContainer, false, false, 5)
	}
}
