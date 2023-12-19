package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Recipe struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Tags        []string `json:"tags"`
	Rating      float64  `json:"rating"`
	TimeToCook  int      `json:"timeToCook"`
	IsFavorited bool     `json:"isFavorited"`
	ImageURL    string   `json:"imageUrl"`
	Steps       []string `json:"steps"`
	Ingredients []string `json:"ingredients"`
}

var client *mongo.Client
var recipesCollection *mongo.Collection

func init() {
	mongoURI := "mongodb+srv://carlosmsoto2:carlos123@recipes.efzgczq.mongodb.net/"
	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	err = client.Ping(context.Background(), nil)
	if err != nil {
		log.Fatal("Could not connect to MongoDB:", err)
	}

	recipesCollection = client.Database("recipe-web-app").Collection("recipes")
}

func main() {
	router := gin.Default()

	router.GET("/recipes", getRecipesHandler)
	router.POST("/recipes", createRecipeHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	address := fmt.Sprintf(":%s", port)
	log.Printf("Server is running on %s...", address)
	log.Fatal(http.ListenAndServe(address, router))
}

func getRecipesHandler(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := recipesCollection.Find(ctx, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch recipes"})
		return
	}

	defer cursor.Close(ctx)

	var recipes []Recipe
	if err := cursor.All(ctx, &recipes); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode recipes"})
		return
	}

	c.JSON(http.StatusOK, recipes)
}

func createRecipeHandler(c *gin.Context) {
	var recipe Recipe
	if err := c.BindJSON(&recipe); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	result, err := recipesCollection.InsertOne(ctx, recipe)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to insert recipe"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"id": result.InsertedID})
}
