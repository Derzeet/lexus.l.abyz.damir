package main

import (
	// New import

	"fmt"
	"net/http"

	"lexus.damir.l.abyx/internal/data"
)

func (app *application) createGunHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Name         string `json:"name"`
		Model        string `json:"model"`
		Caliber      string `json:"caliber"`
		Price        int    `json:"price"`
		Availability bool   `json:"availability"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.errorResponse(w, r, http.StatusBadRequest, err.Error())
		return
	}

	gun := &data.Gun{
		Name:         input.Name,
		Caliber:      input.Caliber,
		Model:        input.Model,
		Price:        input.Price,
		Availability: input.Availability,
	}

	err = app.models.Guns.Insert(gun)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/gun/%d", gun.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"gun": gun}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
