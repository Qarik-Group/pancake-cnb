package main

import (
	"fmt"
	"os"

	"github.com/buildpack/libbuildpack/buildplan"
	"github.com/cloudfoundry/libcfbuildpack/detect"
	"github.com/starkandwayne/cf-pancake-cnb/pancake"
)

func main() {
	detectionContext, err := detect.DefaultDetect()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to run detection: %s", err)
		os.Exit(101)
	}

	if err := detectionContext.BuildPlan.Init(); err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Failed to initialize Build Plan: %s\n", err)
		os.Exit(101)
	}

	code, err := runDetect(detectionContext)
	if err != nil {
		detectionContext.Logger.Info(err.Error())
	}

	os.Exit(code)
}

func runDetect(context detect.Detect) (int, error) {
	return context.Pass(buildplan.BuildPlan{
		pancake.Dependency: buildplan.Dependency{},
	})
}
