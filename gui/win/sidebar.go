package win

import "github.com/gotk3/gotk3/gtk"

var (
	sampleList = []string{
		"Get user",
		"Get users",
		"Create user",
		"Update user",
		"Delete user",
	}
)

func newSidebar() *gtk.Box {
	container, err := gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0)
	handleErr(err)

	container.SetHExpand(true)

	list, err := gtk.ListBoxNew()
	handleErr(err)

	for _, item := range sampleList {
		row, err := gtk.ListBoxRowNew()
		handleErr(err)

		label, err := gtk.LabelNew(item)
		handleErr(err)

		row.Add(label)
		list.Add(row)
	}

	container.Add(newSidebarHeader())
	container.Add(list)

	return container
}

func newSidebarHeader() *gtk.HeaderBar {
	container, err := gtk.HeaderBarNew()
	handleErr(err)

	container.SetTitle("Sonata")

	settingsBtn, err := gtk.ButtonNewFromIconName("open-menu-symbolic", gtk.ICON_SIZE_BUTTON)
	handleErr(err)

	container.PackStart(settingsBtn)

	return container
}
