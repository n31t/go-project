package main

import (
	"errors"
	"net/http"
	"time"

	"github.com/n31t/go-project/pkg/model"
	"github.com/n31t/go-project/pkg/validator"
)

func (app *application) createAuthenticationTokenHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	v := validator.New()
	model.ValidateEmail(v, input.Email)
	model.ValidatePasswordPlaintext(v, input.Password)

	if !v.Valid() {
		app.respondWithError(w, http.StatusUnprocessableEntity, "Invalid request payload")
		return
	}

	user, err := app.models.Users.GetByEmail(input.Email)

	if err != nil {
		switch {
		case errors.Is(err, model.ErrRecordNotFound):
			app.respondWithError(w, http.StatusUnauthorized, "Invalid credentials")
		default:
			app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		}
		return
	}

	match, err := user.Password.Matches(input.Password)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	if !match {
		app.respondWithError(w, http.StatusUnauthorized, "Invalid credentials")
		return
	}

	token, err := app.models.Tokens.New(user.Id, 24*time.Hour, model.ScopeAuthentication)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	var output struct {
		Token   string `json:"token"`
		Expires string `json:"expires"`
	}

	output.Token = token.PlainText
	output.Expires = token.ExpiresAt.Format(time.RFC3339)

	app.respondWithJSON(w, http.StatusCreated, output)
}
