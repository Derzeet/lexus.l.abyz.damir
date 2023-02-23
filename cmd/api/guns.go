package main

import (
	// New import

	"errors"
	"fmt"
	"net/http"

	"lexus.damir.l.abyx/internal/data"
	"lexus.damir.l.abyx/validator"
)

func (app *application) createGunHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Name         string `json:"name"`
		Model        string `json:"model"`
		Caliber      string `json:"caliber"`
		Price        int    `json:"price"`
		Availability bool   `json:"availability"`
		Type         string `json:"type"`
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
		Type:         input.Type,
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

func (app *application) getGunHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	gun, err := app.models.Guns.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"gun": gun}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateGunHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	gun, err := app.models.Guns.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	var input struct {
		Name         *string `json:"name"`
		Model        *string `json:"model"`
		Caliber      *string `json:"caliber"`
		Price        *int    `json:"price"`
		Availability *bool   `json:"availability"`
	}
	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	if input.Name != nil {
		gun.Name = *input.Name
	}
	if input.Model != nil {
		gun.Model = *input.Model
	}
	if input.Caliber != nil {
		gun.Caliber = *input.Caliber
	}
	if input.Price != nil {
		gun.Price = *input.Price
	}
	if input.Availability != nil {
		gun.Availability = *input.Availability
	}

	err = app.models.Guns.Update(gun)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrEditConflict):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{"gun": gun}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listGunsHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name string
		Type string
		data.Filters
	}
	// Initialize a new Validator instance.
	v := validator.New()
	// Call r.URL.Query() to get the url.Values map containing the query string data.
	qs := r.URL.Query()
	// Use our helpers to extract the title and genres query string values, falling back
	// to defaults of an empty string and an empty slice respectively if they are not
	// provided by the client.
	input.Name = app.readString(qs, "name", "")
	input.Type = app.readString(qs, "type", "")
	// Get the page and page_size query string values as integers. Notice that we set
	// the default page value to 1 and default page_size to 20, and that we pass the
	// validator instance as the final argument here.
	input.Filters.Page = app.readInt(qs, "page", 1)
	input.Filters.PageSize = app.readInt(qs, "page_size", 20)
	// Extract the sort query string value, falling back to "id" if it is not provided
	// by the client (which will imply a ascending sort on movie ID).
	input.Filters.Sort = app.readString(qs, "sort", "id")

	input.Filters.SortSafelist = []string{"id", "name", "type", "price", "model", "-id", "-name", "-type", "-price", "-model"}
	// Check the Validator instance for any errors and use the failedValidationResponse()
	// helper to send the client a response if necessary.
	if data.ValidateFilters(v, input.Filters); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	guns, err := app.models.Guns.List(input.Name, input.Type, input.Filters)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	// Send a JSON response containing the movie data.
	err = app.writeJSON(w, http.StatusOK, envelope{"guns": guns}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
