package command

import (
	"fmt"
	"runtime"

	"github.com/pauloo27/sonata/common/data"
	"github.com/spf13/cobra"
)

var (
	Version = &cobra.Command{
		Use:   "version",
		Short: "Print the version",
		Long:  "Print the current version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf(
				"Sonata v%s (%s, compiled with %s for %s/%s)\n",
				data.CurrentVersion,
				runtime.Version(), runtime.Compiler, runtime.GOOS, runtime.GOARCH,
			)
		},
	}
)
