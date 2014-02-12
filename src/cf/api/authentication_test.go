package api_test

import (
	. "cf/api"
	"cf/configuration"
	"cf/net"
	"encoding/base64"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"
	mr "github.com/tjarratt/mr_t"
	"net/http"
	"net/http/httptest"
	testconfig "testhelpers/configuration"
	testnet "testhelpers/net"
)

var _ = Describe("AuthenticationRepository", func() {
	It("TestSuccessfullyLoggingIn", func() {
		deps := setupAuthDependencies(mr.T(), successfulLoginRequest)
		defer teardownAuthDependencies(deps)

		auth := NewUAAAuthenticationRepository(deps.gateway, deps.config)
		apiResponse := auth.Authenticate("foo@example.com", "bar")

		assert.True(mr.T(), deps.handler.AllRequestsCalled())
		assert.False(mr.T(), apiResponse.IsError())
		Expect(deps.config.AuthorizationEndpoint()).To(Equal(deps.ts.URL))
		assert.Equal(mr.T(), deps.config.AccessToken(), "BEARER my_access_token")
		assert.Equal(mr.T(), deps.config.RefreshToken(), "my_refresh_token")
	})

	It("TestUnsuccessfullyLoggingIn", func() {
		deps := setupAuthDependencies(mr.T(), unsuccessfulLoginRequest)
		defer teardownAuthDependencies(deps)

		auth := NewUAAAuthenticationRepository(deps.gateway, deps.config)
		apiResponse := auth.Authenticate("foo@example.com", "oops wrong pass")

		assert.True(mr.T(), deps.handler.AllRequestsCalled())
		assert.True(mr.T(), apiResponse.IsNotSuccessful())
		assert.Equal(mr.T(), apiResponse.Message, "Password is incorrect, please try again.")
		assert.Empty(mr.T(), deps.config.AccessToken())
	})

	It("TestServerErrorLoggingIn", func() {
		deps := setupAuthDependencies(mr.T(), errorLoginRequest)
		defer teardownAuthDependencies(deps)

		auth := NewUAAAuthenticationRepository(deps.gateway, deps.config)
		apiResponse := auth.Authenticate("foo@example.com", "bar")

		assert.True(mr.T(), deps.handler.AllRequestsCalled())
		assert.True(mr.T(), apiResponse.IsError())
		assert.Equal(mr.T(), apiResponse.Message, "Server error, status code: 500, error code: , message: ")
		assert.Empty(mr.T(), deps.config.AccessToken())
	})

	It("TestLoggingInWithErrorMaskedAsSuccess", func() {
		deps := setupAuthDependencies(mr.T(), errorMaskedAsSuccessLoginRequest)
		defer teardownAuthDependencies(deps)

		auth := NewUAAAuthenticationRepository(deps.gateway, deps.config)
		apiResponse := auth.Authenticate("foo@example.com", "bar")

		assert.True(mr.T(), deps.handler.AllRequestsCalled())
		assert.True(mr.T(), apiResponse.IsError())
		assert.Equal(mr.T(), apiResponse.Message, "Authentication Server error: I/O error: uaa.10.244.0.22.xip.io; nested exception is java.net.UnknownHostException: uaa.10.244.0.22.xip.io")
		assert.Empty(mr.T(), deps.config.AccessToken())
	})
})

var authHeaders = http.Header{
	"accept":        {"application/json"},
	"content-type":  {"application/x-www-form-urlencoded"},
	"authorization": {"Basic " + base64.StdEncoding.EncodeToString([]byte("cf:"))},
}

var successfulLoginRequest = testnet.TestRequest{
	Method:  "POST",
	Path:    "/oauth/token",
	Header:  authHeaders,
	Matcher: successfulLoginMatcher,
	Response: testnet.TestResponse{
		Status: http.StatusOK,
		Body: `
{
  "access_token": "my_access_token",
  "token_type": "BEARER",
  "refresh_token": "my_refresh_token",
  "scope": "openid",
  "expires_in": 98765
} `},
}

var successfulLoginMatcher = func(t mr.TestingT, request *http.Request) {
	err := request.ParseForm()
	if err != nil {
		assert.Fail(t, "Failed to parse form: %s", err)
		return
	}

	assert.Equal(t, request.Form.Get("username"), "foo@example.com", "Username did not match.")
	assert.Equal(t, request.Form.Get("password"), "bar", "Password did not match.")
	assert.Equal(t, request.Form.Get("grant_type"), "password", "Grant type did not match.")
	assert.Equal(t, request.Form.Get("scope"), "", "Scope did not mathc.")
}

var unsuccessfulLoginRequest = testnet.TestRequest{
	Method: "POST",
	Path:   "/oauth/token",
	Response: testnet.TestResponse{
		Status: http.StatusUnauthorized,
	},
}

var errorLoginRequest = testnet.TestRequest{
	Method: "POST",
	Path:   "/oauth/token",
	Response: testnet.TestResponse{
		Status: http.StatusInternalServerError,
	},
}

var errorMaskedAsSuccessLoginRequest = testnet.TestRequest{
	Method: "POST",
	Path:   "/oauth/token",
	Response: testnet.TestResponse{
		Status: http.StatusOK,
		Body: `
{"error":{"error":"rest_client_error","error_description":"I/O error: uaa.10.244.0.22.xip.io; nested exception is java.net.UnknownHostException: uaa.10.244.0.22.xip.io"}}
`},
}

type authDependencies struct {
	ts      *httptest.Server
	handler *testnet.TestHandler
	config  configuration.ReadWriter
	gateway net.Gateway
}

func setupAuthDependencies(t mr.TestingT, request testnet.TestRequest) (deps authDependencies) {
	deps.ts, deps.handler = testnet.NewTLSServer(t, []testnet.TestRequest{request})

	deps.config = testconfig.NewRepository()
	deps.config.SetAuthorizationEndpoint(deps.ts.URL)

	deps.gateway = net.NewUAAGateway()
	return
}

func teardownAuthDependencies(deps authDependencies) {
	deps.ts.Close()
}
