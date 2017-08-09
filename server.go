package main

import (
	_ "github.com/lib/pq"
	"encoding/json"
	"net/http"
	"path"
	"strconv"
	"database/sql"
	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"
	"os"
	"github.com/auth0-community/auth0"
	"fmt"
	jose "gopkg.in/square/go-jose.v2"
)

type Post struct {
	Id      int    `json:"id"`
	Content string `json:"content"`
	Author  string `json:"author"`
}

type Pet struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Owner   string `json:"owner"`
}


// connect to the Db
func init() {
	var err error
	Db, err = sql.Open("postgres", "user=fido dbname=petstay password=woof sslmode=disable")
	if err != nil {
		panic(err)
	}
}

func main() {

	r := mux.NewRouter()

	r.Handle("/", http.FileServer(http.Dir("./views/")))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	r.Handle("/pet/{id}", authMiddleware(GetPetHandler)).Methods("GET")
	r.Handle("/pets", authMiddleware(GetPetsHandler)).Methods("GET")
	r.Handle("/pet", authMiddleware(PostPetHandler)).Methods("POST")

	http.ListenAndServe(":3000", handlers.LoggingHandler(os.Stdout, r))

}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// JC: these are from Auth0 management console
		secret := []byte("pwsy9HvbACAKQYlw1Rp1EAL6ej2OfCZ3")
		audience := "https://jamez.com"
		issuer := "https://speak2jezza.au.auth0.com/"
		secretProvider := auth0.NewKeyProvider(secret)

		var audiences []string
		configuration := auth0.NewConfiguration(secretProvider, append(audiences, audience), issuer, jose.HS256)
		validator := auth0.NewValidator(configuration)

		token, err := validator.ValidateRequest(r)

		if err != nil {
			fmt.Println(err)
			fmt.Println("Token is not valid:", token)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized"))
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

// main handler function
func handlePetRequest(w http.ResponseWriter, r *http.Request) {
	var err error
	switch r.Method {
	case "GET":
		//err = handleGetPet(w, r)
	case "POST":
		err = handlePostPet(w, r)
	case "PUT":
		err = handlePut(w, r)
	case "DELETE":
		err = handleDelete(w, r)
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

var GetPetHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		return
	}

	pet, err := retrievePet(id)
	if err != nil {
		return
	}

	output, err := json.MarshalIndent(&pet, "", "\t\t")
	if err != nil {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(output)
	return
})

var PostPetHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	len := r.ContentLength
	body := make([]byte, len)
	r.Body.Read(body)
	var pet Pet
	json.Unmarshal(body, &pet)
	err := pet.create()
	if err != nil {
		return
	}
	w.WriteHeader(200)
	return
})

var GetPetsHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	vals := r.URL.Query()
	offsets :=  vals["offset"]

	offset := 0
	if offsets != nil {
		offset1, err := strconv.Atoi(offsets[0])
		if err != nil {
			return
		}
		offset = offset1
	}
	limits :=  vals["limit"]

	limit := 1
	if limits != nil {
		limit1, err := strconv.Atoi(limits[0])
		if err != nil {
			return
		}
		limit = limit1
	}

	pets, err := retrievePets(offset, limit)
	if err != nil {
		return
	}
	output, err := json.MarshalIndent(&pets, "", "\t\t")
	if err != nil {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write(output)
	return

})


// Create a post
// POST /post/

func handlePostPet(w http.ResponseWriter, r *http.Request) (err error) {
	len := r.ContentLength
	body := make([]byte, len)
	r.Body.Read(body)
	var pet Pet
	json.Unmarshal(body, &pet)
	err = pet.create()
	if err != nil {
		return
	}
	w.WriteHeader(200)
	return
}

// Update a post
// PUT /post/1
func handlePut(w http.ResponseWriter, r *http.Request) (err error) {
	id, err := strconv.Atoi(path.Base(r.URL.Path))
	if err != nil {
		return
	}
	post, err := retrieve(id)
	if err != nil {
		return
	}
	len := r.ContentLength
	body := make([]byte, len)
	r.Body.Read(body)
	json.Unmarshal(body, &post)
	err = post.update()
	if err != nil {
		return
	}
	w.WriteHeader(200)
	return
}

// Delete a post
// DELETE /post/1
func handleDelete(w http.ResponseWriter, r *http.Request) (err error) {
	id, err := strconv.Atoi(path.Base(r.URL.Path))
	if err != nil {
		return
	}
	post, err := retrieve(id)
	if err != nil {
		return
	}
	err = post.delete()
	if err != nil {
		return
	}
	w.WriteHeader(200)
	return
}

