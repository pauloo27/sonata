package command

import (
	"os"

	"github.com/pauloo27/sonata/common/data"
	"github.com/spf13/cobra"
)

var NewProject = &cobra.Command{
	Use:   "project",
	Short: "Create a new project",
	Long:  "Create a new project in the current folder",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := args[0]

		dir, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		project, err := data.NewProject(dir, name)
		if err != nil {
			panic(err)
		}

		err = project.Save()
		if err != nil {
			panic(err)
		}
	},
}
