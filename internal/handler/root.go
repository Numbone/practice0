package handler

import (
	"net/http"
)

func RootHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.ServeFile(w, r, "web/index.html")
			return
		}
		http.FileServer(http.Dir("web")).ServeHTTP(w, r)
	}
}
