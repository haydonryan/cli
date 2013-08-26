package api

import (
	"cf"
	"cf/configuration"
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

var multipleOrgEndpoint = func(writer http.ResponseWriter, request *http.Request) {
	acceptHeaderMatches := request.Header.Get("accept") == "application/json"
	methodMatches := request.Method == "GET"
	pathMatches := request.URL.Path == "/v2/organizations"
	authMatches := request.Header.Get("authorization") == "BEARER my_access_token"

	if !(acceptHeaderMatches && methodMatches && pathMatches && authMatches) {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonResponse := `
{
  "total_results": 2,
  "total_pages": 1,
  "prev_url": null,
  "next_url": null,
  "resources": [
    {
      "metadata": {
        "guid": "org1-guid"
      },
      "entity": {
        "name": "Org1"
      }
    },
    {
      "metadata": {
        "guid": "org2-guid"
      },
      "entity": {
        "name": "Org2"
      }
    }
  ]
}`
	fmt.Fprintln(writer, jsonResponse)
}

func TestOrganizationsFindAll(t *testing.T) {
	ts := httptest.NewTLSServer(http.HandlerFunc(multipleOrgEndpoint))
	defer ts.Close()

	repo := CloudControllerOrganizationRepository{}

	config := &configuration.Configuration{AccessToken: "BEARER my_access_token", Target: ts.URL}
	organizations, err := repo.FindAll(config)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(organizations))

	firstOrg := organizations[0]
	assert.Equal(t, firstOrg.Name, "Org1")
	assert.Equal(t, firstOrg.Guid, "org1-guid")
	secondOrg := organizations[1]
	assert.Equal(t, secondOrg.Name, "Org2")
	assert.Equal(t, secondOrg.Guid, "org2-guid")
}

func TestOrganizationsFindAllWithIncorrectToken(t *testing.T) {
	ts := httptest.NewTLSServer(http.HandlerFunc(multipleOrgEndpoint))
	defer ts.Close()

	repo := CloudControllerOrganizationRepository{}

	config := &configuration.Configuration{AccessToken: "BEARER incorrect_access_token", Target: ts.URL}
	organizations, err := repo.FindAll(config)

	assert.Error(t, err)
	assert.Equal(t, 0, len(organizations))
}

func TestOrganizationsFindByName(t *testing.T) {
	ts := httptest.NewTLSServer(http.HandlerFunc(multipleOrgEndpoint))
	defer ts.Close()

	repo := CloudControllerOrganizationRepository{}
	config := &configuration.Configuration{AccessToken: "BEARER my_access_token", Target: ts.URL}
	existingOrg := cf.Organization{Guid: "org1-guid", Name: "Org1"}

	org, err := repo.FindByName(config, "Org1")
	assert.NoError(t, err)
	assert.Equal(t, org, existingOrg)

	org, err = repo.FindByName(config, "org1")
	assert.NoError(t, err)
	assert.Equal(t, org, existingOrg)

	org, err = repo.FindByName(config, "org that does not exist")
	assert.Error(t, err)
}
