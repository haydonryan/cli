package organization_test

import (
	"cf/commands/organization"
	"cf/configuration"
	"cf/models"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	testapi "testhelpers/api"
	testassert "testhelpers/assert"
	testcmd "testhelpers/commands"
	testconfig "testhelpers/configuration"
	testreq "testhelpers/requirements"
	testterm "testhelpers/terminal"
)

var _ = Describe("org command", func() {
	var (
		ui                  *testterm.FakeUI
		orgRepo             *testapi.FakeOrgRepository
		configRepo          configuration.ReadWriter
		requirementsFactory *testreq.FakeReqFactory
	)

	runCommand := func() {
		cmd := organization.NewListOrgs(ui, configRepo, orgRepo)
		testcmd.RunCommand(cmd, testcmd.NewContext("orgs", []string{}), requirementsFactory)
	}

	BeforeEach(func() {
		ui = &testterm.FakeUI{}
		configRepo = testconfig.NewRepositoryWithDefaults()
		orgRepo = &testapi.FakeOrgRepository{}
		requirementsFactory = &testreq.FakeReqFactory{LoginSuccess: true}
	})

	Describe("requirements", func() {
		It("fails when not logged in", func() {
			requirementsFactory.LoginSuccess = false
			runCommand()
			Expect(testcmd.CommandDidPassRequirements).To(BeFalse())
		})
	})

	Context("when there are orgs to be listed", func() {
		BeforeEach(func() {
			org1 := models.Organization{}
			org1.Name = "Organization-1"

			org2 := models.Organization{}
			org2.Name = "Organization-2"

			org3 := models.Organization{}
			org3.Name = "Organization-3"

			orgRepo.Organizations = []models.Organization{org1, org2, org3}
		})

		It("lists orgs", func() {
			runCommand()

			testassert.SliceContains(ui.Outputs, testassert.Lines{
				{"Getting orgs as my-user"},
				{"Organization-1"},
				{"Organization-2"},
				{"Organization-3"},
			})
		})
	})

	It("tells the user when no orgs were found", func() {
		runCommand()

		testassert.SliceContains(ui.Outputs, testassert.Lines{
			{"Getting orgs as my-user"},
			{"No orgs found"},
		})
	})
})
