package env

import (
	"fmt"
	"os"

	"github.com/gotk3/gotk3/gtk"
	"github.com/joho/godotenv"
	"github.com/pauloo27/sonata/common/client"
	"github.com/pauloo27/sonata/gui/utils"
)

func newContentContainer(store *EnvStore) *gtk.Box {
	container := utils.Must(gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 5))

	container.SetMarginTop(5)
	container.SetMarginBottom(5)
	container.SetMarginStart(5)
	container.SetMarginEnd(5)

	label := utils.Must(gtk.LabelNew(fmt.Sprintf("Editing env %s", store.EnvName)))
	utils.AddCSSClass(label.Widget, "welcome-subtitle")

	container.PackStart(label, false, false, 0)

	editor := utils.NewKeyValueEditor(store.Envs)

	container.PackStart(editor, true, true, 0)

	saveBtn := utils.Must(gtk.ButtonNewWithLabel("Save"))
	closeBtn := utils.Must(gtk.ButtonNewWithLabel("Close"))

	saveBtn.Connect("clicked", func() {
		envs := editor.Store.ToMap()
		raw, err := godotenv.Marshal(envs)
		if err != nil {
			utils.ShowErrorDialog(store.Window, "Failed to generate env")
			return
		}

		envFilePath := fmt.Sprintf("%s/%s", store.Project.RootDir, store.EnvName)

		err = os.WriteFile(envFilePath, []byte(raw), 0644)
		if err != nil {
			utils.ShowErrorDialog(store.Window, "Failed to save env")
			return
		}

		client.UseMap(envs)
	})

	closeBtn.Connect("clicked", func() {
		store.Window.Destroy()
	})

	btnContainer := utils.Must(gtk.BoxNew(gtk.ORIENTATION_HORIZONTAL, 5))

	btnContainer.PackStart(saveBtn, false, false, 0)
	btnContainer.PackEnd(closeBtn, false, false, 0)

	container.PackEnd(btnContainer, false, false, 0)

	return container
}
