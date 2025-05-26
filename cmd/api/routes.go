package main

import (
	"net/http"
	"time"

	"github.com/fajardwntara/vue-api/internal/data"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

func (app *application) routes() http.Handler {
	mux := chi.NewRouter()
	mux.Use(middleware.Recoverer)
	mux.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://*", "https://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	mux.Get("/users/login", app.Login)
	mux.Post("/users/login", app.Login)

	mux.Get("/users/all", func(w http.ResponseWriter, r *http.Request) {
		var users data.User

		data, err := users.GetAll()

		if err != nil {
			app.errorLog.Println(err)
		}

		payload := JsonResponse{
			Error:   false,
			Message: "success",
			Data:    envelope{"user": data},
		}

		app.writeJSON(w, http.StatusOK, payload)

	})

	mux.Get("/users/add", func(w http.ResponseWriter, r *http.Request) {
		var u = data.User{
			Email:     "you@there.com",
			FirstName: "You",
			LastName:  "There",
			Password:  "password",
		}

		app.infoLog.Println("adding user...")

		id, err := app.models.User.Insert(u)

		if err != nil {
			app.errorLog.Println(err)
			app.errorJSON(w, err, http.StatusBadRequest)
			return
		}

		app.infoLog.Println("got back id of", id)

		newUser, _ := app.models.User.GetOne(id)
		app.writeJSON(w, http.StatusOK, newUser)
	})

	mux.Get("/test-generate-token", func(w http.ResponseWriter, r *http.Request) {
		token, err := app.models.User.Token.GenerateToken(2, 60*time.Minute)

		if err != nil {
			app.errorLog.Println(err)
			return
		}

		token.Email = "admin@example.com"
		token.CreatedAt = time.Now()
		token.UpdatedAt = time.Now()

		payload := JsonResponse{
			Error:   false,
			Message: "success",
			Data:    token,
		}

		app.writeJSON(w, http.StatusOK, payload)

	})

	mux.Get("/test-save-token", func(w http.ResponseWriter, r *http.Request) {
		token, err := app.models.User.Token.GenerateToken(2, 60*time.Minute)

		if err != nil {
			app.errorLog.Println(err)
			return
		}

		user, err := app.models.User.GetOne(2)

		if err != nil {
			app.errorLog.Println(err)
		}

		token.UserID = user.ID
		token.CreatedAt = time.Now()
		token.UpdatedAt = time.Now()

		err = token.Insert(*token, *user)

		payload := JsonResponse{
			Error:   false,
			Message: "success",
			Data:    token,
		}

		app.writeJSON(w, http.StatusOK, payload)

	})

	mux.Get("/test-validate-token", func(w http.ResponseWriter, r *http.Request) {
		tokenValidate := r.URL.Query().Get("token")

		valid, err := app.models.Token.ValidToken(tokenValidate)

		if err != nil {
			app.errorJSON(w, err, http.StatusBadRequest)
			return
		}

		var payload JsonResponse

		payload.Error = false
		payload.Data = valid
		payload.Message = "succesfully generated token"

		app.writeJSON(w, http.StatusOK, payload)

	})

	return mux
}
