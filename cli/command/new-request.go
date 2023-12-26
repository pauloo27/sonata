package command

import (
	"os"

	"github.com/charmbracelet/huh"
	"github.com/pauloo27/sonata/common/data"
	"github.com/spf13/cobra"
)

var NewRequest = &cobra.Command{
	Use:   "request",
	Short: "Create a new request",
	Long:  "Create a new request in interactive mode",
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		dir, err := os.Getwd()
		if err != nil {
			panic(err)
		}

		project, err := data.LoadProject(dir)
		if err != nil {
			panic(err)
		}

		var name, url, body string
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
				huh.NewText().Title("Body").Value(&body).Run()
			}
		}

		request := project.NewRequest(
			name,
			"", // TODO: add description
			method,
			url,
			bodyType,
			body,
		)
		if err := request.Save(); err != nil {
			panic(err)
		}
	},
}
