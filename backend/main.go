package main

import (
	"backend/users"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	// Регистрируем маршруты
	router.HandleFunc("/reg", users.Registration)
	router.HandleFunc("/auth", users.Authorisation)
	router.HandleFunc("/add_bank_data", users.Add_bank_data)
	router.HandleFunc("/add_to_cart", users.Add_to_cart)
	router.HandleFunc("/check_bank_data", users.Check_bank_data)
	router.HandleFunc("/sale", users.Sale)
	router.HandleFunc("/bought", users.Bought)
	http.ListenAndServe(":8080", router)
}
