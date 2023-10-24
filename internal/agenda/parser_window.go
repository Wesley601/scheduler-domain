package agenda

import (
	"encoding/json"
	"time"

	"alinea.com/internal/core"
	"alinea.com/pkg/utils"
)

type WindowJSON struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type ListWindowJSON struct {
	Data []WindowJSON `json:"data"`
}

type ListWindowParser struct {
	windows []core.Window
}

func (p ListWindowParser) FromWindow(a ListWindowJSON) (ListWindowParser, error) {
	var parser ListWindowParser

	for _, a := range a.Data {
		parser.windows = append(parser.windows, core.Window{
			From: utils.Must(time.Parse(time.RFC3339Nano, a.From)),
			To:   utils.Must(time.Parse(time.RFC3339Nano, a.To)),
		})
	}

	return parser, nil
}

func (p ListWindowParser) FromJSON(b []byte) (ListWindowParser, error) {
	var a ListWindowJSON

	err := json.Unmarshal(b, &a)
	if err != nil {
		return ListWindowParser{}, err
	}

	return p.FromWindow(a)
}

func (p *ListWindowParser) ToJSON() ([]byte, error) {
	a, err := p.ToJSONStruct()
	if err != nil {
		return []byte{}, err
	}

	return json.MarshalIndent(a, "", "  ")
}

func (p *ListWindowParser) ToWindow() ([]core.Window, error) {
	return p.windows, nil
}

func (p *ListWindowParser) ToJSONStruct() ([]WindowJSON, error) {
	var listWindow []WindowJSON

	for _, window := range p.windows {
		listWindow = append(listWindow, WindowJSON{
			From: window.From.Format(time.RFC3339Nano),
			To:   window.To.Format(time.RFC3339Nano),
		})
	}

	return listWindow, nil
}
