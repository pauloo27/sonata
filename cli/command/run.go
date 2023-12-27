package command

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/glamour"
	gloss "github.com/charmbracelet/lipgloss"
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

		name = strings.TrimSuffix(name, ".request.json")

		request, found := project.GetRequest(name)
		if !found {
			panic("Request not found")
		}

		variables := parseRunVariables(runKeyValuePairs)

		client := client.NewClient()
		res, err := client.Run(request, variables)
		if err != nil {
			panic(err)
		}

		showRes(request, res)
	},
}

func init() {
	Run.Flags().StringSliceVarP(
		&runKeyValuePairs, "variables", "v", []string{}, "-v key=value -v key2=value2",
	)
}

func parseRunVariables(keyValuePairs []string) map[string]string {
	variables := make(map[string]string)

	for _, pair := range keyValuePairs {
		splitted := strings.SplitN(pair, "=", 2)

		key := splitted[0]
		value := splitted[1]

		variables[key] = value
	}

	return variables
}

var (
	keyStyle   = gloss.NewStyle().Foreground(gloss.Color("4"))
	valueStyle = gloss.NewStyle().Foreground(gloss.Color("7"))
	Success    = gloss.NewStyle().Foreground(gloss.Color("10")).Background(gloss.Color("0"))
	Failure    = gloss.NewStyle().Foreground(gloss.Color("9")).Background(gloss.Color("0"))
)

func showRes(req *data.Request, res *client.Response) {
	printField("Name", req.Name)
	printField("URL", res.CalledURL)
	printField("Method", string(req.Method))

	statusCodeStyle := Failure
	if res.StatusCode >= 100 && res.StatusCode < 300 {
		statusCodeStyle = Success
	}

	printStyledField(
		keyStyle.Render("Status Code"),
		statusCodeStyle.Render(
			fmt.Sprintf("%d %s", res.StatusCode, http.StatusText(res.StatusCode)),
		),
	)

	printField("Took", res.Time.Truncate(time.Millisecond).String())

	var sb strings.Builder
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
}

func printField(key string, value string) {
	fmt.Printf("%s: %s\n", keyStyle.Render(key), valueStyle.Render(value))
}

func printStyledField(key string, value string) {
	fmt.Printf("%s: %s\n", key, value)
}
