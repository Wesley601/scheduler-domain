package server

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"

	"alinea.com/internal/app"
	"alinea.com/internal/core"
	"alinea.com/internal/service"
	"alinea.com/pkg/mongo"
	"alinea.com/pkg/utils"
)

type SuccessResponse struct {
	Message string `json:"message"`
}

func (s SuccessResponse) ToJSON() ([]byte, error) {
	return json.MarshalIndent(s, "", "  ")
}

type ServiceController struct {
	s service.ServiceService
}

func NewServiceController() *ServiceController {

	return &ServiceController{
		s: *app.ServiceService,
	}
}

func (c ServiceController) List(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles(
		"./web/templates/base.html",
		"./web/templates/services/table-row.html",
		"./web/templates/services/list.html",
		"./web/templates/fragment/pagination.html",
	)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "internal Server Error", http.StatusInternalServerError)
		return
	}

	q := r.URL.Query().Get("q")
	page, err := parseOptionalIntQueryParam(r.URL.Query().Get("page"), 1)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "internal Server Error", http.StatusInternalServerError)
		return
	}
	perPage, err := parseOptionalIntQueryParam(r.URL.Query().Get("per_page"), 10)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "internal Server Error", http.StatusInternalServerError)
		return
	}

	services, err := c.s.List(context.Background(), mongo.ListFilter{
		Q: q,
		Pagination: mongo.Pagination{
			Page:    page,
			PerPage: perPage,
		},
	})
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "internal Server Error", http.StatusInternalServerError)
		return
	}

	hx := r.Header.Get("hx-request")
	tpl := "base"

	if hx == "true" {
		tpl = "content"
		if r.Header.Get("Hx-Trigger-Name") != "" {
			tpl = "service-row"
		}
	}

	j, err := services.ToJSONStruct()
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "internal Server Error", http.StatusInternalServerError)
		return
	}

	ts.ExecuteTemplate(w, tpl, struct {
		Services []core.Service
		Pages    []PageInfo
	}{
		Services: utils.Must(services.ToService()),
		Pages:    genPageInfo(j.Meta.Total, int64(j.Meta.PerPage)),
	})

	fmt.Printf("r.Header.Get(\"Accept\"): %v\n", r.Header.Get("Accept"))

	// parseResponse(w, r.Header.Get("Accept"), tpl, ts, &services)
}

func (c ServiceController) New(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("./web/templates/base.html", "./web/templates/services/new.html")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "internal Server Error", http.StatusInternalServerError)
		return
	}

	hx := r.Header.Get("hx-request")
	tpl := "base"

	if hx == "true" {
		tpl = "content"
	}

	ts.ExecuteTemplate(w, tpl, nil)
}

func (c ServiceController) Create(w http.ResponseWriter, r *http.Request) {
	ts, err := template.ParseFiles("./web/templates/fragment/success.html", "./web/templates/fragment/danger.html")
	if err != nil {
		log.Println(err.Error())
		http.Error(w, "internal Server Error", http.StatusInternalServerError)
		return
	}

	n := r.FormValue("name")
	d := r.FormValue("duration")

	_, err = c.s.Create(context.Background(), service.CreateServiceDTO{
		Name:     n,
		Duration: d,
	})
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusBadRequest)
		err := ts.ExecuteTemplate(w, "danger", struct {
			Message string
		}{
			Message: "Oops! Something get wrong",
		})
		if err != nil {
			log.Panic(err)
		}
		return
	}

	fmt.Printf("r.Header.Get(\"Accept\"): %v\n", r.Header.Get("Accept"))

	parseResponse(w, r.Header.Get("Accept"), "success", ts, SuccessResponse{
		Message: "Service created!",
	})
}
