package environment_variable_groups_test

import (
	"net/http"
	"net/http/httptest"
	"time"

	testapi "github.com/cloudfoundry/cli/cf/api/fakes"
	"github.com/cloudfoundry/cli/cf/configuration"
	"github.com/cloudfoundry/cli/cf/models"
	"github.com/cloudfoundry/cli/cf/net"
	testconfig "github.com/cloudfoundry/cli/testhelpers/configuration"
	testnet "github.com/cloudfoundry/cli/testhelpers/net"

	. "github.com/cloudfoundry/cli/cf/api/environment_variable_groups"
	. "github.com/cloudfoundry/cli/testhelpers/matchers"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("CloudControllerEnvironmentVariableGroupsRepository", func() {
	var (
		testServer  *httptest.Server
		testHandler *testnet.TestHandler
		configRepo  configuration.ReadWriter
		repo        CloudControllerEnvironmentVariableGroupsRepository
	)

	BeforeEach(func() {
		configRepo = testconfig.NewRepositoryWithDefaults()
		gateway := net.NewCloudControllerGateway((configRepo), time.Now)
		repo = NewCloudControllerEnvironmentVariableGroupsRepository(configRepo, gateway)
	})

	AfterEach(func() {
		testServer.Close()
	})

	setupTestServer := func(reqs ...testnet.TestRequest) {
		testServer, testHandler = testnet.NewServer(reqs)
		configRepo.SetApiEndpoint(testServer.URL)
	}

	Describe("ListRunning", func() {
		BeforeEach(func() {
			setupTestServer(listRunningRequest)
		})

		It("lists the environment variables in the running group", func() {
			envVars, err := repo.ListRunning()

			Expect(err).NotTo(HaveOccurred())
			Expect(testHandler).To(HaveAllRequestsCalled())
			Expect(envVars).To(Equal([]models.EnvironmentVariable{
				models.EnvironmentVariable{Name: "abc", Value: "123"},
				models.EnvironmentVariable{Name: "do-re-mi", Value: "fa-sol-la-ti"},
			}))
		})
	})
})

var listRunningRequest = testapi.NewCloudControllerTestRequest(testnet.TestRequest{
	Method: "GET",
	Path:   "/v2/config/environment_variable_groups/running",
	Response: testnet.TestResponse{
		Status: http.StatusOK,
		Body: `{
  "abc": 123,
  "do-re-mi": "fa-sol-la-ti"
}`,
	},
})
