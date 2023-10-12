package server

import (
	"context"
	"log"
	"net/http"
	"text/template"

	"alinea.com/pkg/mongo"
	"alinea.com/pkg/utils"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
)

var r chi.Router

func init() {
	r = chi.NewRouter()
	r.Use(middleware.Logger)

	client := utils.Must(mongo.NewClient(context.Background()))
	serviceRepository := mongo.NewServiceRepository(client)
	serviceController := NewServiceController(serviceRepository)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		ts, err := template.ParseFiles("./web/templates/base.html", "./web/templates/home.html")
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "internal Server Error", http.StatusInternalServerError)
			return
		}

		ts.ExecuteTemplate(w, "base", nil)
	})

	r.Route("/services", func(r chi.Router) {
		r.Get("/", serviceController.List)
		r.Get("/new", serviceController.New)
		r.Post("/", serviceController.Create)
	})
}
