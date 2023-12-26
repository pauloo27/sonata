package command

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/pauloo27/sonata/common/client"
	"github.com/pauloo27/sonata/common/data"
	"github.com/spf13/cobra"
)

var (
	runKeyValuePairs []string
)

var Run = &cobra.Command{
	Use:   "run",
	Short: "Run a request",
	Long:  "Run a request",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		name := strings.Join(args, " ")
		dir, err := os.Getwd()
		if err != nil {
			panic(err)
		}

		project, err := data.LoadProject(dir)
		if err != nil {
			panic(err)
		}

		name = strings.TrimSuffix(name, ".json")

		request, found := project.GetRequest(name)
		if !found {
			panic("Request not found")
		}

		params := parseRunParams(runKeyValuePairs)

		client := client.NewClient()
		res, err := client.Run(request, params)
		if err != nil {
			panic(err)
		}

		var sb strings.Builder
		title := fmt.Sprintf("# Status code %d\n", res.StatusCode)

		sb.WriteString(title)
		sb.WriteString("```")

		contentType := res.Headers.Get("Content-Type")

		if strings.HasPrefix(contentType, "application/json") {
			sb.WriteString("json")
		}

		sb.WriteString("\n")

		sb.WriteString(res.Body)
		sb.WriteString("\n```")

		in := sb.String()

		out, err := glamour.Render(in, "dark")
		if err != nil {
			panic(err)
		}
		fmt.Print(out)
	},
}

func init() {
	Run.Flags().StringSliceVarP(
		&runKeyValuePairs, "params", "p", []string{}, "-p key=value -p key2=value2",
	)
}

func parseRunParams(keyValuePairs []string) map[string]any {
	params := make(map[string]any)

	for _, pair := range keyValuePairs {
		splitted := strings.SplitN(pair, "=", 2)

		key := splitted[0]
		value := splitted[1]

		params[key] = value
	}

	return params
}
