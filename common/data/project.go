package model

import (
	"encoding/json"
	"os"
	"path"
)

const (
	CurrentVersion = "0.0.1"
)

type Project struct {
	Name    string `json:"name"`
	Version string `json:"version"`

	rootDir string `json:"-"`

	requests map[string]*Request `json:"-"`
}

func LoadProject(rootDir string) (*Project, error) {
	var project Project

	data, err := os.ReadFile(path.Join(rootDir, "sonata.json"))
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &project)
	if err != nil {
		return nil, err
	}

	project.rootDir = rootDir
	return &project, nil
}

func NewProject(rootDir string, name string) (*Project, error) {
	project := Project{
		Name:    name,
		Version: CurrentVersion,
		rootDir: rootDir,
	}

	return &project, nil
}

func (p *Project) Save() error {
	data, err := json.Marshal(p)
	if err != nil {
		return err
	}

	return os.WriteFile(path.Join(p.rootDir, "sonata.json"), data, 420)
}
