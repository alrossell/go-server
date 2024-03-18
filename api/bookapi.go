package api

import (
    "fmt"
    "context"
    "encoding/json"
    "log"
    "net/http"
    "time"

    "example.com/my-go-project/global"
    
    "github.com/gorilla/mux"
    "go.mongodb.org/mongo-driver/bson"
)

const DataBase = global.DataBase
const Collection = global.Collection 

type Book struct {
    ID     string `json:"id,omitempty" bson:"_id,omitempty"`
    Title  string `json:"title,omitempty" bson:"title,omitempty"`
    Author string `json:"author,omitempty" bson:"author,omitempty"`
    Year   string `json:"year,omitempty" bson:"year,omitempty"`
}

func DeleteAllBooks(response http.ResponseWriter, request *http.Request) {
    fmt.Println("Deleting all the books")

    client := global.GetClient() 
    fmt.Println("Deleting all the books")
    response.Header().Set("content-type", "application/json")
    
    collection := client.Database(DataBase).Collection(Collection)
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
	_, err := collection.DeleteMany(ctx, bson.D{})
    if err != nil {
        response.WriteHeader(http.StatusInternalServerError)
        response.Write([]byte(`{"message": "` + err.Error() + `"}`))
        return
    }
    response.WriteHeader(http.StatusNoContent)
}

func DeleteBook(response http.ResponseWriter, request *http.Request) {
    fmt.Println("Deleting a book")

    client := global.GetClient() 
    response.Header().Set("content-type", "application/json")
    params := mux.Vars(request)
    id, _ := params["id"]

    collection := client.Database(DataBase).Collection(Collection)
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

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
    client := global.GetClient()

    response.Header().Set("content-type", "application/json")
    var book Book
    _ = json.NewDecoder(request.Body).Decode(&book)
    collection := client.Database(DataBase).Collection(Collection)
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    result, _ := collection.InsertOne(ctx, book)
    json.NewEncoder(response).Encode(result)
}

// GetBooks endpoint retrieves all books from the database.
func GetBooks(response http.ResponseWriter, request *http.Request) {
    fmt.Println("Geting the books") 
    client := global.GetClient()

    response.Header().Set("content-type", "application/json")
    var books []Book
    collection := client.Database(DataBase).Collection(Collection)
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

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
