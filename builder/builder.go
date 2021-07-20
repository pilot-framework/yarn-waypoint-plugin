package builder

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
)

type BuildConfig struct {
	ExecDirectory string `hcl:"directory,optional"`
	OutputDir string `hcl:"output,optional"`
	BaseDir string `hcl:"base,optional"`
}

type Builder struct {
	config BuildConfig
}

// Implement Configurable
func (b *Builder) Config() (interface{}, error) {
	return &b.config, nil
}

// Implement ConfigurableNotify
func (b *Builder) ConfigSet(config interface{}) error {
	c, ok := config.(*BuildConfig)
	if !ok {
		// The Waypoint SDK should ensure this never gets hit
		return fmt.Errorf("Expected *BuildConfig as parameter")
	}

	tmpFiles, err := os.ReadDir("/tmp")
	if err != nil {
		return fmt.Errorf("Error accessing tmp directory")
	}

	tmpDir := ""

	for _, file := range tmpFiles {
		if file.IsDir() && strings.Contains(file.Name(), "waypoint") {
			tmpDir = file.Name()
			break
		}
	}

	if tmpDir == "" {
		return fmt.Errorf("Could not find tmp directory for this project")
	}

	c.BaseDir = path.Join("/tmp", tmpDir)

	return nil
}

// Implement Builder
func (b *Builder) BuildFunc() interface{} {
	// return a function which will be called by Waypoint
	return b.build
}

// A BuildFunc does not have a strict signature, you can define the parameters
// you need based on the Available parameters that the Waypoint SDK provides.
// Waypoint will automatically inject parameters as specified
// in the signature at run time.
//
// Available input parameters:
// - context.Context
// - *component.Source
// - *component.JobInfo
// - *component.DeploymentConfig
// - *datadir.Project
// - *datadir.App
// - *datadir.Component
// - hclog.Logger
// - terminal.UI
// - *component.LabelSet
//
// The output parameters for BuildFunc must be a Struct which can
// be serialzied to Protocol Buffers binary format and an error.
// This Output Value will be made available for other functions
// as an input parameter.
// If an error is returned, Waypoint stops the execution flow and
// returns an error to the user.
func (b *Builder) build(ctx context.Context, ui terminal.UI) (*Binary, error) {
	u := ui.Status()
	defer u.Close()

	if b.config.ExecDirectory == "" {
		b.config.ExecDirectory = b.config.BaseDir
	} else {
		b.config.ExecDirectory = path.Join(b.config.BaseDir, strings.TrimLeft(b.config.ExecDirectory, "./"))
	}

	if b.config.OutputDir == "" {
		b.config.OutputDir = "build"
	}

	u.Update("Installing dependencies required for build process...")

	i := exec.Command(
		"yarn",
		"install",
	)

	i.Dir = b.config.ExecDirectory

	err := i.Run()
	if err != nil {
		u.Step(terminal.StatusError, "Build failed")

		return nil, err
	}

	u.Step("", "Successfully installed dependencies")

	u.Update("Building optimized static files...")

	c := exec.Command(
		"yarn",
		"build",
	)

	c.Dir = b.config.ExecDirectory

	err = c.Run()
	if err != nil {
		u.Step(terminal.StatusError, "Build failed")

		return nil, err
	}

	u.Step(terminal.StatusOK, "Static files built successfully")

	return &Binary{
		Location: path.Join(b.config.ExecDirectory, b.config.OutputDir),
	}, nil
}
