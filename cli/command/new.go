package command

import "github.com/spf13/cobra"

var New = &cobra.Command{
	Use:   "new",
	Short: "Create a new thing",
	Long:  "Create a new project or request",
	Run: func(cmd *cobra.Command, args []string) {
		_ = cmd.Help()
	},
}

func init() {
	New.AddCommand(NewProject)
	New.AddCommand(NewRequest)
}
