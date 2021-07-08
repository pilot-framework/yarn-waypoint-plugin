package builder

import (
	"context"
	"exec"
	"fmt"

	"github.com/hashicorp/waypoint-plugin-sdk/terminal"
)

type BuildConfig struct {
	ExecDirectory string `hcl:"directory,optional"`
	OutputDir string `hcl:"output,optional"`
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

	_, err := os.Stat(c.ExecDirectory)

	// validate the config
	if err != nil {
		return fmt.Errorf("Directory you specified Yarn to be executed in does not exist")
	}

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
	u.Update("Building application")

	if b.config.ExecDirectory == "" {
		b.config.ExecDirectory = "./"
	}

	if b.config.OutputDir == "" {
		b.config.OutputDir = "build"
	}

	c := exec.Command(
		"yarn"
		"build"
	)

	c.Path = b.config.ExecDirectory

	err := c.Run()
	if err != nil {
		u.Step(terminal.StatusError, "Build failed")

		return nil, err
	}

	u.Step(terminal.StatusOK, "Static files build successfully")

	return &Binary{
		Location: path.Join(b.config.Source, b.config.OutputDir),
	}, nil
}
