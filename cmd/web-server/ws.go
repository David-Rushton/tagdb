package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		fmt.Println(r.Method)
		fmt.Println(r.URL.Path)
		fmt.Println(r.URL.Query())
		fmt.Println(r.URL.Fragment)

		if tryServeStaticFile(w, r) {
			return
		}

		if tryServeApiRequest(w, r) {
			return
		}

		http.Error(w, "404", http.StatusNotFound)
	})

	fmt.Println("Starting server on http://localhost:31979")
	http.ListenAndServe(":31979", nil)
}

func tryServeStaticFile(w http.ResponseWriter, r *http.Request) bool {
	filePath := path.Join("wwwRoot", r.URL.Path+".html")
	if _, err := http.Dir("wwwRoot").Open(r.URL.Path + ".html"); err != nil {
		// TODO: 500.
		return false
	}

	w.Header().Set("Content-Type", "text/html")
	http.ServeFile(w, r, filePath)
	return true
}

func tryServeApiRequest(w http.ResponseWriter, r *http.Request) bool {
	if r.URL.Path == "/api/issues" {
		if r.Method == "GET" {

			resp := struct {
				Name string `json:"name"`
				Age  string `json:"age"`
			}{
				Name: "David",
				Age:  "30",
			}

			w.Header().Set("Content-Type", "application/json")
			err := json.NewEncoder(w).Encode(resp)
			if err != nil {
				http.Error(w, "500", http.StatusInternalServerError)
			}

			return true
		}
	}

	return false
}
