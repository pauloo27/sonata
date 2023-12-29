package utils

import (
	"os"

	"github.com/pauloo27/sonata/common/data"
)

func LoadProject() (*data.Project, error) {
	dir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	return data.LoadProject(dir)
}
