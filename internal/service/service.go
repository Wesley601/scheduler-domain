package service

import (
	"context"
	"time"

	"alinea.com/internal/core"
	"alinea.com/pkg/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ServiceRepository interface {
	Save(c context.Context, s core.Service) error
	List(c context.Context, page mongo.ListFilter) (mongo.ServicePage, error)
	FindByID(c context.Context, id string) (core.Service, error)
}

type ServiceService struct {
	serviceRepository ServiceRepository
}

func NewServiceService(serviceRepository ServiceRepository) *ServiceService {
	return &ServiceService{
		serviceRepository: serviceRepository,
	}
}

type CreateServiceDTO struct {
	Name     string
	Duration string
}

func (useCase *ServiceService) Create(c context.Context, dto CreateServiceDTO) (Parser, error) {
	d, err := time.ParseDuration(dto.Duration)
	if err != nil {
		return Parser{}, err
	}

	service := core.Service{
		ID:       primitive.NewObjectID().Hex(),
		Name:     dto.Name,
		Duration: d,
	}

	if err := useCase.serviceRepository.Save(c, service); err != nil {
		return Parser{}, err
	}

	parser := Parser{
		service: service,
	}

	return parser, nil
}

func (useCase *ServiceService) FindByID(c context.Context, id string) (Parser, error) {
	service, err := useCase.serviceRepository.FindByID(c, id)
	if err != nil {
		return Parser{}, err
	}

	parser := Parser{
		service: service,
	}

	return parser, nil
}

func (useCase *ServiceService) List(c context.Context, page mongo.ListFilter) (ListParser, error) {
	services, err := useCase.serviceRepository.List(c, page)
	if err != nil {
		return ListParser{}, err
	}

	parser := ListParser{
		services: services,
	}

	return parser, nil
}
