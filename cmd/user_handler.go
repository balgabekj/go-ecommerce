package main

import (
	"database/sql"
	"encoding/json"
	"github.com/balgabekj/go-ecommerce/pkg/model"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"
)

func (app *application) respondWithError(w http.ResponseWriter, code int, message string) {
	app.respondWithJson(w, code, map[string]string{"error": message})
}

func (app *application) respondWithJson(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)

	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "500 Internal Server Error")
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}

func (app *application) createUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
	}

	hashed, _ := bcrypt.GenerateFromPassword([]byte(input.Password), 8)

	user := &model.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: string(hashed),
	}

	app.respondWithJson(w, http.StatusCreated, user)
}

func (app *application) getAllUsersHandler(w http.ResponseWriter, r *http.Request) {
	users, err := app.models.Users.GetAll()
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "Failed to fetch users")
		return
	}
	app.respondWithJson(w, http.StatusOK, users)
}

func (app *application) getUserByIdHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["userID"])
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	user, err := app.models.Users.GetById(id)
	if err != nil {
		if err == sql.ErrNoRows {
			app.respondWithError(w, http.StatusNotFound, "User not found")
			return
		}
		app.respondWithError(w, http.StatusInternalServerError, "Failed to fetch user")
		return
	}

	app.respondWithJson(w, http.StatusOK, user)
}

func (app *application) updateUserHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["userID"])
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	user, err := app.models.Users.GetById(id)
	if err != nil {
		if err == sql.ErrNoRows {
			app.respondWithError(w, http.StatusNotFound, "User not found")
			return
		}
		app.respondWithError(w, http.StatusInternalServerError, "Failed to fetch user")
		return
	}

	user.Name = input.Name
	user.Email = input.Email
	user.Password = input.Password

	err = app.models.Users.Update(user)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "Failed to update user")
		return
	}

	app.respondWithJson(w, http.StatusOK, user)
}

func (app *application) deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id, err := strconv.Atoi(params["userID"])
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	err = app.models.Users.Delete(id)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "Failed to delete user")
		return
	}

	app.respondWithJson(w, http.StatusOK, map[string]string{"message": "User deleted successfully"})
}
func (app *application) readJSON(w http.ResponseWriter, r *http.Request, dst interface{}) error {
	dec := json.NewDecoder(r.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(dst)
	if err != nil {
		return err
	}

	return nil
}
