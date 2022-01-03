package main

import (
	"log"
	"net/http"
	"net/url"
	"sync"

	"github.com/gorilla/mux"
)

var wg sync.WaitGroup

func handleRequests(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	path := ""
	if r.Form.Has("path") {
		path, _ = url.QueryUnescape(r.Form["path"][0])
		switch r.Method {
		case "GET":
			secret, err := client.Logical().Read("secret/data/" + path)
			if err != nil {
				log.Println("ERROR: Failed to read secret: ", err)
			}

			_, ok := secret.Data["data"].(map[string]interface{})
			if !ok {
				log.Println("ERROR: Type assertion failed! No \"data\" field")
			} else {
				w.Write([]byte("Success"))
				log.Println(secret)
				log.Println("GET request successful!")

			}
		case "POST", "PUT":
			tempMap := map[string][]string{}
			for k, v := range r.PostForm {
				tempMap[k] = v
			}
			secret, err := client.Logical().Write("secret/data/"+path, map[string]interface{}{"data": tempMap})
			if err != nil {
				log.Println("ERROR: Failed to write secret: ", err)
			} else {
				w.Write([]byte("Success"))
				log.Println(secret)
				log.Println("POST/PUT request successful!")
			}
		case "DELETE":
			_, err := client.Logical().Delete("secret/data/something")
			if err != nil {
				log.Println("ERROR: Failed to delete secret: ", err)
			} else {
				w.Write([]byte("Success"))
				log.Println("DELETE request successful!")
			}
		}
	} else {
		w.Write([]byte("path param required!"))
	}
}

func startServer() {
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/kv", handleRequests)
	wg.Add(1)
	go http.ListenAndServe(":8000", router)
	log.Println("Server started, listening on port 8000")
	wg.Wait()
}
