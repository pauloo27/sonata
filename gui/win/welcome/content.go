package welcome

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/gotk3/gotk3/gtk"
	"github.com/pauloo27/sonata/common/data"
	"github.com/pauloo27/sonata/gui/utils"
	"github.com/pauloo27/sonata/gui/win"
)

const (
	cacheDirPath       = "$HOME/.cache/sonata"
	recentProjectsPath = cacheDirPath + "/recent_projects"
	maxRecentProjects  = 5
)

func newContentContainer(gtkWin *gtk.Window) *gtk.Box {
	contentContainer := utils.Must(gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 0))
	contentContainer.SetVExpand(true)
	contentContainer.SetHAlign(gtk.ALIGN_CENTER)
	contentContainer.SetVAlign(gtk.ALIGN_CENTER)

	welcomeLbl := utils.Must(gtk.LabelNew("Welcome to"))
	utils.AddCSSClass(welcomeLbl.Widget, "welcome-subtitle")

	sonataLbl := utils.Must(gtk.LabelNew("💤 SONATA 💤"))
	utils.AddCSSClass(sonataLbl.Widget, "welcome-title")

	createNewBtn := utils.Must(gtk.ButtonNewWithLabel("Create new project"))
	createNewBtn.Connect("clicked", func() {
		selectedPath := utils.ChooseFolder(gtkWin, "Select project folder")
		if selectedPath == "" {
			return
		}

		name := path.Base(selectedPath)

		project, err := data.NewProject(selectedPath, name)
		if err != nil {
			utils.ShowErrorDialog(gtkWin, "Failed to create project")
			return
		}

		if err = project.Save(); err != nil {
			utils.ShowErrorDialog(gtkWin, "Failed to save project")
			return
		}

		err = addProjectToRecent(selectedPath)
		if err != nil {
			utils.ShowErrorDialog(gtkWin, "Failed to add project to recents")
		}

		win.Replace("project", project)
	})

	openBtn := utils.Must(gtk.ButtonNewWithLabel("Open project"))
	openBtn.Connect("clicked", func() {
		selectedPath := utils.ChooseFolder(gtkWin, "Open project")
		if selectedPath == "" {
			return
		}

		project, err := data.LoadProject(selectedPath)
		if err != nil {
			utils.ShowErrorDialog(
				gtkWin, fmt.Sprintf(
					"Failed to load project: %s",
					selectedPath,
				),
			)
			return
		}

		err = addProjectToRecent(selectedPath)
		if err != nil {
			println(err.Error())
			utils.ShowErrorDialog(gtkWin, "Failed to add project to recents")
		}
		win.Replace("project", project)
	})

	contentContainer.PackStart(welcomeLbl, false, false, 5)
	contentContainer.PackStart(sonataLbl, false, false, 5)
	contentContainer.PackStart(createNewBtn, false, false, 5)
	contentContainer.PackStart(openBtn, false, false, 5)
	contentContainer.PackStart(newRecentProjectsContainer(gtkWin), false, false, 0)

	return contentContainer
}

func newRecentProjectsContainer(parent *gtk.Window) *gtk.Box {
	container := utils.Must(gtk.BoxNew(gtk.ORIENTATION_VERTICAL, 5))

	container.PackStart(
		utils.Must(gtk.LabelNew("Recent projects:")),
		false, false, 5,
	)

	recentProjectsPath := listRecentProjects()

	for _, projectPath := range recentProjectsPath {
		projectBtn := utils.Must(gtk.ButtonNewWithLabel(projectPath))

		projectPathCopy := projectPath

		projectBtn.Connect("clicked", func() {
			project, err := data.LoadProject(projectPathCopy)
			if err != nil {
				utils.ShowErrorDialog(
					parent, fmt.Sprintf(
						"Failed to load project: %s",
						projectPathCopy,
					),
				)
				return
			}

			err = addProjectToRecent(projectPathCopy)
			if err != nil {
				utils.ShowErrorDialog(parent, "Failed to update project to recents")
			}
			win.Replace("project", project)
		})

		container.PackStart(projectBtn, false, false, 0)
	}

	return container
}

func listRecentProjects() []string {
	finalPath := os.ExpandEnv(recentProjectsPath)

	/* #nosec G304 */
	/* #nosec G302 */
	file, err := os.OpenFile(finalPath, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		return nil
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	var projects []string

	for i := 0; i < maxRecentProjects; i++ {
		data, _, err := reader.ReadLine()
		if err != nil {
			break
		}
		line := string(data)
		if line == "" {
			continue
		}

		projects = append(projects, line)
	}

	return projects
}

func addProjectToRecent(projectPath string) error {
	/* #nosec G304 */
	/* #nosec G302 */

	err := os.MkdirAll(os.ExpandEnv(cacheDirPath), 0755)
	if err != nil {
		return err
	}

	file, err := os.OpenFile(
		os.ExpandEnv(recentProjectsPath),
		os.O_RDWR|os.O_CREATE, 0644,
	)
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		return err
	}
	lines := strings.Split(string(data), "\n")

	newRecentProjectFile := strings.Builder{}

	counter := 1
	newRecentProjectFile.WriteString(projectPath + "\n")

	for _, line := range lines {
		if line == projectPath || line == "" {
			continue
		}
		newRecentProjectFile.WriteString(line + "\n")
		counter++

		if counter >= maxRecentProjects {
			break
		}
	}

	if err = file.Truncate(0); err != nil {
		return err
	}

	if _, err = file.Seek(0, 0); err != nil {
		return err
	}
	_, err = file.WriteString(newRecentProjectFile.String())
	return err
}
