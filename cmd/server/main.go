package main

import "net/http"

func main() {
	http.HandleFunc("/login", login)
	http.HandleFunc("/register", register)
}
