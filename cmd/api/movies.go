package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/sxc/greenlight/internal/data"
)

func (app *application) createMovieHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title   string       `json:"title"`
		Year    int32        `json:"year"`
		Runtime data.Runtime `json:"runtime"`
		Genres  []string     `json:"genres"`
	}
	err := app.readJSON(w, r, &input)

	// err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		// app.errorResponse(w, r, http.StatusBadRequest, err.Error())
		app.badRequestResponse(w, r, err)
		return
	}
	fmt.Fprintf(w, "%+v\n", input)
	// fmt.Fprintln(w, "createMovieHandler")
}

func (app *application) showMovieHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil || id < 1 {
		// http.NotFound(w, r)
		app.notFoundResponse(w, r)
		return
	}

	// Craete a new instance of the Movie struct ,
	// containg the ID we extracted from the URL and some dummy data.
	movie := data.Movie{
		ID:        id,
		CreatedAt: time.Now(),
		Title:     "Casablanca",
		Runtime:   102,
		Genres:    []string{"drama", "romance", "war"},
		Version:   1,
	}

	// Encode the struct as JSON and send the result back to the client.
	// Create an envelope {"movie": movie}
	err = app.writeJSON(w, http.StatusOK, envelope{"movie": movie}, nil)

	if err != nil {
		// app.logger.Print(err)
		// http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		app.serverErrorResponse(w, r, err)
	}
}
