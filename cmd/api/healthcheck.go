package main

import (
	"net/http"
)

func (app *application) healthcheckHandler(w http.ResponseWriter, r *http.Request) {

	env := envelope{

		"status": "available",
		"system_info": map[string]string{
			"environment": app.config.env,
			"version":     version,
		},
	}
	err := app.writeJSON(w, http.StatusOK, env, nil)
	if err != nil {
		// app.logger.Print(err)
		// http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		app.serverErrorResponse(w, r, err)
	}
}
