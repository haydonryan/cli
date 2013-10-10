package api

import (
	"cf"
	"cf/configuration"
	"cf/net"
	"fmt"
	"strings"
)

type ServiceAuthTokenRepository interface {
	Create(authToken cf.ServiceAuthToken) (apiResponse net.ApiResponse)
	Update(authToken cf.ServiceAuthToken) (apiResponse net.ApiResponse)
	FindAll() (authTokens []cf.ServiceAuthToken, apiResponse net.ApiResponse)
}

type CloudControllerServiceAuthTokenRepository struct {
	gateway net.Gateway
	config  *configuration.Configuration
}

func NewCloudControllerServiceAuthTokenRepository(config *configuration.Configuration, gateway net.Gateway) (repo CloudControllerServiceAuthTokenRepository) {
	repo.gateway = gateway
	repo.config = config
	return
}

func (repo CloudControllerServiceAuthTokenRepository) Create(authToken cf.ServiceAuthToken) (apiResponse net.ApiResponse) {
	path := fmt.Sprintf("%s/v2/service_auth_tokens", repo.config.Target)
	body := fmt.Sprintf(`{"label":"%s","provider":"%s","token":"%s"}`, authToken.Label, authToken.Provider, authToken.Token)

	request, apiResponse := repo.gateway.NewRequest("POST", path, repo.config.AccessToken, strings.NewReader(body))
	if apiResponse.IsNotSuccessful() {
		return
	}

	apiResponse = repo.gateway.PerformRequest(request)
	return
}

func (repo CloudControllerServiceAuthTokenRepository) FindAll() (authTokens []cf.ServiceAuthToken, apiResponse net.ApiResponse) {
	path := fmt.Sprintf("%s/v2/service_auth_tokens", repo.config.Target)

	request, apiResponse := repo.gateway.NewRequest("GET", path, repo.config.AccessToken, nil)
	if apiResponse.IsNotSuccessful() {
		return
	}

	response := &ApiResponse{}
	_, apiResponse = repo.gateway.PerformRequestForJSONResponse(request, response)
	if apiResponse.IsNotSuccessful() {
		return
	}

	for _, resource := range response.Resources {
		authTokens = append(authTokens, cf.ServiceAuthToken{
			Guid:     resource.Metadata.Guid,
			Label:    resource.Entity.Label,
			Provider: resource.Entity.Provider,
		})
	}

	return
}

func (repo CloudControllerServiceAuthTokenRepository) Update(authToken cf.ServiceAuthToken) (apiResponse net.ApiResponse) {
	tokens, apiResponse := repo.FindAll()
	if apiResponse.IsNotSuccessful() {
		return
	}

	i := indexOfToken(tokens, authToken)
	if i == -1 {
		apiResponse = net.NewNotFoundApiResponse("Service Token", authToken.Label)
		return
	}

	tokenGuid := tokens[i].Guid

	path := fmt.Sprintf("%s/v2/service_auth_tokens/%s", repo.config.Target, tokenGuid)
	body := fmt.Sprintf(`{"token":"%s"}`, authToken.Token)
	println(path)
	println(body)
	request, apiResponse := repo.gateway.NewRequest("PUT", path, repo.config.AccessToken, strings.NewReader(body))
	if apiResponse.IsNotSuccessful() {
		return
	}

	apiResponse = repo.gateway.PerformRequest(request)
	return
}

func indexOfToken(tokens []cf.ServiceAuthToken, matcher cf.ServiceAuthToken) int {
	key := matcher.FindByNameKey()
	for i, token := range tokens {
		if token.FindByNameKey() == key {
			return i
		}
	}

	return -1
}
