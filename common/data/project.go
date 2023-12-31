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

	RootDir string `json:"-"`

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

	project.RootDir = rootDir
	if err := project.ReloadRequests(); err != nil {
		return nil, err
	}

	return &project, nil
}

func (p *Project) ReloadRequests() error {
	p.requestMap = make(map[string]*Request)
	p.requests = make([]*Request, 0)

	files, err := os.ReadDir(p.RootDir)
	if err != nil {
		return err
	}

	for _, file := range files {
		if !strings.HasSuffix(file.Name(), ".request.json") {
			continue
		}

		data, err := os.ReadFile(path.Join(p.RootDir, file.Name()))
		if err != nil {
			return err
		}

		var request Request
		err = json.Unmarshal(data, &request)
		if err != nil {
			return err
		}

		p.requests = append(p.requests, &request)

		request.path = path.Join(p.RootDir, file.Name())
		request.p = p
		request.Name = strings.TrimSuffix(file.Name(), ".request.json")

		p.requestMap[request.Name] = &request
	}

	return nil
}

func NewProject(rootDir string, name string) (*Project, error) {
	project := Project{
		Name:    name,
		Version: CurrentVersion,
		RootDir: rootDir,
	}

	return &project, nil
}

func (p *Project) Save() error {
	data, err := json.Marshal(p)
	if err != nil {
		return err
	}

	return os.WriteFile(path.Join(p.RootDir, "sonata.json"), data, 420)
}

func (p *Project) GetRequest(name string) (*Request, bool) {
	r, found := p.requestMap[name]
	return r, found
}

func (p *Project) ListRequests() []*Request {
	return p.requests
}

func (p *Project) ListEnvironments() ([]string, error) {
	entries, err := os.ReadDir(p.RootDir)
	if err != nil {
		return nil, err
	}

	environments := make([]string, 0, len(entries))

	for _, entry := range entries {
		if !strings.HasSuffix(entry.Name(), ".env") {
			continue
		}

		environments = append(environments, entry.Name())
	}

	return environments, nil
}
