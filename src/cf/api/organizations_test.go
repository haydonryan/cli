package api_test

import (
	. "cf/api"
	"cf/models"
	"cf/net"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"
	mr "github.com/tjarratt/mr_t"
	"net/http"
	"net/http/httptest"
	testapi "testhelpers/api"
	testconfig "testhelpers/configuration"
	testnet "testhelpers/net"
)

var _ = Describe("Testing with ginkgo", func() {
	It("TestOrganizationsListOrgs", func() {
		firstPageOrgsRequest := testapi.NewCloudControllerTestRequest(testnet.TestRequest{
			Method: "GET",
			Path:   "/v2/organizations",
			Response: testnet.TestResponse{Status: http.StatusOK, Body: `{
	"next_url": "/v2/organizations?page=2",
	"resources": [
		{
		  "metadata": { "guid": "org1-guid" },
		  "entity": { "name": "Org1" }
		},
		{
		  "metadata": { "guid": "org2-guid" },
		  "entity": { "name": "Org2" }
		}
	]}`},
		})

		secondPageOrgsRequest := testapi.NewCloudControllerTestRequest(testnet.TestRequest{
			Method: "GET",
			Path:   "/v2/organizations?page=2",
			Response: testnet.TestResponse{Status: http.StatusOK, Body: `{"resources": [
		{
		  "metadata": { "guid": "org3-guid" },
		  "entity": { "name": "Org3" }
		}
	]}`},
		})

		ts, handler, repo := createOrganizationRepo(firstPageOrgsRequest, secondPageOrgsRequest)
		defer ts.Close()

		orgs := []models.Organization{}
		apiResponse := repo.ListOrgs(func(o models.Organization) bool {
			orgs = append(orgs, o)
			return true
		})

		Expect(len(orgs)).To(Equal(3))
		Expect(orgs[0].Guid).To(Equal("org1-guid"))
		Expect(orgs[1].Guid).To(Equal("org2-guid"))
		Expect(orgs[2].Guid).To(Equal("org3-guid"))
		assert.True(mr.T(), apiResponse.IsSuccessful())
		assert.True(mr.T(), handler.AllRequestsCalled())
	})

	It("TestOrganizationsListOrgsWithNoOrgs", func() {
		emptyOrgsRequest := testapi.NewCloudControllerTestRequest(testnet.TestRequest{
			Method:   "GET",
			Path:     "/v2/organizations",
			Response: testnet.TestResponse{Status: http.StatusOK, Body: `{"resources": []}`},
		})

		ts, handler, repo := createOrganizationRepo(emptyOrgsRequest)
		defer ts.Close()

		wasCalled := false
		apiResponse := repo.ListOrgs(func(o models.Organization) bool {
			wasCalled = true
			return false
		})

		assert.False(mr.T(), wasCalled)
		assert.True(mr.T(), apiResponse.IsSuccessful())
		assert.True(mr.T(), handler.AllRequestsCalled())
	})

	It("TestOrganizationsFindByName", func() {
		req := testapi.NewCloudControllerTestRequest(testnet.TestRequest{
			Method: "GET",
			Path:   "/v2/organizations?q=name%3Aorg1&inline-relations-depth=1",
			Response: testnet.TestResponse{Status: http.StatusOK, Body: `{"resources": [{
	  "metadata": { "guid": "org1-guid" },
	  "entity": {
		"name": "Org1",
		"quota_definition": {
		  "entity": {
			"name": "not-your-average-quota",
			"memory_limit": 128
		  }
		},
		"spaces": [{
		  "metadata": { "guid": "space1-guid" },
		  "entity": { "name": "Space1" }
		}],
		"domains": [{
		  "metadata": { "guid": "domain1-guid" },
		  "entity": { "name": "cfapps.io" }
		}]
	  }
	}]}`},
		})

		ts, handler, repo := createOrganizationRepo(req)
		defer ts.Close()
		existingOrg := models.Organization{}
		existingOrg.Guid = "org1-guid"
		existingOrg.Name = "Org1"

		org, apiResponse := repo.FindByName("Org1")
		assert.True(mr.T(), handler.AllRequestsCalled())
		assert.False(mr.T(), apiResponse.IsNotSuccessful())

		Expect(org.Name).To(Equal(existingOrg.Name))
		Expect(org.Guid).To(Equal(existingOrg.Guid))
		Expect(org.QuotaDefinition.Name).To(Equal("not-your-average-quota"))
		Expect(org.QuotaDefinition.MemoryLimit).To(Equal(uint64(128)))
		Expect(len(org.Spaces)).To(Equal(1))
		Expect(org.Spaces[0].Name).To(Equal("Space1"))
		Expect(org.Spaces[0].Guid).To(Equal("space1-guid"))
		Expect(len(org.Domains)).To(Equal(1))
		Expect(org.Domains[0].Name).To(Equal("cfapps.io"))
		Expect(org.Domains[0].Guid).To(Equal("domain1-guid"))
	})

	It("TestOrganizationsFindByNameWhenDoesNotExist", func() {
		req := testapi.NewCloudControllerTestRequest(testnet.TestRequest{
			Method:   "GET",
			Path:     "/v2/organizations?q=name%3Aorg1&inline-relations-depth=1",
			Response: testnet.TestResponse{Status: http.StatusOK, Body: `{"resources": []}`},
		})

		ts, handler, repo := createOrganizationRepo(req)
		defer ts.Close()

		_, apiResponse := repo.FindByName("org1")
		assert.True(mr.T(), handler.AllRequestsCalled())
		assert.False(mr.T(), apiResponse.IsError())
		assert.True(mr.T(), apiResponse.IsNotFound())
	})

	It("TestCreateOrganization", func() {
		req := testapi.NewCloudControllerTestRequest(testnet.TestRequest{
			Method:   "POST",
			Path:     "/v2/organizations",
			Matcher:  testnet.RequestBodyMatcher(`{"name":"my-org"}`),
			Response: testnet.TestResponse{Status: http.StatusCreated},
		})

		ts, handler, repo := createOrganizationRepo(req)
		defer ts.Close()

		apiResponse := repo.Create("my-org")
		assert.True(mr.T(), handler.AllRequestsCalled())
		assert.False(mr.T(), apiResponse.IsNotSuccessful())
	})

	It("TestRenameOrganization", func() {
		req := testapi.NewCloudControllerTestRequest(testnet.TestRequest{
			Method:   "PUT",
			Path:     "/v2/organizations/my-org-guid",
			Matcher:  testnet.RequestBodyMatcher(`{"name":"my-new-org"}`),
			Response: testnet.TestResponse{Status: http.StatusCreated},
		})

		ts, handler, repo := createOrganizationRepo(req)
		defer ts.Close()

		apiResponse := repo.Rename("my-org-guid", "my-new-org")
		assert.True(mr.T(), handler.AllRequestsCalled())
		assert.False(mr.T(), apiResponse.IsNotSuccessful())
	})

	It("TestDeleteOrganization", func() {

		req := testapi.NewCloudControllerTestRequest(testnet.TestRequest{
			Method:   "DELETE",
			Path:     "/v2/organizations/my-org-guid?recursive=true",
			Response: testnet.TestResponse{Status: http.StatusOK},
		})

		ts, handler, repo := createOrganizationRepo(req)
		defer ts.Close()

		apiResponse := repo.Delete("my-org-guid")
		assert.True(mr.T(), handler.AllRequestsCalled())
		assert.False(mr.T(), apiResponse.IsNotSuccessful())
	})
})

func createOrganizationRepo(reqs ...testnet.TestRequest) (ts *httptest.Server, handler *testnet.TestHandler, repo OrganizationRepository) {
	ts, handler = testnet.NewTLSServer(GinkgoT(), reqs)

	configRepo := testconfig.NewRepositoryWithDefaults()
	configRepo.SetApiEndpoint(ts.URL)
	gateway := net.NewCloudControllerGateway()
	repo = NewCloudControllerOrganizationRepository(configRepo, gateway)
	return
}
