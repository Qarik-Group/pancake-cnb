package pancake

import (
	"github.com/cloudfoundry/libcfbuildpack/build"
	"github.com/cloudfoundry/libcfbuildpack/helper"
	"github.com/cloudfoundry/libcfbuildpack/layers"
)

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
		return helper.CopyFile(artifact, layer.Root)
	}, c.flags()...)
}

func (c Contributor) flags() []layers.Flag {
	return []layers.Flag{layers.Cache}
}
