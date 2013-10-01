package testhelpers

import (
	"cf"
	"cf/net"
)

type FakeRouteRepository struct {
	FindByHostHost       string
	FindByHostErr        bool
	FindByHostRoute      cf.Route

	CreatedRoute       cf.Route
	CreatedRouteDomain cf.Domain

	BoundRoute cf.Route
	BoundApp   cf.Application

	FindAllErr    bool
	FindAllRoutes []cf.Route
}

func (repo *FakeRouteRepository) FindAll() (routes []cf.Route, apiStatus net.ApiStatus) {
	if repo.FindAllErr {
		apiStatus = net.NewApiStatusWithMessage("Error finding all routes")
	}

	routes = repo.FindAllRoutes
	return
}

func (repo *FakeRouteRepository) FindByHost(host string) (route cf.Route, apiStatus net.ApiStatus) {
	repo.FindByHostHost = host

	if repo.FindByHostErr {
		apiStatus = net.NewApiStatusWithMessage("Route not found")
	}

	route = repo.FindByHostRoute
	return
}

func (repo *FakeRouteRepository) Create(newRoute cf.Route, domain cf.Domain) (createdRoute cf.Route, apiStatus net.ApiStatus) {
	repo.CreatedRoute = newRoute
	repo.CreatedRouteDomain = domain

	createdRoute = cf.Route{
		Host: newRoute.Host,
		Guid: newRoute.Host + "-guid",
	}
	return
}

func (repo *FakeRouteRepository) Bind(route cf.Route, app cf.Application) (apiStatus net.ApiStatus) {
	repo.BoundRoute = route
	repo.BoundApp = app
	return
}



