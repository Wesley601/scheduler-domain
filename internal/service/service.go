package service

import (
	"context"
	"time"

	"alinea.com/internal/core"
	"alinea.com/pkg/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ServiceService struct {
	serviceRepository *mongo.ServiceRepository
}

func NewServiceService(serviceRepository *mongo.ServiceRepository) *ServiceService {
	return &ServiceService{
		serviceRepository: serviceRepository,
	}
}

type CreateServiceDTO struct {
	Name     string
	Duration string
}

func (useCase *ServiceService) Create(c context.Context, dto CreateServiceDTO) (*core.Service, error) {
	d, err := time.ParseDuration(dto.Duration)
	if err != nil {
		return nil, err
	}

	service := core.Service{
		ID:       primitive.NewObjectID().Hex(),
		Name:     dto.Name,
		Duration: d,
	}

	if err := useCase.serviceRepository.Save(c, service); err != nil {
		return nil, err
	}

	return &service, nil
}

func (useCase *ServiceService) FindByID(c context.Context, id string) (*core.Service, error) {
	service, err := useCase.serviceRepository.FindByID(c, id)
	if err != nil {
		return nil, err
	}

	return &service, nil
}

func (useCase *ServiceService) List(c context.Context, page mongo.ListFilter) (mongo.ServicePage, error) {
	services, err := useCase.serviceRepository.List(c, page)
	if err != nil {
		return mongo.ServicePage{}, err
	}

	return services, nil
}
