package pancake

import (
	// "fmt"
	// "path/filepath"

	"github.com/buildpack/libbuildpack/application"
	"github.com/cloudfoundry/libcfbuildpack/build"
	"github.com/cloudfoundry/libcfbuildpack/helper"
	"github.com/cloudfoundry/libcfbuildpack/layers"
)

// Contributor is responsibile for deciding what this buildpack will contribute during build
type Contributor struct {
	app                application.Application
	launchContribution bool
	launchLayer        layers.Layers
	httpdLayer         layers.DependencyLayer
}

// NewContributor will create a new Contributor object
func NewContributor(context build.Build) (c Contributor, willContribute bool, err error) {
	plan, wantDependency := context.BuildPlan[Dependency]
	if !wantDependency {
		return Contributor{}, false, nil
	}

	deps, err := context.Buildpack.Dependencies()
	if err != nil {
		return Contributor{}, false, err
	}

	dep, err := deps.Best(Dependency, plan.Version, context.Stack)
	if err != nil {
		return Contributor{}, false, err
	}

	contributor := Contributor{
		app:         context.Application,
		launchLayer: context.Layers,
		httpdLayer:  context.Layers.DependencyLayer(dep),
	}

	return contributor, true, nil
}

// Contribute will install cf-pancake, create profile.d
func (c Contributor) Contribute() error {
	return c.httpdLayer.Contribute(func(artifact string, layer layers.DependencyLayer) error {
		layer.Logger.SubsequentLine("Installing to %s", layer.Root)
		if err := helper.CopyFile(artifact, layer.Root); err != nil {
			return err
		}
		// if err := helper.ExtractTarXz(artifact, layer.Root, 1); err != nil {
		// 	return err
		// }

		return c.launchLayer.WriteApplicationMetadata(layers.Metadata{})
	}, c.flags()...)
}

func (c Contributor) flags() []layers.Flag {
	var flags []layers.Flag

	if c.launchContribution {
		flags = append(flags, layers.Launch)
	}

	return flags
}
