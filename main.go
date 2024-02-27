package main

import (
	"context"
	"net/http"
	"log"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/bson"
)

var collection *mongo.Collection

type book struct {
	ID       string `json:"id" bson:"_id"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	Quantity int    `json:"quantity"`
	MaxPresent int    `json:"max_present"`
}


func initMongoDB() *mongo.Client {
	// Set up MongoDB client options
	clientOptions := options.Client().ApplyURI("mongodb+srv://adhayan436:hd8zLAhYllDgsGPu@cluster0.aq7ozho.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0")

	// Connect to MongoDB
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Check the connection
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}

	// Set the collection
	collection = client.Database("database1").Collection("books")

	return client
}

func getBookByID(id string) (*book, error) {
	var result book
	err := collection.FindOne(context.Background(), bson.M{"_id": id}).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func bookByID(c *gin.Context) {
	id := c.Param("id")
	book, err := getBookByID(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found"})
		return
	}
	c.IndentedJSON(http.StatusOK, book)
}


func getBooks(c *gin.Context) {
	cursor, err := collection.Find(context.Background(), bson.M{})
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch books"})
		return
	}
	defer cursor.Close(context.Background())

	var books []book
	err = cursor.All(context.Background(), &books)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode books"})
		return
	}

	c.IndentedJSON(http.StatusOK, books)
}

func checkoutBook(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
		return
	}

	book, err := getBookByID(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found"})
		return
	}

	if book.Quantity <= 0 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Books Not Available"})
		return
	}

	// Decrement the quantity by 1 in MongoDB
	updateResult, err := collection.UpdateOne(
		context.Background(),
		bson.M{"_id": id},
		bson.D{{"$inc", bson.D{{"quantity", -1}}}},
	)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Failed to update quantity"})
		return
	}

	if updateResult.ModifiedCount == 0 {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Failed to update quantity"})
		return
	}

	// Update local book object
	book.Quantity--

	c.IndentedJSON(http.StatusOK, book)
}


func returnBook(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Bad request"})
		return
	}

	book, err := getBookByID(id)
	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found"})
		return
	}

	// Increment the quantity
	book.Quantity++

	// Check if quantity exceeds MaxPresent
	if book.Quantity > book.MaxPresent {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Maximum present limit exceeded"})
		return
	}

	// Update the book in MongoDB
	updateResult, err := collection.UpdateOne(
		context.Background(),
		bson.M{"_id": id},
		bson.D{{"$set", bson.D{{"quantity", book.Quantity}}}},
	)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Failed to update quantity"})
		return
	}

	if updateResult.ModifiedCount == 0 {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "Failed to update quantity"})
		return
	}

	c.IndentedJSON(http.StatusOK, book)
}


func createBook(c *gin.Context) {
	var newBook book
	if err := c.BindJSON(&newBook); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON payload"})
		return
	}

	_, err := collection.InsertOne(context.Background(), newBook)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to create book"})
		return
	}

	c.IndentedJSON(http.StatusCreated, newBook)
}

func createMultiBook(c *gin.Context) {
	var newBooks []book
	if err := c.BindJSON(&newBooks); err != nil {
		// If the JSON payload cannot be bound to []book, try binding to a single book
		var singleBook book
		if err := c.BindJSON(&singleBook); err != nil {
			c.IndentedJSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON payload"})
			return
		}
		// If single book successfully bound, add it to the newBooks slice
		newBooks = append(newBooks, singleBook)
	}

	var insertResult *mongo.InsertManyResult
	bookDocuments := make([]interface{}, len(newBooks))
	for i, bk := range newBooks {
		bookDocuments[i] = bk
	}

	insertResult, err := collection.InsertMany(context.Background(), bookDocuments)
	if err != nil {
		c.IndentedJSON(http.StatusInternalServerError, gin.H{"error": "Failed to create books"})
		return
	}

	c.IndentedJSON(http.StatusCreated, gin.H{"message": fmt.Sprintf("Created %d book(s)", len(insertResult.InsertedIDs))})
}


func main(){
	mongoClient := initMongoDB()
	defer mongoClient.Disconnect(context.Background())

	router := gin.Default()
	fmt.Println("API has been started!")

	router.GET("/books", getBooks)
	router.GET("/books/:id", bookByID)
	router.POST("/book", createBook)
	router.POST("/books", createMultiBook)
	router.PATCH("/checkout/:id", checkoutBook)
	router.PATCH("/return/:id", returnBook)
	router.Run("localhost:8080")
}