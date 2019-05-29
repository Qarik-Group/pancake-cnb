package integration

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/sclevine/spec"
	"github.com/sclevine/spec/report"

	"github.com/cloudfoundry/dagger"
	. "github.com/onsi/gomega"
)

var (
	bpDir, phpBP, phpWebBP, httpdBP, pancakeBP string
)

func TestIntegration(t *testing.T) {
	var err error
	Expect := NewWithT(t).Expect
	bpDir, err = dagger.FindBPRoot()
	Expect(err).NotTo(HaveOccurred())
	pancakeBP, err = dagger.PackageBuildpack(bpDir)
	Expect(err).ToNot(HaveOccurred())
	defer os.RemoveAll(pancakeBP)

	phpBP, err = dagger.GetLatestBuildpack("php-cnb")
	phpWebBP, err = dagger.GetLatestBuildpack("php-web-cnb")
	httpdBP, err = dagger.GetLatestBuildpack("httpd-cnb")
	Expect(err).ToNot(HaveOccurred())
	defer os.RemoveAll(phpBP)
	defer os.RemoveAll(phpWebBP)
	defer os.RemoveAll(httpdBP)

	spec.Run(t, "Integration", testIntegration, spec.Report(report.Terminal{}))
}

func testIntegration(t *testing.T, when spec.G, it spec.S) {
	var Expect func(interface{}, ...interface{}) GomegaAssertion
	it.Before(func() {
		Expect = NewWithT(t).Expect
	})

	it("should build a working OCI image for a simple app", func() {
		app, err := dagger.PackBuild(filepath.Join("fixtures", "phpapp"), pancakeBP, phpBP, httpdBP, phpWebBP)
		vcapServices, err := ioutil.ReadFile(filepath.Join("fixtures", "vcap_services", "p-mysql.json"))
		Expect(err).ToNot(HaveOccurred())
		app.Env["VCAP_APPLICATION"] = "{}"
		app.Env["VCAP_SERVICES"] = string(vcapServices)
		Expect(err).ToNot(HaveOccurred())
		// TODO: restore app.Destroy when no longer debugging
		// defer app.Destroy()

		Expect(app.Start()).To(Succeed())
		Expect(app.HTTPGetBody("/")).To(ContainSubstring("PHP Version"))
		Expect(app.HTTPGetBody("/")).To(ContainSubstring("MYSQL_HOSTNAME"))
		Expect(app.HTTPGetBody("/")).To(ContainSubstring("P_MYSQL_PASSWORD"))
	})

	// when("the app is pushed twice", func() {
	// 	it("does not reinstall go modules", func() {
	// 		app, err := dagger.PackBuild(filepath.Join("fixtures", "phpapp"), phpURI, pancakeURI)
	// 		Expect(err).ToNot(HaveOccurred())
	// 		defer app.Destroy()

	// 		Expect(app.Start()).To(Succeed())
	// 		Expect(app.HTTPGetBody("/")).To(ContainSubstring("PHP Version"))

	// 		_, imageID, _, err := app.Info()
	// 		Expect(err).NotTo(HaveOccurred())

	// 		app, err = dagger.PackBuildNamedImage(imageID, appDir, phpURI, pancakeURI)
	// 		Expect(err).ToNot(HaveOccurred())

	// 		Expect(app.Start()).To(Succeed())
	// 		Expect(app.HTTPGetBody("/")).To(ContainSubstring("PHP Version"))
	// 	})
	// })

	// when("the app is vendored", func() {
	// 	it("builds an OCI image without downloading any extra packages", func() {
	// 		app, err := dagger.PackBuild(filepath.Join("fixtures", "phpapp"), phpURI, pancakeURI)
	// 		Expect(err).ToNot(HaveOccurred())

	// 		// Expect(app.BuildLogs()).NotTo(MatchRegexp(goDownloading))

	// 		Expect(app.Start()).To(Succeed())
	// 		Expect(app.HTTPGetBody("/")).To(ContainSubstring("PHP Version"))
	// 	})
	// })
}
