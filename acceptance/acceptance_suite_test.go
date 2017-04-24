package acceptance_test

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"

	"testing"
)

const name = "marketplace-v2"
const timeout = time.Second * 60

var (
	service           string
	planWithSchema    string
	planWithoutSchema string
)

var _ = BeforeSuite(func() {
	api := os.Getenv("CF_API")
	username := os.Getenv("CF_USERNAME")
	password := os.Getenv("CF_PASSWORD")
	org := os.Getenv("CF_ORG")
	space := os.Getenv("CF_SPACE")
	service = os.Getenv("CF_SERVICE")
	planWithSchema = os.Getenv("CF_PLAN_WITH_SCHEMA")
	planWithoutSchema = os.Getenv("CF_PLAN_WITHOUT_SCHEMA")

	mustHaveAll(api, username, password, org, space, service, planWithSchema, planWithoutSchema)

	var err error
	path, err := gexec.Build("github.com/pivotal-cf/marketplacev2")
	Expect(err).NotTo(HaveOccurred())

	cmd := exec.Command("cf", "login", "-a", api, "-u", username, "-p", password, "-o", org, "-s", space)
	session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
	Expect(err).ShouldNot(HaveOccurred())
	session.Wait(timeout)

	cmd = exec.Command("cf", "uninstall-plugin", name)
	session, err = gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
	Expect(err).ShouldNot(HaveOccurred())
	session.Wait(timeout)

	cmd = exec.Command("cf", "install-plugin", path, "-f")
	session, err = gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
	Expect(err).ShouldNot(HaveOccurred())
	session.Wait(timeout)
})

var _ = AfterSuite(func() {
	cmd := exec.Command("cf", "uninstall-plugin", name)
	session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
	Expect(err).ShouldNot(HaveOccurred())
	session.Wait(timeout)

	gexec.CleanupBuildArtifacts()
})

func TestAcceptance(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Acceptance Suite")
}

func mustHaveAll(items ...string) {
	for _, s := range items {
		mustHave(s)
	}
}

func mustHave(s string) {
	if s == "" {
		fmt.Println("All env vars must be set")
		os.Exit(1)
	}
}
