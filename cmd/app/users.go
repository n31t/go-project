package main

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/n31t/go-project/pkg/model"
	"github.com/n31t/go-project/pkg/validator"
)

func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	user := &model.User{
		Name:      input.Name,
		Email:     input.Email,
		Activated: false,
	}

	err = user.Password.Set(input.Password)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	v := validator.New()

	if model.ValidateUser(v, user); !v.Valid() {
		var errors []string
		for field, err := range v.Errors {
			errors = append(errors, field+": "+err)
		}
		errorMessage := strings.Join(errors, ", ")
		app.respondWithError(w, http.StatusBadRequest, errorMessage)
		return
	}

	err = app.models.Users.Insert(user)
	if err != nil {
		if errors.Is(err, model.ErrDuplicateEmail) {
			app.respondWithError(w, http.StatusBadRequest, "Email already in use")
		} else {
			app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		}
		return
	}

	token, err := app.models.Tokens.New(user.Id, 3*24*time.Hour, model.ScopeActivation)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	var res struct {
		Token *string     `json:"token"`
		User  *model.User `json:"user"`
	}
	res.Token = &token.PlainText
	res.User = user

	app.respondWithJSON(w, http.StatusCreated, res)

}

func (app *application) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		TokenPlainText string `json:"token"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	v := validator.New()

	if model.ValidateTokenPlainText(v, input.TokenPlainText); !v.Valid() {
		var errors []string
		for field, err := range v.Errors {
			errors = append(errors, field+": "+err)
		}
		errorMessage := strings.Join(errors, ", ")
		app.respondWithError(w, http.StatusBadRequest, errorMessage)
		return
	}

	user, err := app.models.Users.GetForToken(model.ScopeActivation, input.TokenPlainText)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrRecordNotFound):
			v.AddError("token", "invalid or expired activation token")
			var errors []string
			for field, err := range v.Errors {
				errors = append(errors, field+": "+err)
			}
			errorMessage := strings.Join(errors, ", ")
			app.respondWithError(w, http.StatusBadRequest, errorMessage)
		default:
			app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		}
		return
	}

	user.Activated = true

	err = app.models.Users.Update(user)
	if err != nil {
		switch {
		case errors.Is(err, model.ErrEditConflict):
			app.respondWithError(w, http.StatusConflict, "edit conflict")

		default:
			app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		}
		return
	}

	err = app.models.Tokens.DeleteAllForUser(model.ScopeActivation, user.Id)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	app.respondWithJSON(w, http.StatusCreated, user)
}
