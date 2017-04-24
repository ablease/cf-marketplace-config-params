package acceptance_test

import (
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
	"github.com/onsi/gomega/gexec"
)

var _ = Describe("Marketplace-v2", func() {
	It("lists the marketplace", func() {
		cmd := exec.Command("cf", "marketplace-v2")
		session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
		session.Wait(timeout)

		Expect(err).ToNot(HaveOccurred())
		Expect(session.ExitCode()).To(BeZero())
		Expect(session).To(gbytes.Say("Getting services from marketplace"))
	})

	It("lists a service", func() {
		cmd := exec.Command("cf", "marketplace-v2", "-s", service)
		session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
		session.Wait(timeout)

		Expect(err).ToNot(HaveOccurred())
		Expect(session.ExitCode()).To(BeZero())
		Expect(session).To(gbytes.Say("Getting service plan information for service %s", service))
	})

	Context("for a plan with a schema", func() {
		It("outputs the schema", func() {
			cmd := exec.Command("cf", "marketplace-v2", "-s", service, "-p", planWithSchema)
			session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
			session.Wait(timeout)

			Expect(err).ToNot(HaveOccurred())
			Expect(session.ExitCode()).To(BeZero())
			Expect(session).To(gbytes.Say("Create Service Configuration Parameters"))
			Expect(session).To(gbytes.Say("\\$schema"))
		})
	})

	Context("for a plan with no schema", func() {
		It("shows a message", func() {
			cmd := exec.Command("cf", "marketplace-v2", "-s", service, "-p", planWithoutSchema)
			session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
			session.Wait(timeout)

			Expect(err).ToNot(HaveOccurred())
			Expect(session.ExitCode()).To(BeZero())
			Expect(session).To(gbytes.Say("Create Service Configuration Parameters"))
			Expect(session).To(gbytes.Say("Not available"))
		})
	})
})
