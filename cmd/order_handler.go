package main

import (
	"github.com/balgabekj/go-ecommerce/pkg/model"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
)

func (app *application) createOrderHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		CustomerName string `json:"customerName"`
		TotalAmount  int    `json:"totalAmount"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	order := &model.Order{
		CustomerName: input.CustomerName,
		TotalAmount:  input.TotalAmount,
	}

	err = app.models.Orders.Insert(order)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "Failed to create order")
		return
	}

	app.respondWithJson(w, http.StatusCreated, order)
}

func (app *application) getOrderHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["orderId"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid order ID")
		return
	}

	order, err := app.models.Orders.Get(id)
	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "Order not found")
		return
	}

	app.respondWithJson(w, http.StatusOK, order)
}

func (app *application) updateOrderHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["orderId"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid order ID")
		return
	}

	order, err := app.models.Orders.Get(id)
	if err != nil {
		app.respondWithError(w, http.StatusNotFound, "Order not found")
		return
	}

	var input struct {
		CustomerName *string `json:"customerName"`
		TotalAmount  *int    `json:"totalAmount"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	if input.CustomerName != nil {
		order.CustomerName = *input.CustomerName
	}
	if input.TotalAmount != nil {
		order.TotalAmount = *input.TotalAmount
	}

	err = app.models.Orders.Update(order)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "Failed to update order")
		return
	}

	app.respondWithJson(w, http.StatusOK, order)
}

func (app *application) deleteOrderHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	param := vars["orderId"]

	id, err := strconv.Atoi(param)
	if err != nil || id < 1 {
		app.respondWithError(w, http.StatusBadRequest, "Invalid order ID")
		return
	}

	err = app.models.Orders.Delete(id)
	if err != nil {
		app.respondWithError(w, http.StatusInternalServerError, "Failed to delete order")
		return
	}

	app.respondWithJson(w, http.StatusOK, map[string]string{"message": "Order deleted successfully"})
}
