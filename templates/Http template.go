package main

import (
	"fmt"
	"net/http"
	"text/template"
)

type User struct {
	Id     int
	Name   string
	Active bool
}

func (u *User) PrintActive() string {
	if !u.Active {
		return ""
	}
	return "is active"
}

func main() {
	tmpl, err := template.New("").ParseFiles("func.html")
	if err != nil {
		panic(err)
	}

	users := []User{
		{1, "Niko", true},
		{2, "Dasha", true},
		{3, "Julia", true},
		{4, "Klim", false},
	}

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err := tmpl.ExecuteTemplate(w, "func.html",
			struct {
				Users []User
			}{
				users,
			})
		if err != nil {
			panic(err)
		}
	})
	fmt.Println("starting server at :8080")
	http.ListenAndServe(":8080", nil)
}
