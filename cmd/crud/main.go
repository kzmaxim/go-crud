package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"html/template"
	"net/http"
)

var database *sql.DB

type Game struct {
	Id    int
	Title string
	Genre string
	Price int
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	rows, err := database.Query("SELECT * FROM gamesdb")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	games := []Game{}
	for rows.Next() {
		g := Game{}
		err := rows.Scan(&g.Id, &g.Title, &g.Genre, &g.Price)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		games = append(games, g)
	}
	tmpl, _ := template.ParseFiles("../../templates/index.html")
	tmpl.Execute(w, games)
}

func CreateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()
		title := r.FormValue("title")
		genre := r.FormValue("genre")
		price := r.FormValue("price")

		_, err := database.Exec("INSERT INTO gamesdb (Title, Genre, Price) VALUES (?, ?, ?)", title, genre, price)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		http.Redirect(w, r, "/", 301)
	} else {
		http.ServeFile(w, r, "../../templates/create.html")
	}
}

func EditPage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	row := database.QueryRow("SELECT * FROM gamesdb WHERE Id = ?", id)
	g := Game{}
	row.Scan(&g.Id, &g.Title, &g.Genre, &g.Price)
	tmpl, _ := template.ParseFiles("../../templates/edit.html")
	tmpl.Execute(w, g)
}

func EditHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	r.ParseForm()
	title := r.FormValue("title")
	genre := r.FormValue("genre")
	price := r.FormValue("price")
	_, err := database.Exec("UPDATE gamesdb SET Title = ?, Genre = ?, Price = ? WHERE Id = ?", title, genre, price, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", 301)
}

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	_, err := database.Exec("DELETE FROM gamesdb WHERE Id = ?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", 301)
}

func main() {
	db, err := sql.Open("mysql", "root:3maxim14@tcp(127.0.0.1:3306)/gocrud")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer db.Close()
	database = db
	router := mux.NewRouter()

	router.HandleFunc("/", IndexHandler)
	router.HandleFunc("/create", CreateHandler).Methods("GET", "POST")
	router.HandleFunc("/edit/{id:[0-9]+}", EditPage).Methods("GET")
	router.HandleFunc("/edit/{id:[0-9]+}", EditHandler).Methods("POST")
	router.HandleFunc("/delete/{id:[0-9]+}", DeleteHandler)

	http.Handle("/", router)
	fmt.Println("Server listening on port 8080")
	http.ListenAndServe("localhost:8080", router)
}
