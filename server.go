package main

import (
    "fmt"
    "net/http"

    "example.com/my-go-project/global"
    "example.com/my-go-project/api"

    "github.com/gorilla/mux"
)

type Book struct {
    ID     string `json:"id,omitempty" bson:"_id,omitempty"`
    Title  string `json:"title,omitempty" bson:"title,omitempty"`
    Author string `json:"author,omitempty" bson:"author,omitempty"`
    Year   string `json:"year,omitempty" bson:"year,omitempty"`
}

const DataBase = "mySeriesDatabase"
const Collection = "mySeries"

// corsMiddleware sets up the CORS headers for response, allowing cross-origin requests.
func corsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Access-Control-Allow-Origin", "http://localhost:8081") // Allow all origins
        w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE") // Allowed methods
        w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization") // Allowed headers

         w.Header().Set("Access-Control-Allow-Credentials", "true")
         // If it's a preflight request, respond with 200 without calling next handler.
        if r.Method == "OPTIONS" {
            w.WriteHeader(http.StatusOK)
            return
        }

        next.ServeHTTP(w, r)
    })
}

func main() {
    fmt.Println("Starting Server") 

    global.CreateClient()

    router := mux.NewRouter()

    // Wrap your handlers with the CORS middleware
    apiRouter := router.PathPrefix("/").Subrouter()
    apiRouter.Use(corsMiddleware)
    apiRouter.HandleFunc("/books", api.CreateBook).Methods("POST", "OPTIONS") // Include OPTIONS to handle preflight requests
    apiRouter.HandleFunc("/books", api.GetBooks).Methods("GET", "OPTIONS")    // Include OPTIONS to handle preflight requests
    apiRouter.HandleFunc("/books/{id}", api.DeleteBook).Methods("DELETE", "OPTIONS") // Include OPTIONS to handle preflight requests
    apiRouter.HandleFunc("/books", api.DeleteAllBooks).Methods("DELETE", "OPTIONS") // Include OPTIONS to handle preflight requests
    
    http.ListenAndServe(":5000", router)}

