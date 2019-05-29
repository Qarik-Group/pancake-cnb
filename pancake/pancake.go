package pancake

import (
	"os"
	"path/filepath"
	"github.com/cloudfoundry/libcfbuildpack/build"
	"github.com/cloudfoundry/libcfbuildpack/helper"
	"github.com/cloudfoundry/libcfbuildpack/layers"
)

// Dependency is the key used in the build plan by this buildpack
const Dependency = "cf-pancake"

// Contributor is responsibile for deciding what this buildpack will contribute during build
type Contributor struct {
	layer layers.DependencyLayer
}

// NewContributor will create a new Contributor object
func NewContributor(context build.Build) (c Contributor, willContribute bool, err error) {
	plan, wantLayer := context.BuildPlan[Dependency]
	if !wantLayer {
		return Contributor{}, false, nil
	}

	deps, err := context.Buildpack.Dependencies()
	if err != nil {
		return Contributor{}, false, err
	}

	version := plan.Version
	if version == "" {
		if version, err = context.Buildpack.DefaultVersion(Dependency); err != nil {
			return Contributor{}, false, err
		}
	}

	dep, err := deps.Best(Dependency, version, context.Stack)
	if err != nil {
		return Contributor{}, false, err
	}

	contributor := Contributor{
		layer: context.Layers.DependencyLayer(dep),
	}

	return contributor, true, nil
}


// Contribute will install cf-pancake, create profile.d
func (c Contributor) Contribute() error {
	return c.layer.Contribute(func(artifact string, layer layers.DependencyLayer) error {
		layer.Logger.SubsequentLine("Installing to %s", layer.Root)
		if err := helper.ExtractTarXz(artifact, layer.Root, 0); err != nil {
			return err
		}

		pancakeBin, err := filepath.Glob(filepath.Join(layer.Root, "cf-pancake*"))
		if err != nil {
			return err
		}

		err = os.MkdirAll(filepath.Join(layer.Root, "bin"), 0755)
		if err != nil {
			return err
		}

		err = os.Rename(pancakeBin[0], filepath.Join(layer.Root, "bin", "cf-pancake"))
		if err != nil {
			return err
		}

		if err := layer.WriteProfile("0_pancake.sh", runCFPancakeOnStart()); err != nil {
			return err
		}
		return nil
	}, c.flags()...)
}

func (c Contributor) flags() []layers.Flag {
	return []layers.Flag{layers.Cache, layers.Launch}
}

func runCFPancakeOnStart() string {
	return `#!/bin/bash

eval "$(/layers/com.starkandwayne.cf-pancake/cf-pancake/bin/cf-pancake exports)"
`
}