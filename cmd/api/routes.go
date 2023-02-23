package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() *httprouter.Router {
	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)

	router.HandlerFunc(http.MethodGet, "/v1/guns", app.listGunsHandler)
	router.HandlerFunc(http.MethodPost, "/v1/gun", app.createGunHandler)
	router.HandlerFunc(http.MethodGet, "/v1/gun/:id", app.getGunHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/gun/:id", app.updateGunHandler)

	// router.HandlerFunc(http.MethodPost, "/v1/movies", app.createMovieHandler)
	// router.HandlerFunc(http.MethodGet, "/v1/movies/:id", app.showMovieHandler)

	return router
}
