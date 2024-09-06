package main

import (
	"fmt"
	"net/http"
	"time"
)

func mainPage(w http.ResponseWriter, r *http.Request) {
	session, err := r.Cookie("session_id")

	isLogin := err != http.ErrNoCookie
	if isLogin {
		fmt.Fprintf(w, `<a href="/logout">logout</a><br>`)
		fmt.Fprintf(w, "Welcome, %s", session.Value)
	} else {
		fmt.Fprintf(w, `<a href="/login">login</a><br>`)
		fmt.Fprintf(w, `You need to login`)
	}
}

func loginPage(w http.ResponseWriter, r *http.Request) {
	expiration := time.Now().Add(10 * time.Hour)
	cookie := &http.Cookie{
		Name:    "session_id",
		Value:   "niko",
		Expires: expiration,
	}
	http.SetCookie(w, cookie)
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

func adminIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, `<a href="/">site index</a>`)
	fmt.Fprintln(w, "Admin main page")
}

func panicPage(w http.ResponseWriter, r *http.Request) {
	panic("this must be recovered")
}

func main() {
	adminMux := http.NewServeMux()
	adminMux.HandleFunc("/admin/", adminIndex)
	adminMux.HandleFunc("/admin/panic", panicPage)

	// set middleware
	adminHandler := adminAuthMiddleware(adminMux)

	siteMux := http.NewServeMux()
	siteMux.Handle("/admin/", adminHandler)
	siteMux.HandleFunc("/login", loginPage)
	siteMux.HandleFunc("/logout", logoutPage)
	siteMux.HandleFunc("/", mainPage)

	// set middleware
	siteHandler := accessLogMiddleware(siteMux)
	siteHandler = panicMiddleware(siteHandler)

	addr := ":8080"
	fmt.Println("Starting server at", addr)
	http.ListenAndServe(addr, siteHandler)
}
