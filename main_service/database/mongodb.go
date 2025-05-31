package database

import (
	"context"
	"fmt"
	"log"
	"restaurant_reviews/internal"
	"time"

	"go.mongodb.org/mongo-driver/bson"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoDB *mongo.Client
var Cxt context.Context

func ConnectMongo(ctx context.Context) error {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %s", err.Error())
	}

	MongoDB = client
	Cxt = ctx
	return nil
}

func GetUserID(email string) (string, error) {
	collection := MongoDB.Database("restaurantdb_1").Collection("users")
	filter := bson.M{
		"email": email,
	}

	var user internal.User

	err := collection.FindOne(Cxt, filter).Decode(&user)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return "", nil
		}
		return "", err
	}

	return user.ID, nil
}

func GetUser(email string, password string) (internal.User, error) {
	collection := MongoDB.Database("restaurantdb_1").Collection("users")
	filter := bson.M{
		"email":        email,
		"passwordHash": password,
	}

	var user internal.User

	err := collection.FindOne(Cxt, filter).Decode(&user)
	if err != nil {
		return user, err
	}

	return user, nil
}

func RegisterUser(email string, password string, role string) (interface{}, error) {
	collection := MongoDB.Database("restaurantdb_1").Collection("users")
	now := time.Now().UTC()

	user := bson.D{
		{Key: "email", Value: email},
		{Key: "role", Value: role},
		{Key: "passwordHash", Value: password},
		{Key: "registerAt", Value: now},
	}

	insertResult, err := collection.InsertOne(Cxt, user)
	if err != nil {
		log.Fatal(err)
		return "Not Created", fmt.Errorf("failed to create user")
	}

	return insertResult.InsertedID, nil
}

func CreateFeedBack(id string, userId string, restaurantId string, text string, rating float64) (internal.Review, error) {
	collection := MongoDB.Database("restaurantdb_1").Collection("reviews")
	now := time.Now().UTC()

	review := bson.D{
		{Key: "_id", Value: id},
		{Key: "userId", Value: userId},
		{Key: "restaurantId", Value: restaurantId},
		{Key: "text", Value: text},
		{Key: "rating", Value: rating},
		{Key: "createdAt", Value: now},
	}

	result := internal.Review{
		ID:           id,
		UserID:       userId,
		RestaurantID: restaurantId,
		Text:         text,
		Rating:       rating,
		CreatedAt:    now,
	}

	_, err := collection.InsertOne(Cxt, review)
	if err != nil {
		log.Fatal(err)
		return result, fmt.Errorf("failed to create review: %s", err)
	}

	// Update restaurant rating
	ratingsCollection := MongoDB.Database("restaurantdb_1").Collection("ratings")

	// Get current rating or create new if not exists
	var currentRating internal.Rating
	err = ratingsCollection.FindOne(Cxt, bson.D{{Key: "restaurantId", Value: restaurantId}}).Decode(&currentRating)
	if err != nil && err != mongo.ErrNoDocuments {
		return result, fmt.Errorf("failed to get current rating: %s", err)
	}

	if err == mongo.ErrNoDocuments {
		// First review for this restaurant
		newRating := internal.Rating{
			RestaurantID:  restaurantId,
			AverageRating: rating,
			ReviewCount:   1,
		}
		_, err = ratingsCollection.InsertOne(Cxt, newRating)
		if err != nil {
			return result, fmt.Errorf("failed to create initial rating: %s", err)
		}
	} else {
		// Update existing rating
		newCount := currentRating.ReviewCount + 1
		newAverage := (currentRating.AverageRating*float64(currentRating.ReviewCount) + rating) / float64(newCount)

		_, err = ratingsCollection.UpdateOne(
			Cxt,
			bson.D{{Key: "restaurantId", Value: restaurantId}},
			bson.D{{Key: "$set", Value: bson.D{
				{Key: "averageRating", Value: newAverage},
				{Key: "reviewCount", Value: newCount},
			}}},
		)
		if err != nil {
			return result, fmt.Errorf("failed to update rating: %s", err)
		}
	}

	return result, nil
}

func DeleteUser(id int) (int, error) {
	collection := MongoDB.Database("restaurantdb_1").Collection("users")

	review := bson.D{{Key: "user_id", Value: id}}
	_, err := collection.DeleteOne(Cxt, review)

	if err != nil {
		return id, fmt.Errorf("failed to delete user: %s", err)
	}

	return id, nil
}
