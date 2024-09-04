package main

import (
	"fmt"
	"net/http"
	"time"
)

func mainPage(w http.ResponseWriter, r *http.Request) {
	session, err := r.Cookie("session_id")
	loggedIn := (err != http.ErrNoCookie)

	if loggedIn {
		fmt.Fprintf(w, `<a href="/logout">logout</a>`)
		fmt.Fprintf(w, "Welcome,"+session.Value)
	} else {
		fmt.Fprintf(w, `<a href="/login">login</a>`)
		fmt.Fprintf(w, "Hi, you need to login!")
	}
}

func loginPage(w http.ResponseWriter, r *http.Request) {
	expiration := time.Now().Add(10 * time.Hour)

	cookie := http.Cookie{
		Name:    "session_id",
		Value:   "niko",
		Expires: expiration,
	}
	http.SetCookie(w, &cookie)
	http.Redirect(w, r, "/", http.StatusFound)
}

func logoutPage(w http.ResponseWriter, r *http.Request) {
	session, err := r.Cookie("session_id")

	if err == http.ErrNoCookie {
		http.Redirect(w, r, "/", http.StatusFound)
		return
	}

	session.Expires = time.Now().AddDate(0, 0, -1)
	http.SetCookie(w, session)
	http.Redirect(w, r, "/", http.StatusFound)
}

func main() {
	http.HandleFunc("/", mainPage)
	http.HandleFunc("/login", loginPage)
	http.HandleFunc("/logout", logoutPage)

	addr := ":8080"
	fmt.Println("Starting server at", addr)
	http.ListenAndServe(addr, nil)
}