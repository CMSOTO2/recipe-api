package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Collection struct {
	ID      string   `json:"_id"`
	Recipes []Recipe `json:"recipes"`
}
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

var mongoClient *mongo.Client

const mongoURI = "mongodb+srv://carlosmsoto2:carlos123@recipes.efzgczq.mongodb.net/?retryWrites=true&w=majority"

func init() {
	if err := connectToMongoDB(); err != nil {
		log.Fatal("Could not connect to MongoDB")
	}
}

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("loading environment variables:", err)
	}
	r := gin.Default()
	r.GET("/recipes", getAllRecipes)
	r.GET("/recipe/:id", getSingleRecipe)

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	http.ListenAndServe(":"+port, r)
}

func connectToMongoDB() error {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(mongoURI).SetServerAPIOptions(serverAPI)
	client, err := mongo.Connect(context.TODO(), opts)
	if err == nil {
		mongoClient = client
		return nil
	}

	return err
}

func getAllRecipes(c *gin.Context) {
	cursor, err := mongoClient.Database("recipe-web-app").Collection("recipes").Find(context.TODO(), bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}
	var recipes []bson.M
	if err = cursor.All(context.TODO(), &recipes); err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, recipes)
}

func getSingleRecipe(c *gin.Context) {
	id := c.Param("id")

	objectId, err := primitive.ObjectIDFromHex(id)

	if err != nil {
		c.AbortWithStatus(400)
		return
	}
	var recipe bson.M

	err = mongoClient.Database("recipe-web-app").Collection("recipes").FindOne(context.TODO(), bson.D{{"_id", objectId}}).Decode(&recipe)
	if err != nil {
		c.AbortWithStatus(404)
		return
	}

	c.JSON(http.StatusOK, recipe)
}
