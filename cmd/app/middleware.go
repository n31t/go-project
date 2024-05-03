package main

import (
	"errors"
	"net/http"
	"strings"

	"github.com/n31t/go-project/pkg/model"
	"github.com/n31t/go-project/pkg/validator"
)

func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Add the "Vary: Authorization" header to the response. This indicates to any
		// caches that the response may vary based on the value of the Authorization
		// header in the request.
		w.Header().Add("Vary", "Authorization")
		authorizationHeader := r.Header.Get("Authorization")

		if authorizationHeader == "" {
			r = app.contextSetUser(r, model.AnonymousUser)
			next.ServeHTTP(w, r)
			return
		}

		headerParts := strings.Split(authorizationHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			app.respondWithError(w, http.StatusUnauthorized, "Invalid or missing authorization token")
			return
		}

		token := headerParts[1]

		v := validator.New()

		if model.ValidateTokenPlainText(v, token); !v.Valid() {
			app.respondWithError(w, http.StatusUnauthorized, "Invalid or missing authorization token")
			return
		}

		user, err := app.models.Users.GetForToken(model.ScopeAuthentication, token)
		if err != nil {
			switch {
			case errors.Is(err, model.ErrRecordNotFound):
				app.respondWithError(w, http.StatusUnauthorized, "Invalid or missing authorization token")
			default:
				app.respondWithError(w, http.StatusMovedPermanently, "500 Internal Server Error")
			}
			return
		}

		r = app.contextSetUser(r, user)

		next.ServeHTTP(w, r)
	})

}

func (app *application) requireAuthenticatedUser(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := app.contextGetUser(r)
		if user.IsAnonymous() {
			app.authenticationRequiredResponse(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// Checks that a user is both authenticated and activated.
func (app *application) requireActivatedUser(next http.HandlerFunc) http.HandlerFunc {
	// Rather than returning this http.HandlerFunc we assign it to the variable fn.
	fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := app.contextGetUser(r)
		// Check that a user is activated.
		if !user.Activated {
			app.inactiveAccountResponse(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
	// Wrap fn with the requireAuthenticatedUser() middleware before returning it.
	return app.requireAuthenticatedUser(fn)
}

func (app *application) requirePermission(code string, next http.HandlerFunc) http.HandlerFunc {
	fn := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := app.contextGetUser(r)
		permissions, err := app.models.Permissions.GetAllForUser(user.Id)

		if err != nil {
			app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
			return
		}

		if !permissions.Include(code) {
			app.notPermittedResponse(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
	return app.requireActivatedUser(fn)

}
