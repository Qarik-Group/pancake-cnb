module github.com/starkandwayne/cf-pancake-cnb

go 1.12

require (
	github.com/buildpack/libbuildpack v1.16.0
	github.com/cloudfoundry/dagger v0.0.0-20190521201554-93417312948c
	github.com/cloudfoundry/libcfbuildpack v1.44.0
	github.com/onsi/gomega v1.5.0
	github.com/sclevine/spec v1.2.0
	google.golang.org/appengine v1.5.0 // indirect
)

replace github.com/cloudfoundry/libcfbuildpack => github.com/drnic/libcfbuildpack v1.55.1-0.20190529014517-75dbd8e77483
