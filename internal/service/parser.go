package service

import (
	"encoding/json"
	"time"

	"alinea.com/internal/core"
	"alinea.com/pkg/mongo"
	"alinea.com/pkg/utils"
)

type ServiceJSON struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Duration string `json:"duration"`
}

type PageJson struct {
	Total   int64 `json:"total"`
	Page    int   `json:"page"`
	PerPage int   `json:"per_page"`
}

type ListServiceJSON struct {
	Data []ServiceJSON `json:"data"`
	Meta PageJson      `json:"meta"`
}

type ListParser struct {
	services mongo.ServicePage
}

func (p ListParser) FromService(a ListServiceJSON) (ListParser, error) {
	var parser ListParser

	parser.services.Pagination = mongo.Pagination{
		Page:    a.Meta.Page,
		PerPage: a.Meta.PerPage,
	}
	parser.services.Total = a.Meta.Total

	for _, a := range a.Data {
		parser.services.Data = append(parser.services.Data, utils.Must(assembleService(a.ID, a.Name, a.Duration)))
	}

	return parser, nil
}

func (p ListParser) FromJSON(b []byte) (ListParser, error) {
	var a ListServiceJSON

	err := json.Unmarshal(b, &a)
	if err != nil {
		return ListParser{}, err
	}

	return p.FromService(a)
}

func (p *ListParser) ToJSON() ([]byte, error) {
	a, err := p.ToJSONStruct()
	if err != nil {
		return []byte{}, err
	}

	return json.MarshalIndent(a, "", "  ")
}

func (p *ListParser) ToService() ([]core.Service, error) {
	return p.services.Data, nil
}

func (p *ListParser) ToJSONStruct() (ListServiceJSON, error) {
	var json ListServiceJSON
	json.Meta = PageJson{
		Total:   p.services.Total,
		Page:    p.services.Page,
		PerPage: p.services.PerPage,
	}

	for _, service := range p.services.Data {
		json.Data = append(json.Data, ServiceJSON{
			ID:       service.ID,
			Name:     service.Name,
			Duration: service.Duration.String(),
		})
	}

	return json, nil
}

type Parser struct {
	service core.Service
}

func (p Parser) FromJSON(b []byte) (Parser, error) {
	var a ServiceJSON

	err := json.Unmarshal(b, &a)
	if err != nil {
		return Parser{}, err
	}

	return p.FromService(a)
}

func (p Parser) FromService(a ServiceJSON) (Parser, error) {
	var parser Parser

	parser.service = utils.Must(assembleService(a.ID, a.Name, a.Duration))

	return parser, nil
}

func (p *Parser) ToJSON() ([]byte, error) {
	a, err := p.ToJSONStruct()
	if err != nil {
		return []byte{}, err
	}

	return json.MarshalIndent(a, "", "  ")
}

func (p *Parser) ToService() (core.Service, error) {
	return p.service, nil
}

func (p *Parser) ToJSONStruct() (ServiceJSON, error) {
	return ServiceJSON{
		ID:       p.service.ID,
		Name:     p.service.Name,
		Duration: p.service.Duration.String(),
	}, nil
}

func assembleService(id, name, duration string) (core.Service, error) {
	d, err := time.ParseDuration(duration)
	if err != nil {
		return core.Service{}, err
	}
	return core.Service{
		ID:       id,
		Name:     name,
		Duration: d,
	}, nil
}
