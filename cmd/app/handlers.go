package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/n31t/go-project/pkg/model"
	"github.com/n31t/go-project/pkg/validator"
)

// Anime handlers
func (app *application) animeCreate(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title       string  `json:"title"`
		Episodes    int     `json:"episodes"`
		Studio      string  `json:"studio"`
		Description string  `json:"description"`
		ReleaseYear int     `json:"releaseYear"`
		Genre       string  `json:"genre"`
		Rating      float64 `json:"rating"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	anime := &model.Anime{
		Title:       input.Title,
		Episodes:    input.Episodes,
		Studio:      input.Studio,
		Description: input.Description,
		ReleaseYear: input.ReleaseYear,
		Genre:       input.Genre,
		Rating:      input.Rating,
	}

	err = app.models.Animes.Insert(anime)
	if err != nil {
		// app.respondWithError(w, http.StatusInternalServerError, "TEST")
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}
	app.respondWithJSON(w, http.StatusCreated, anime)
}

//	func (app *application) animesList(w http.ResponseWriter, r *http.Request) {
//		animes, err := app.models.Animes.SelectAll()
//		if err != nil {
//			app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
//			return
//		}
//		app.respondWithJSON(w, http.StatusOK, animes)
//	}
func (app *application) animesList(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Title string
		Genre string
		model.Filters
	}

	v := validator.New()
	qs := r.URL.Query()

	input.Title = app.readString(qs, "title", "")
	input.Genre = app.readString(qs, "genre", "")

	input.Filters.Page = app.readInt(qs, "page", 1, v)
	input.Filters.PageSize = app.readInt(qs, "page_size", 10, v)

	input.Filters.Sort = app.readString(qs, "sort", "id")
	input.Filters.SortSafeList = []string{"id", "title", "episodes", "studio", "releaseYear", "genre", "rating"}

	if model.ValidateFilters(v, input.Filters); !v.Valid() {
		app.respondWithError(w, http.StatusBadRequest, fmt.Sprintf("%v", v.Errors))
		return
	}

	animes, metadata, err := app.models.Animes.SelectAll(input.Title, input.Genre, input.Filters)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}
	fmt.Println(metadata)
	app.respondWithJSONMetadata(w, http.StatusOK, animes, metadata)

}

func (app *application) animeRetrieve(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["id"]

	id, err := strconv.Atoi(param)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid anime ID")
		return
	}

	anime, err := app.models.Animes.Select(id)
	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "Anime not found")
		return
	}
	app.respondWithJSON(w, http.StatusOK, anime)
}

func (app *application) animeUpdate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["id"]

	id, err := strconv.Atoi(param)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid anime ID")
		return
	}
	anime, err := app.models.Animes.Select(id)
	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "Anime not found")
		return
	}

	var input struct {
		Title       *string  `json:"title"`
		Episodes    *int     `json:"episodes"`
		Studio      *string  `json:"studio"`
		Description *string  `json:"description"`
		ReleaseYear *int     `json:"releaseYear"`
		Genre       *string  `json:"genre"`
		Rating      *float64 `json:"rating"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if input.Title != nil {
		anime.Title = *input.Title
	}
	if input.Episodes != nil {
		anime.Episodes = *input.Episodes
	}
	if input.Studio != nil {
		anime.Studio = *input.Studio
	}
	if input.Description != nil {
		anime.Description = *input.Description
	}
	if input.ReleaseYear != nil {
		anime.ReleaseYear = *input.ReleaseYear
	}
	if input.Genre != nil {
		anime.Genre = *input.Genre
	}
	if input.Rating != nil {
		anime.Rating = *input.Rating
	}
	err = app.models.Animes.Update(anime)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}
	app.respondWithJSON(w, http.StatusOK, anime)
}

func (app *application) animeDelete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["id"]
	id, err := strconv.Atoi(param)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid anime ID")
		return
	}
	err = app.models.Animes.Delete(id)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}
	app.respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

// User handlers
// func (app *application) userCreate(w http.ResponseWriter, r *http.Request) {
// 	var input struct {
// 		Username string `json:"username"`
// 		Password string `json:"password"`
// 		Email    string `json:"email"`
// 	}

// 	err := app.readJSON(w, r, &input)
// 	if err != nil {
// 		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
// 		return
// 	}
// 	user := &model.User{
// 		Username: input.Username,
// 		Password: input.Password,
// 		Email:    input.Email,
// 	}

// 	err = app.models.Users.Insert(user)
// 	if err != nil {
// 		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
// 		return
// 	}
// 	app.respondWithJSON(w, http.StatusCreated, user)
// }
