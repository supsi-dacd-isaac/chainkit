package project

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

// ChainkitManifest defines the name of the manifest file
const ChainkitManifest = "chainkit.yml"

type projectBinaries struct {
	CLI    string
	Daemon string
}

// Project represents a project
type Project struct {
	Name     string
	RootDir  string `yaml:"-"`
	Binaries *projectBinaries
	Image    string
}

// New will create a new project in the given directory.
func New(dir, name string) *Project {
	p := &Project{
		Name: name,
		Binaries: &projectBinaries{
			CLI:    name + "cli",
			Daemon: name + "d",
		},
		RootDir: path.Join(dir, name),
	}
	p.SetDefaults()
	return p
}

// Save serializes the project data on disk
func (p *Project) Save() error {
	ybuf, err := yaml.Marshal(p)
	if err != nil {
		return err
	}
	fp, err := os.Create(path.Join(p.RootDir, ChainkitManifest))
	if err != nil {
		return err
	}
	if _, err = fp.Write(ybuf); err != nil {
		return err
	}
	return nil
}

// Validate runs sanity checks against the project
func (p *Project) Validate() error {
	errorOut := func(field string) error {
		return fmt.Errorf("missing required field %q", field)
	}

	switch {
	case p.Name == "":
		return errorOut("name")
	case p.Binaries == nil:
		return errorOut("binaries")
	case p.Binaries.CLI == "":
		return errorOut("binaries.cli")
	case p.Binaries.Daemon == "":
		return errorOut("binaries.daemon")
	}

	return nil
}

// SetDefaults sets the project default values.
func (p *Project) SetDefaults() {
	switch {
	case p.Image == "":
		p.Image = fmt.Sprintf("chainkit-%s", p.Name)
	}
}

// Load will load a project from a given directory
func Load(dir string) (*Project, error) {
	errMsg := fmt.Sprintf("Cannot read manifest %q", ChainkitManifest)
	data, err := ioutil.ReadFile(path.Join(dir, ChainkitManifest))
	if err != nil {
		return nil, errors.Wrap(err, errMsg)
	}
	p := &Project{}
	if err = yaml.Unmarshal(data, p); err != nil {
		return nil, errors.Wrap(err, errMsg)
	}
	p.RootDir = dir

	if err := p.Validate(); err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("%s validation", ChainkitManifest))
	}

	p.SetDefaults()

	return p, nil
}

// StateDir returns the state directory within the project.
func (p *Project) StateDir() string {
	return path.Join(p.RootDir, "state")
}

// DataDir returns the data directory within the project state.
func (p *Project) DataDir() string {
	return path.Join(p.StateDir(), "data")
}

// ConfigDir returns the config directory within the project state.
func (p *Project) ConfigDir() string {
	return path.Join(p.StateDir(), "config")
}

// GenesisPath returns the genesis path for the project.
func (p *Project) GenesisPath() string {
	return path.Join(p.ConfigDir(), "genesis.json")
}

// CLIDir returns the CLI directory within the project state.
func (p *Project) CLIDir() string {
	return path.Join(p.StateDir(), "cli")
}

// IPFSDir returns the IPFS data directory within the project state.
func (p *Project) IPFSDir() string {
	return path.Join(p.StateDir(), "ipfs")
}