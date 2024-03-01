package main

import (
	"github.com/balgabekj/go-ecommerce/pkg/model"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func (app *application) createProductHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name        string `json:"name"`
		Price       int    `json:"price"`
		Description string `json:"description"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	product := &model.Product{
		Name:        input.Name,
		Price:       input.Price,
		Description: input.Description,
	}

	err = app.models.Products.Insert(product)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "Failed to create product")
		return
	}

	app.respondWithJson(w, http.StatusCreated, product)
}

func (app *application) getProductHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["productId"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	product, err := app.models.Products.Get(id)
	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "Product not found")
		return
	}

	app.respondWithJson(w, http.StatusOK, product)
}

func (app *application) updateProductHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["productId"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	product, err := app.models.Products.Get(id)
	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "Product not found")
		return
	}

	var input struct {
		Name        *string `json:"name"`
		Price       *int    `json:"price"`
		Description *string `json:"description"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if input.Name != nil {
		product.Name = *input.Name
	}
	if input.Price != nil {
		product.Price = *input.Price
	}
	if input.Description != nil {
		product.Description = *input.Description
	}

	err = app.models.Products.Update(product)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "Failed to update product")
		return
	}

	app.respondWithJson(w, http.StatusOK, product)
}

func (app *application) deleteProductHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["productId"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	err = app.models.Products.Delete(id)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "Failed to delete product")
		return
	}

	app.respondWithJson(w, http.StatusOK, map[string]string{"message": "Product deleted successfully"})
}
