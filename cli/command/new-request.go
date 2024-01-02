package command

import (
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/pauloo27/sonata/cli/editor"
	"github.com/pauloo27/sonata/cli/utils"
	"github.com/pauloo27/sonata/common/data"
	"github.com/spf13/cobra"
)

var NewRequest = &cobra.Command{
	Use:   "request",
	Short: "Create a new request",
	Long:  "Create a new request in interactive mode",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		project, err := utils.LoadProject()
		if err != nil {
			panic(err)
		}

		var name, url, body, rawHeaders string
		var bodyType data.BodyType
		var method data.HTTPMethod

		basicInfoForm := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().Title("Name").Value(&name),
				huh.NewInput().Title("URL").Value(&url),
				huh.NewSelect[data.HTTPMethod]().Title("Method").Value(&method).
					Options(
						huh.NewOption("GET", data.HTTPMethodGet),
						huh.NewOption("POST", data.HTTPMethodPost),
						huh.NewOption("PUT", data.HTTPMethodPut),
						huh.NewOption("DELETE", data.HTTPMethodDelete),
						huh.NewOption("PATCH", data.HTTPMethodPatch),
					),
				huh.NewText().Title("Headers").Placeholder("Header=value").Value(&rawHeaders),
			).Title("Create new sonata HTTP request"),
		)

		if err = basicInfoForm.Run(); err != nil {
			panic(err)
		}

		if method == data.HTTPMethodGet {
			bodyType = data.BodyTypeNone
		} else {
			err = huh.NewSelect[data.BodyType]().Title("Body type").Value(&bodyType).
				Options(
					huh.NewOption("None", data.BodyTypeJSON),
					huh.NewOption("JSON", data.BodyTypeJSON),
					huh.NewOption("Form", data.BodyTypeForm),
					huh.NewOption("Text", data.BodyTypeText),
				).
				Run()
			if err != nil {
				panic(err)
			}
			if bodyType != data.BodyTypeNone {
				body, err = editor.ReadFromEditor(data.BodyTypeExtensions[bodyType])
				if err != nil {
					panic(err)
				}
			}
		}

		headers := parseHeaders(rawHeaders)

		request := project.NewRequest(
			name,
			"", // TODO: add description
			method,
			url,
			bodyType,
			body,
			headers,
		)
		if err := request.Save(); err != nil {
			panic(err)
		}
	},
}

func parseHeaders(raw string) map[string]string {
	headers := make(map[string]string)

	lines := strings.Split(raw, "\n")

	for _, line := range lines {
		parts := strings.SplitN(line, "=", 2)

		name := parts[0]
		value := parts[1]

		headers[name] = value
	}

	return headers
}
