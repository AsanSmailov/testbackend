package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"

	_ "github.com/lib/pq"
)

type Database struct {
	db *sql.DB
}

func NewDatabase() *Database {
	connStr := "user= password= dbname= sslmode="
	db, err := sql.Open("", connStr)
	if err != nil {
		log.Fatal(err)
	}

	return &Database{
		db: db,
	}
}

func (d *Database) Close() {
	d.db.Close()
}

func (d *Database) ExecQuery(query string) error {
	// Execute the query
	_, err := d.db.Exec(query)
	if err != nil {
		// Check if the error is pq.ErrNoRows
		if err == sql.ErrNoRows {
			return err
		}
		return fmt.Errorf("Error executing query: %s", err)
	}

	return nil
}

func (d *Database) QueryRow(query string) (string, error) {
	var result string
	err := d.db.QueryRow(query).Scan(&result)
	if err != nil {
		// Check if the error is pq.ErrNoRows
		if err == sql.ErrNoRows {
			return "", err
		}
		return "", fmt.Errorf("Error executing query: %s", err)
	}

	return result, nil
}

type UserService struct {
	db *Database
}

func NewUserService(db *Database) *UserService {
	return &UserService{
		db: db,
	}
}

func (us *UserService) Registration(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	login := r.FormValue("login")
	password := r.FormValue("password")

	query := fmt.Sprintf("INSERT INTO users (name, login, password) VALUES ('%s', '%s', '%s')", name, login, password)
	err := us.db.ExecQuery(query)

	if err != nil {
		fmt.Fprintf(w, "%s", "Error executing query")
		return
	}

	fmt.Fprintf(w, "%s", "Success")
}

func (us *UserService) Authorisation(w http.ResponseWriter, r *http.Request) {
	login := r.FormValue("login")
	password := r.FormValue("password")

	query := fmt.Sprintf("SELECT * FROM users WHERE login='%s' AND password='%s'", login, password)
	_, err := us.db.QueryRow(query)
	if err != nil {
		// Check if the error is pq.ErrNoRows
		if err == sql.ErrNoRows {
			fmt.Fprintf(w, "%s", "Invalid login or password")
		} else {
			fmt.Fprintf(w, "%s", "Error executing query")
		}
		return
	}

	fmt.Fprintf(w, "%s", "Valid login and password")
}

func (us *UserService) AddToCart(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	goods := r.FormValue("goods")

	query := fmt.Sprintf("UPDATE users SET goods = array_append(goods, '%s') WHERE name='%s'", goods, name)
	err := us.db.ExecQuery(query)
	if err != nil {
		fmt.Fprintf(w, "%s", "Error executing query")
		return
	}

	fmt.Fprintf(w, "%s", "Success")
}

func (us *UserService) Sale(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")

	query := fmt.Sprintf("SELECT total_sum FROM users WHERE name='%s'", name)
	result, err := us.db.QueryRow(query)
	if err != nil {
		fmt.Fprintf(w, "%s", "Error executing query")
		return
	}

	fmt.Fprintf(w, "%s", result)
}

func (us *UserService) Bought(w http.ResponseWriter, r *http.Request) {
	name := r.FormValue("name")
	sum, _ := strconv.Atoi(r.FormValue("sum"))

	query := fmt.Sprintf("UPDATE users SET total_sum = total_sum + %d WHERE name='%s'", sum, name)
	err := us.db.ExecQuery(query)
	if err != nil {
		fmt.Fprintf(w, "%s", "Error executing query")
		return
	}

	fmt.Fprintf(w, "%s", "Success")
}

func main() {
	db := NewDatabase()
	defer db.Close()

	userService := NewUserService(db)

	http.HandleFunc("/registration", userService.Registration)
	http.HandleFunc("/authorisation", userService.Authorisation)
	http.HandleFunc("/add-to-cart", userService.AddToCart)
	http.HandleFunc("/sale", userService.Sale)
	http.HandleFunc("/bought", userService.Bought)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
