package data

import (
	"encoding/json"
	"os"
	"path"
	"strings"
)

const (
	CurrentVersion = "0.0.1"
)

type Project struct {
	Name    string `json:"name"`
	Version string `json:"version"`

	rootDir string `json:"-"`

	requests   []*Request          `json:"-"`
	requestMap map[string]*Request `json:"-"`
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
	if err := loadProjectRequests(&project); err != nil {
		return nil, err
	}

	return &project, nil
}

func loadProjectRequests(project *Project) error {
	project.requestMap = make(map[string]*Request)

	files, err := os.ReadDir(project.rootDir)
	if err != nil {
		return err
	}

	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".json") {
			continue
		}

		data, err := os.ReadFile(path.Join(project.rootDir, file.Name()))
		if err != nil {
			return err
		}

		var request Request
		err = json.Unmarshal(data, &request)
		if err != nil {
			return err
		}

		project.requests = append(project.requests, &request)

		request.path = path.Join(project.rootDir, file.Name())
		request.p = project
		request.Name = strings.TrimSuffix(file.Name(), ".json")

		project.requestMap[request.Name] = &request
	}

	return nil
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

func (p *Project) GetRequest(name string) (*Request, bool) {
	r, found := p.requestMap[name]
	return r, found
}

func (p *Project) ListRequests() []*Request {
	return p.requests
}
