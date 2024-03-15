package main

import (
    "fmt"
    "context"
    "encoding/json"
    "log"
    "net/http"
    "time"

    "github.com/gorilla/mux"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

type Book struct {
    ID     string `json:"id,omitempty" bson:"_id,omitempty"`
    Title  string `json:"title,omitempty" bson:"title,omitempty"`
    Author string `json:"author,omitempty" bson:"author,omitempty"`
    Year   string `json:"year,omitempty" bson:"year,omitempty"`
}

const DataBase = "mySeriesDatabase"
const Collection = "mySeries"

var client *mongo.Client

func DeleteAllBooks(response http.ResponseWriter, request *http.Request) {
    fmt.Println("Deleting all the books")
    response.Header().Set("content-type", "application/json")
    
    collection := client.Database(DataBase).Collection(Collection)
    ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
    
	_, err := collection.DeleteMany(ctx, bson.D{})
    if err != nil {
        response.WriteHeader(http.StatusInternalServerError)
        response.Write([]byte(`{"message": "` + err.Error() + `"}`))
        return
    }
    response.WriteHeader(http.StatusNoContent)
}

func DeleteBook(response http.ResponseWriter, request *http.Request) {
    response.Header().Set("content-type", "application/json")
    params := mux.Vars(request)
    id, _ := params["id"]
    collection := client.Database(DataBase).Collection(Collection)
    ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
    _, err := collection.DeleteOne(ctx, bson.M{"_id": id})
    if err != nil {
        response.WriteHeader(http.StatusInternalServerError)
        response.Write([]byte(`{"message": "` + err.Error() + `"}`))
        return
    }
    response.WriteHeader(http.StatusNoContent)
}

// CreateBook endpoint creates a book in the database.
func CreateBook(response http.ResponseWriter, request *http.Request) {
    log.Println("Creating a book")

    response.Header().Set("content-type", "application/json")
    var book Book
    _ = json.NewDecoder(request.Body).Decode(&book)
    collection := client.Database(DataBase).Collection(Collection)
    ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
    result, _ := collection.InsertOne(ctx, book)
    json.NewEncoder(response).Encode(result)
}

// GetBooks endpoint retrieves all books from the database.
func GetBooks(response http.ResponseWriter, request *http.Request) {
    fmt.Println("Geting the books") 

    response.Header().Set("content-type", "application/json")
    var books []Book
    collection := client.Database(DataBase).Collection(Collection)
    ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
    cursor, err := collection.Find(ctx, bson.M{})
    if err != nil {
        response.WriteHeader(http.StatusInternalServerError)
        response.Write([]byte(`{"message": "` + err.Error() + `"}`))
        return
    }
    defer cursor.Close(ctx)
    for cursor.Next(ctx) {
        var book Book
        cursor.Decode(&book)
        books = append(books, book)
    }
    json.NewEncoder(response).Encode(books)
}

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

    ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
    clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
    client, _ = mongo.Connect(ctx, clientOptions)

    router := mux.NewRouter()

    // Wrap your handlers with the CORS middleware
    apiRouter := router.PathPrefix("/").Subrouter()
    apiRouter.Use(corsMiddleware)
    apiRouter.HandleFunc("/books", CreateBook).Methods("POST", "OPTIONS") // Include OPTIONS to handle preflight requests
    apiRouter.HandleFunc("/books", GetBooks).Methods("GET", "OPTIONS")    // Include OPTIONS to handle preflight requests
    apiRouter.HandleFunc("/books/{id}", DeleteBook).Methods("DELETE", "OPTIONS") // Include OPTIONS to handle preflight requests
    apiRouter.HandleFunc("/books", DeleteAllBooks).Methods("DELETE", "OPTIONS") // Include OPTIONS to handle preflight requests
    
    http.ListenAndServe(":5000", router)}

