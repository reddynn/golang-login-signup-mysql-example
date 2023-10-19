package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

type User struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

func Connect() *sql.DB {
	connect, err := sql.Open("mysql", "root:Passwd123@tcp(127.0.0.1:3306)/mydb")
	if err != nil {
		panic(err)
	}
	return connect
}
func main() {
	http.HandleFunc("/", Home)
	http.HandleFunc("/signup", Signup)
	http.HandleFunc("/login", Login)
	log.Fatal(http.ListenAndServe(":8080", nil))

}
func Home(w http.ResponseWriter, r *http.Request) {

	http.ServeFile(w, r, "index.html")

}

func Signup(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.ServeFile(w, r, "signup.html")
	}
	db := Connect()
	defer db.Close()
	username := r.FormValue("username")
	password := r.FormValue("password")
	var user string
	err := db.QueryRow("select username from users where username=?", username).Scan(&user)

	switch {
	case err == sql.ErrNoRows:
		hashedpassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {

			http.Error(w, "unable to create password", 500)
			return
		}
		_, err = db.Query("insert into users(username,password) values(?,?)", username, hashedpassword)
		if err != nil {
			http.Error(w, "unable to create user", 500)
			return
		}

		w.Write([]byte("user created"))
		return
	case err != nil:
		http.Error(w, "unable to create user", 500)
		return
	default:
		http.Redirect(w, r, "/", http.StatusMovedPermanently)

	}

}

func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		http.ServeFile(w, r, "login.html")
	}
	db := Connect()
	defer db.Close()
	username := r.FormValue("username")
	password := r.FormValue("password")
	var user string
	var passwd string
	err := db.QueryRow("select username,password from users where username=?", username).Scan(&user, &passwd)

	if err != nil {
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(passwd), []byte(password))
	if err != nil {
		http.Redirect(w, r, "/", http.StatusMovedPermanently)
		return
	}
	w.Write([]byte("logged in successfully"))
}
