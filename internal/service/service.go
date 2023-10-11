package service

import (
	"context"
	"time"

	"alinea.com/internal/core"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ServiceRepository interface {
	Save(c context.Context, s core.Service) error
	FindByID(c context.Context, id string) (core.Service, error)
}

type serviceService struct {
	serviceRepository ServiceRepository
}

func NewServiceService(serviceRepository ServiceRepository) *serviceService {
	return &serviceService{
		serviceRepository: serviceRepository,
	}
}

type CreateServiceDTO struct {
	Name     string
	Duration string
}

func (useCase *serviceService) Create(c context.Context, dto CreateServiceDTO) (Parser, error) {
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

func (useCase *serviceService) FindByID(c context.Context, id string) (Parser, error) {
	service, err := useCase.serviceRepository.FindByID(c, id)
	if err != nil {
		return Parser{}, err
	}

	parser := Parser{
		service: service,
	}

	return parser, nil
}
