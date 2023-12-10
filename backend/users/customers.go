package users

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	_ "github.com/lib/pq"
)

func Registration(w http.ResponseWriter, r *http.Request) {
	connStr := "user=postgres password=123 dbname=usersdb sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	name := r.FormValue("name")
	login := r.FormValue("login")
	password := r.FormValue("password")

	result, err := db.Exec("insert into customers (Name,Login,Password) values (%s,%s,%s)", name, login, password)
	if err != nil {
		fmt.Println("Ошибка при выполнении запроса:", err)
	}
	fmt.Println(result.RowsAffected())
	fmt.Fprintf(w, "%s", "Success")
}

func Authorisation(w http.ResponseWriter, r *http.Request) {
	connStr := "user=postgres password=123 dbname=usersdb sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	login := r.FormValue("login")
	password := r.FormValue("password")

	var result string
	err = db.QueryRow("SELECT * FROM customers WHERE login=$1 AND password=$2", login, password).Scan(&result)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("Логин или пароль неверны")
		} else {
			fmt.Println("Ошибка при выполнении запроса:", err)
		}
		return
	}

	fmt.Fprintf(w, "%s", "Логин и пароль верны")
}

func Add_bank_data(w http.ResponseWriter, r *http.Request) {
	connStr := "user=postgres password=123 dbname=usersdb sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	name := r.FormValue("name")
	bank_number := r.FormValue("bank_number")
	bank_cvc := r.FormValue("bank_cvc")

	result, err := db.Exec("UPDATE customers SET bank_number = $2, bank_cvc = $3 WHERE name = $1", name, bank_number, bank_cvc) //нужно будет сделать либо имя уникальным либо искать по логину
	if err != nil {
		fmt.Println("Ошибка при выполнении запроса:", err)
	}
	fmt.Println(result.RowsAffected())

	fmt.Fprintf(w, "%s", "Success")
}

func Add_to_cart(w http.ResponseWriter, r *http.Request) {
	connStr := "user=postgres password=123 dbname=usersdb sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	name := r.FormValue("name")
	cart := r.FormValue("cart")

	result, err := db.Exec("UPDATE customers SET cart = array_append(cart, $2) WHERE name = $1", name, cart) //нужно будет сделать либо имя уникальным либо искать по логину
	if err != nil {
		fmt.Println("Ошибка при выполнении запроса:", err)
	}
	fmt.Println(result.RowsAffected())

	fmt.Fprintf(w, "%s", "Success")
}

func Check_bank_data(w http.ResponseWriter, r *http.Request) {
	connStr := "user=postgres password=123 dbname=usersdb sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	name := r.FormValue("name")

	rows, err := db.Query("SELECT * FROM customers WHERE name = $1 AND card_number <> '0' AND cvc <> '0'", name)
	if err != nil {
		fmt.Println("Ошибка при выполнении запроса:", err)
	}
	defer rows.Close()

	if rows.Next() {
		// Найдено совпадение
		fmt.Println("У пользователя есть непустой номер карты и cvc")
	} else {
		// Совпадение не найдено
		fmt.Println("У пользователя отсутствует непустой номер карты или cvc")
	}
}

func Sale(w http.ResponseWriter, r *http.Request) {
	connStr := "user=postgres password=123 dbname=usersdb sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	name := r.FormValue("name")

	var result string
	err = db.QueryRow("SELECT total_amount_of_purchases FROM customers WHERE name=$1 ", name).Scan(&result)
	if err != nil {
		fmt.Println("Ошибка при выполнении запроса:", err)
		return
	}

	fmt.Fprintf(w, "%s", result)
}

func Bought(w http.ResponseWriter, r *http.Request) {
	connStr := "user=postgres password=123 dbname=usersdb sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	name := r.FormValue("name")
	sum, _ := strconv.Atoi(r.FormValue("sum"))

	result, err := db.Exec("UPDATE customers SET total_amount_of_purchases = total_amount_of_purchases + $2 WHERE name = $1", name, sum) //нужно будет сделать либо имя уникальным либо искать по логину
	if err != nil {
		fmt.Println("Ошибка при выполнении запроса:", err)
	}
	fmt.Println(result.RowsAffected())

	fmt.Fprintf(w, "%s", "Success")
}
