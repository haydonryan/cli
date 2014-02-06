package api

import (
	"cf"
	"cf/configuration"
	"cf/net"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
)

type AppRouteEntity struct {
	Host   string
	Domain Resource
}

type AppRouteResource struct {
	Resource
	Entity AppRouteEntity
}

func (resource AppRouteResource) ToFields() (route cf.RouteFields) {
	route.Guid = resource.Metadata.Guid
	route.Host = resource.Entity.Host
	return
}

func (resource AppRouteResource) ToModel() (route cf.RouteSummary) {
	route.RouteFields = resource.ToFields()
	route.Domain.Guid = resource.Entity.Domain.Metadata.Guid
	route.Domain.Name = resource.Entity.Domain.Entity.Name
	return
}

type ApplicationEntity struct {
	Name               *string             `json:"name,omitempty"`
	Command            *string             `json:"command,omitempty"`
	State              *string             `json:"state,omitempty"`
	SpaceGuid          *string             `json:"space_guid,omitempty"`
	Instances          *int                `json:"instances,omitempty"`
	Memory             *uint64             `json:"memory,omitempty"`
	StackGuid          *string             `json:"stack_guid,omitempty"`
	Stack              *StackResource      `json:"stack,omitempty"`
	Routes             *[]AppRouteResource `json:"routes,omitempty"`
	Buildpack          *string             `json:"buildpack,omitempty"`
	EnvironmentJson    *map[string]string  `json:"environment_json,omitempty"`
	HealthCheckTimeout *int                `json:"health_check_timeout,omitempty"`
}

type ApplicationResource struct {
	Resource
	Entity ApplicationEntity
}

func NewApplicationEntityFromAppParams(app cf.AppParams) ApplicationEntity {
	entity := ApplicationEntity{
		Buildpack:          app.BuildpackUrl,
		Name:               app.Name,
		SpaceGuid:          app.SpaceGuid,
		Instances:          app.InstanceCount,
		Memory:             app.Memory,
		StackGuid:          app.StackGuid,
		Command:            app.Command,
		HealthCheckTimeout: app.HealthCheckTimeout,
	}
	if app.State != nil {
		state := strings.ToUpper(*app.State)
		entity.State = &state
	}
	if app.EnvironmentVars != nil && len(*app.EnvironmentVars) > 0 {
		entity.EnvironmentJson = app.EnvironmentVars
	}
	return entity
}

func (resource ApplicationResource) ToFields() (app cf.ApplicationFields) {
	entity := resource.Entity
	app.Guid = resource.Metadata.Guid

	if entity.Name != nil {
		app.Name = *entity.Name
	}
	if entity.Memory != nil {
		app.Memory = uint64(*entity.Memory)
	}
	if entity.Instances != nil {
		app.InstanceCount = *entity.Instances
	}
	if entity.State != nil {
		app.State = strings.ToLower(*entity.State)
	}
	if entity.EnvironmentJson != nil {
		app.EnvironmentVars = *entity.EnvironmentJson
	}
	if entity.SpaceGuid != nil {
		app.SpaceGuid = *entity.SpaceGuid
	}
	return
}

func (resource ApplicationResource) ToModel() (app cf.Application) {
	app.ApplicationFields = resource.ToFields()

	entity := resource.Entity
	if entity.Stack != nil {
		app.Stack = entity.Stack.ToFields()
	}

	if entity.Routes != nil {
		for _, routeResource := range *entity.Routes {
			app.Routes = append(app.Routes, routeResource.ToModel())
		}
	}

	return
}

type PaginatedApplicationResources struct {
	Resources []ApplicationResource
}

type ApplicationRepository interface {
	Create(params cf.AppParams) (createdApp cf.Application, apiResponse net.ApiResponse)
	Read(name string) (app cf.Application, apiResponse net.ApiResponse)
	Update(appGuid string, params cf.AppParams) (updatedApp cf.Application, apiResponse net.ApiResponse)
	Delete(appGuid string) (apiResponse net.ApiResponse)
}

type CloudControllerApplicationRepository struct {
	config  *configuration.Configuration
	gateway net.Gateway
}

func NewCloudControllerApplicationRepository(config *configuration.Configuration, gateway net.Gateway) (repo CloudControllerApplicationRepository) {
	repo.config = config
	repo.gateway = gateway
	return
}

func (repo CloudControllerApplicationRepository) Create(params cf.AppParams) (createdApp cf.Application, apiResponse net.ApiResponse) {
	data, err := repo.formatAppJSON(params)
	if err != nil {
		apiResponse = net.NewApiResponseWithError("Failed to marshal JSON", err)
		return
	}

	path := fmt.Sprintf("%s/v2/apps", repo.config.Target)
	resource := new(ApplicationResource)
	apiResponse = repo.gateway.CreateResourceForResponse(path, repo.config.AccessToken, strings.NewReader(data), resource)
	if apiResponse.IsNotSuccessful() {
		return
	}

	createdApp = resource.ToModel()
	return
}

func (repo CloudControllerApplicationRepository) Read(name string) (app cf.Application, apiResponse net.ApiResponse) {
	path := fmt.Sprintf("%s/v2/spaces/%s/apps?q=%s&inline-relations-depth=1", repo.config.Target, repo.config.SpaceFields.Guid, url.QueryEscape("name:"+name))
	appResources := new(PaginatedApplicationResources)
	apiResponse = repo.gateway.GetResource(path, repo.config.AccessToken, appResources)
	if apiResponse.IsNotSuccessful() {
		return
	}

	if len(appResources.Resources) == 0 {
		apiResponse = net.NewNotFoundApiResponse("%s %s not found", "App", name)
		return
	}

	res := appResources.Resources[0]
	app = res.ToModel()
	return
}

func (repo CloudControllerApplicationRepository) Update(appGuid string, params cf.AppParams) (updatedApp cf.Application, apiResponse net.ApiResponse) {
	data, err := repo.formatAppJSON(params)
	if err != nil {
		apiResponse = net.NewApiResponseWithError("Failed to marshal JSON", err)
		return
	}

	path := fmt.Sprintf("%s/v2/apps/%s?inline-relations-depth=1", repo.config.Target, appGuid)
	resource := new(ApplicationResource)
	apiResponse = repo.gateway.UpdateResourceForResponse(path, repo.config.AccessToken, strings.NewReader(data), resource)
	if apiResponse.IsNotSuccessful() {
		return
	}

	updatedApp = resource.ToModel()
	return
}

func (repo CloudControllerApplicationRepository) formatAppJSON(input cf.AppParams) (data string, err error) {
	appResource := NewApplicationEntityFromAppParams(input)
	bytes, err := json.Marshal(appResource)
	data = string(bytes)
	return
}

func (repo CloudControllerApplicationRepository) Delete(appGuid string) (apiResponse net.ApiResponse) {
	path := fmt.Sprintf("%s/v2/apps/%s?recursive=true", repo.config.Target, appGuid)
	return repo.gateway.DeleteResource(path, repo.config.AccessToken)
}
