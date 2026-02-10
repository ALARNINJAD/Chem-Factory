package main

import "net/http"

func login(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
