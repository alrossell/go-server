package main

import (
    "fmt"
    "context"
    "encoding/json"
    // "log"
    "net/http"
    "time"

    "github.com/gorilla/mux"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
)

// Book struct represents a book document in MongoDB.
type Book struct {
    ID     string `json:"id,omitempty" bson:"_id,omitempty"`
    Title  string `json:"title,omitempty" bson:"title,omitempty"`
    Author string `json:"author,omitempty" bson:"author,omitempty"`
    Year   string `json:"year,omitempty" bson:"year,omitempty"`
}

const DataBase = "mySeriesDatabase"
const Collection = "mySeries"

var client *mongo.Client

// CreateBook endpoint creates a book in the database.
func CreateBook(response http.ResponseWriter, request *http.Request) {
    fmt.Println("Creating a Book")

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

func main() {

   fmt.Println("Starting Server") 

    ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
    clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
    client, _ = mongo.Connect(ctx, clientOptions)

    router := mux.NewRouter()
    router.HandleFunc("/books", CreateBook).Methods("POST")
    router.HandleFunc("/books", GetBooks).Methods("GET")
    http.ListenAndServe(":5000", router)
}

