package testhelpers

import (
	"cf/configuration"
	"cf"
)

type FakeConfigRepository struct {}

var testConfigurationSingleton *configuration.Configuration
var testConfigurationIsInitialized bool

func (repo FakeConfigRepository) Get() (c configuration.Configuration, err error) {
	if !testConfigurationIsInitialized {
		testConfigurationSingleton = &configuration.Configuration{}
		testConfigurationSingleton.Target = "https://api.run.pivotal.io"
		testConfigurationSingleton.ApiVersion = "2"
		testConfigurationSingleton.AuthorizationEndpoint = "https://login.run.pivotal.io"
		testConfigurationSingleton.ApplicationStartTimeout = 30 // seconds
		testConfigurationIsInitialized = true
	}
	return *testConfigurationSingleton, nil
}

func (repo FakeConfigRepository) Delete() {
	testConfigurationSingleton  = &configuration.Configuration{}
	testConfigurationIsInitialized = false
}

func (repo FakeConfigRepository) Save(config configuration.Configuration) (err error) {
	testConfigurationSingleton = &config
	testConfigurationIsInitialized = true
	return
}

func (repo FakeConfigRepository) ClearSession() (err error) {
	c, _ := repo.Get()
	c.AccessToken = ""
	c.Organization = cf.Organization{}
	c.Space = cf.Space{}
	repo.Save(c)
	return nil
}

func (repo FakeConfigRepository) Login() (c configuration.Configuration) {
	c, _ = repo.Get()
	c.AccessToken = `BEARER eyJhbGciOiJSUzI1NiJ9.eyJqdGkiOiJjNDE4OTllNS1kZTE1LTQ5NGQtYWFiNC04ZmNlYzUxN2UwMDUiLCJzdWIiOiI3NzJkZGEzZi02NjlmLTQyNzYtYjJiZC05MDQ4NmFiZTFmNmYiLCJzY29wZSI6WyJjbG91ZF9jb250cm9sbGVyLnJlYWQiLCJjbG91ZF9jb250cm9sbGVyLndyaXRlIiwib3BlbmlkIiwicGFzc3dvcmQud3JpdGUiXSwiY2xpZW50X2lkIjoiY2YiLCJjaWQiOiJjZiIsImdyYW50X3R5cGUiOiJwYXNzd29yZCIsInVzZXJfaWQiOiI3NzJkZGEzZi02NjlmLTQyNzYtYjJiZC05MDQ4NmFiZTFmNmYiLCJ1c2VyX25hbWUiOiJ1c2VyMUBleGFtcGxlLmNvbSIsImVtYWlsIjoidXNlcjFAZXhhbXBsZS5jb20iLCJpYXQiOjEzNzcwMjgzNTYsImV4cCI6MTM3NzAzNTU1NiwiaXNzIjoiaHR0cHM6Ly91YWEuYXJib3JnbGVuLmNmLWFwcC5jb20vb2F1dGgvdG9rZW4iLCJhdWQiOlsib3BlbmlkIiwiY2xvdWRfY29udHJvbGxlciIsInBhc3N3b3JkIl19.kjFJHi0Qir9kfqi2eyhHy6kdewhicAFu8hrPR1a5AxFvxGB45slKEjuP0_72cM_vEYICgZn3PcUUkHU9wghJO9wjZ6kiIKK1h5f2K9g-Iprv9BbTOWUODu1HoLIvg2TtGsINxcRYy_8LW1RtvQc1b4dBPoopaEH4no-BIzp0E5E`
	repo.Save(c)
	c, _ = repo.Get()
	return c
}
