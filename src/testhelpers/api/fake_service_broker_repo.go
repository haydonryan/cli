package api

import (
	"cf"
	"cf/net"
)

type FakeServiceBrokerRepo struct {
	FindByNameName string
	FindByNameServiceBroker cf.ServiceBroker
	FindByNameNotFound bool

	CreatedServiceBroker cf.ServiceBroker
	UpdatedServiceBroker cf.ServiceBroker
	RenamedServiceBroker cf.ServiceBroker
	DeletedServiceBroker cf.ServiceBroker

	FindAllServiceBrokers []cf.ServiceBroker
	FindAllErr bool
}

func (repo *FakeServiceBrokerRepo) FindByName(name string) (serviceBroker cf.ServiceBroker, apiResponse net.ApiResponse) {
	repo.FindByNameName = name
	serviceBroker = repo.FindByNameServiceBroker

	if repo.FindByNameNotFound {
		apiResponse = net.NewNotFoundApiResponse("%s %s not found","Service Broker", name)
	}

	return
}

func (repo *FakeServiceBrokerRepo) FindAll() (serviceBrokers []cf.ServiceBroker, apiResponse net.ApiResponse) {
	if repo.FindAllErr {
		apiResponse = net.NewApiResponseWithMessage("Error finding all service brokers")
	}

	serviceBrokers = repo.FindAllServiceBrokers
	return
}

func (repo *FakeServiceBrokerRepo) Create(serviceBroker cf.ServiceBroker) (apiResponse net.ApiResponse) {
	repo.CreatedServiceBroker = serviceBroker
	return
}

func (repo *FakeServiceBrokerRepo) Update(serviceBroker cf.ServiceBroker) (apiResponse net.ApiResponse) {
	repo.UpdatedServiceBroker = serviceBroker
	return
}

func (repo *FakeServiceBrokerRepo) Rename(serviceBroker cf.ServiceBroker) (apiResponse net.ApiResponse) {
	repo.RenamedServiceBroker = serviceBroker
	return
}

func (repo *FakeServiceBrokerRepo) Delete(serviceBroker cf.ServiceBroker) (apiResponse net.ApiResponse) {
	repo.DeletedServiceBroker = serviceBroker
	return
}
