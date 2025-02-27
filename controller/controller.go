package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/taneja-garvit/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDB connection configuration
const connectionString = "your-mongodb-connection-string-here"
const dbName = "netflix"
const colName = "watchlist"

// Collection variable to hold the MongoDB collection instance
var collection *mongo.Collection

// init function connects to MongoDB and initializes the collection instance
func init() {
	// Set client options
	clientOptions := options.Client().ApplyURI(connectionString)

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Check if the connection was successful
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal("Could not connect to MongoDB:", err)
	}

	fmt.Println("MongoDB connection successful")

	// Get the collection instance
	collection = client.Database(dbName).Collection(colName)
	fmt.Println("Collection instance is ready")
}

//mongoDB helpers -file

// insert 1 record

func insertOneMovie(movie model.Netflix) { // we took model.netflix bec here we takes the name of the packeage instead of file
	inserted, err := collection.InsertOne(context.Background(), movie)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Inserted One Movie with ID:", inserted.InsertedID)
}

func updateOneRecord(movieId string) {
	id, _ := primitive.ObjectIDFromHex(movieId)

	filter := bson.M{"_id": id} // bson.M is used where we have to use key value paires thing or like map thing
	update := bson.M{"$set": bson.M{"watched": true}}

	result, err := collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Modifeid Count", result.ModifiedCount)
}

func deleteOneRecord(movieId string) {
	id, _ := primitive.ObjectIDFromHex(movieId)
	filter := bson.M{"_id": id}
	deletedCount, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Deleted Count", deletedCount)
}

func deleteMany() int64 {
	// filter:= bson.D{{}}  it means selecting all values
	deleteResult, err := collection.DeleteMany(context.Background(), bson.D{{}}, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Deleted Count", deleteResult.DeletedCount)
	return deleteResult.DeletedCount
}

// get all movies from database

// GetAllMovies retrieves all movie records from MongoDB
func getAllMovies() []primitive.M {
	// Find all documents in the collection
	cur, err := collection.Find(context.Background(), bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}
	defer cur.Close(context.Background()) // Close cursor after use

	// Create a slice to hold the result
	var movies []primitive.M

	// Iterate through the cursor
	for cur.Next(context.Background()) {
		var movie bson.M
		err := cur.Decode(&movie)
		if err != nil {
			log.Fatal(err)
		}
		movies = append(movies, movie)
	}

	// Check for any errors after iteration
	if err := cur.Err(); err != nil {
		log.Fatal(err)
	}

	// Print the result or use it
	fmt.Println("Retrieved Movies:", movies)

	// Return the list of movies
	return movies
}

// Actual export controllers

func GetAllMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencode")
	allMovies := getAllMovies()
	json.NewEncoder(w).Encode(allMovies)
}

func CreateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencode")
	w.Header().Set("Alow-Control-Allow-Methods", "POST")

	var movie model.Netflix
	h, _ := josn.NewDecoder(r.Body).Decode(&movie)
	insertOneMovie(movie)
	json.NewEncoder(w).Encode(movie)

}

func MarkAsWatched(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/x-www-form-urlencode")
	w.Header().Set("Alow-Control-Allow-Methods", "PUT")

	params := mux.Vars(r)
	updateOneRecord(params["id"])
	json.NewEncoder(w).Encode(params["id"])
}

func DeleteOneMovie(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/x-www-form-urlencode")
	w.Header().Set("Alow-Control-Allow-Methods", "DELETE")

	params := mux.Vars(r)
	deleteOneRecord(params["id"])
	json.NewEncoder(w).Encode(params["id"])
}

func DeleteAllMovies(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/x-www-form-urlencode")
	w.Header().Set("Alow-Control-Allow-Methods", "DELETE")

	count := deleteMany()
	json.NewEncoder(w).Encode(count)
}