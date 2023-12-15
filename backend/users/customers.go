package users

import (
	//"backend/bd"
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	_ "github.com/lib/pq"
)

func bd(str string, str1 string) error {
	connStr := "user=postgres password=123 dbname=wzpwshep sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	if str == "Exec" {
		result, err := db.Exec(str1)
		fmt.Println(result.RowsAffected())
		return err
	}
	if str == "QueryRow" {
		var resultbd string
		err := db.QueryRow(str1).Scan(&resultbd)
		return err
	}
	return err
}

func Registration(w http.ResponseWriter, r *http.Request) {

	name := r.FormValue("name")
	login := r.FormValue("login")
	password := r.FormValue("password")

	str1 := fmt.Sprintf("insert into users (name,login,password) values (%s,%s,%s)", name, login, password)
	err := bd("Exec", str1)

	if err != nil {
		fmt.Fprintf(w, "%s", "Ошибка при выполнении запроса")
	}
	fmt.Fprintf(w, "%s", "Success")
}

func Authorisation(w http.ResponseWriter, r *http.Request) {

	login := r.FormValue("login")
	password := r.FormValue("password")

	str1 := fmt.Sprintf("SELECT * FROM users WHERE login=%s AND password=%s", login, password)
	err := bd("Exec", str1)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Fprintf(w, "%s", "Логин и пароль неверны")
		} else {
			fmt.Fprintf(w, "%s", "Ошибка при выполнении запроса")
		}
		return
	}
	fmt.Fprintf(w, "%s", "Логин и пароль верны")
}

func Add_to_cart(w http.ResponseWriter, r *http.Request) {

	name := r.FormValue("name")
	goods := r.FormValue("goods")

	str1 := fmt.Sprintf("UPDATE users SET goods = array_append(goods, %s) WHERE name = %s", goods, name) //нужно будет сделать либо имя уникальным либо искать по логину
	err := bd("Exec", str1)
	if err != nil {
		fmt.Fprintf(w, "%s", "Ошибка при выполнении запроса")
	}

	fmt.Fprintf(w, "%s", "Success")
}

func Sale(w http.ResponseWriter, r *http.Request) {

	name := r.FormValue("name")

	var result string
	str1 := fmt.Sprintf("SELECT total_sum FROM users WHERE name=%s ", name)
	err := bd("QueryRow", str1)
	if err != nil {
		fmt.Fprintf(w, "%s", "Ошибка при выполнении запроса")
		return
	}

	fmt.Fprintf(w, "%s", result)
}

func Bought(w http.ResponseWriter, r *http.Request) {

	name := r.FormValue("name")
	sum, _ := strconv.Atoi(r.FormValue("sum"))

	str1 := fmt.Sprintf("UPDATE users SET total_sum = total_sum  + %d WHERE name = %s", sum, name) //нужно будет сделать либо имя уникальным либо искать по логину
	err := bd("Exec", str1)
	if err != nil {
		fmt.Fprintf(w, "%s", "Ошибка при выполнении запроса")
		return
	}

	fmt.Fprintf(w, "%s", "Success")
}
