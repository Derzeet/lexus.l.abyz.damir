package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/julienschmidt/httprouter"
)

// again, in the book you have "any" type, but if you use go 1.17 and lower
// you will use interface{} instead of any
type envelope map[string]interface{}

// Retrieve the "id" URL parameter from the current request context, then convert it to
// an integer and return it. If the operation isn't successful, return 0 and an error.
func (app *application) readIDParam(r *http.Request) (int64, error) {
	params := httprouter.ParamsFromContext(r.Context())
	id, err := strconv.ParseInt(params.ByName("id"), 10, 64)
	// if err != nil || id < 1 {
	// 	return 0, errors.New("invalid id parameter")
	// }
	if err != nil {
		return 0, errors.New("invalid id parameter")
	}
	return id, nil
}

func (app *application) readNameParam(r *http.Request) (string, error) {
	params := httprouter.ParamsFromContext(r.Context())
	name := params.ByName("name")
	if name == "" {
		return "", errors.New("Empty name")
	}
	return name, nil
}

func (app *application) readString(qs url.Values, key string, defaultValue string) string {
	s := qs.Get(key)
	if s == "" {
		return defaultValue
	}
	return s
}

// in my version of go there is no type as 'any', and instead of it I used interface{},
// cuz Marshal actually accepts it as a parameter and map is implementing interface.
// on your side data interface{} must be data any if you are using go version 1.18 or newer
// any is a type alias of interface
func (app *application) writeJSON(w http.ResponseWriter, status int, data interface{}, headers http.Header) error {
	js, err := json.Marshal(data)
	if err != nil {
		return err
	}

	js = append(js, '\n')

	//adding additional headers if there are any to be added
	for key, value := range headers {
		w.Header()[key] = value
	}

	// Adding Content-Type and status code to header and response as json
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)
	return nil
}

func (app *application) readCSV(qs url.Values, key string, defaultValue []string) []string {
	// Extract the value from the query string.
	csv := qs.Get(key)
	// If no key exists (or the value is empty) then return the default value.
	if csv == "" {
		return defaultValue
	}
	// Otherwise parse the value into a []string slice and return it.
	return strings.Split(csv, ",")
}

func (app *application) readInt(qs url.Values, key string, defaultValue int) int {
	// Extract the value from the query string.
	s := qs.Get(key)
	// If no key exists (or the value is empty) then return the default value.
	if s == "" {
		return defaultValue
	}
	// Try to convert the value to an int. If this fails, add an error message to the
	// validator instance and return the default value.
	i, err := strconv.Atoi(s)
	if err != nil {
		// v.AddError(key, "must be an integer value")
		return defaultValue
	}
	// Otherwise, return the converted integer value.
	return i
}

func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	err := json.NewDecoder(r.Body).Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError
		var invalidUnmarshalError *json.InvalidUnmarshalError
		if errors.As(err, &syntaxError) {
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", syntaxError.Offset)
		} else if errors.As(err, &unmarshalTypeError) {
			if unmarshalTypeError.Field != "" {
				return fmt.Errorf("body contains incorrect JSON type for field %q", unmarshalTypeError.Field)
			}
			return fmt.Errorf("body contains badly-formed JSON (at character %d)", unmarshalTypeError.Offset)

		} else if errors.As(err, &invalidUnmarshalError) {
			panic(err) //If our program reaches a point where it cannot be recovered due to some major errors

		} else if errors.Is(err, io.ErrUnexpectedEOF) {
			return errors.New("body contains badly-formed JSON")

		} else if errors.Is(err, io.EOF) {
			return errors.New("body must not be empty")

		} else {
			return err
		}
	}

	return nil
}