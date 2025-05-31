package database

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var migrations = []struct {
	Version int
	Name    string
	Up      func(context.Context, *mongo.Database) error
}{
	{
		Version: 1,
		Name:    "Create users collection",
		Up: func(ctx context.Context, db *mongo.Database) error {
			err := db.CreateCollection(ctx, "users")
			if err != nil {
				if cmdErr, ok := err.(mongo.CommandError); ok && cmdErr.Code == 48 {
					return nil
				}
				return err
			}

			_, err = db.Collection("users").Indexes().CreateMany(ctx, []mongo.IndexModel{
				{
					Keys:    bson.D{{Key: "email", Value: 1}},
					Options: options.Index().SetUnique(true),
				},
				{
					Keys: bson.D{{Key: "registerAt", Value: 1}},
				},
			})
			return err
		},
	},
	{
		Version: 2,
		Name:    "Create categories collection",
		Up: func(ctx context.Context, db *mongo.Database) error {
			err := db.CreateCollection(ctx, "categories")
			if err != nil {
				if cmdErr, ok := err.(mongo.CommandError); ok && cmdErr.Code == 48 {
					return nil
				}
				return err
			}

			_, err = db.Collection("categories").Indexes().CreateOne(ctx, mongo.IndexModel{
				Keys:    bson.D{{Key: "name", Value: 1}},
				Options: options.Index().SetUnique(true),
			})
			return err
		},
	},
	{
		Version: 3,
		Name:    "Create restaurants collection",
		Up: func(ctx context.Context, db *mongo.Database) error {
			err := db.CreateCollection(ctx, "restaurants")
			if err != nil {
				if cmdErr, ok := err.(mongo.CommandError); ok && cmdErr.Code == 48 {
					return nil
				}
				return err
			}

			_, err = db.Collection("restaurants").Indexes().CreateMany(ctx, []mongo.IndexModel{
				{
					Keys: bson.D{{Key: "name", Value: 1}},
				},
				{
					Keys: bson.D{{Key: "categoryId", Value: 1}},
				},
				{
					Keys: bson.D{
						{Key: "location.latitude", Value: "2dsphere"},
						{Key: "location.longitude", Value: "2dsphere"},
					},
				},
			})
			return err
		},
	},
	{
		Version: 4,
		Name:    "Create reviews collection",
		Up: func(ctx context.Context, db *mongo.Database) error {
			err := db.CreateCollection(ctx, "reviews")
			if err != nil {
				if cmdErr, ok := err.(mongo.CommandError); ok && cmdErr.Code == 48 {
					return nil
				}
				return err
			}

			_, err = db.Collection("reviews").Indexes().CreateMany(ctx, []mongo.IndexModel{
				{
					Keys: bson.D{
						{Key: "userId", Value: 1},
						{Key: "restaurantId", Value: 1},
					},
				},
				{
					Keys: bson.D{{Key: "createdAt", Value: 1}},
				},
			})
			return err
		},
	},
	{
		Version: 5,
		Name:    "Create nlp_results collection",
		Up: func(ctx context.Context, db *mongo.Database) error {
			err := db.CreateCollection(ctx, "nlp_results")
			if err != nil {
				if cmdErr, ok := err.(mongo.CommandError); ok && cmdErr.Code == 48 {
					return nil
				}
				return err
			}

			_, err = db.Collection("nlp_results").Indexes().CreateMany(ctx, []mongo.IndexModel{
				{
					Keys:    bson.D{{Key: "reviewId", Value: 1}},
					Options: options.Index().SetUnique(true),
				},
				{
					Keys: bson.D{{Key: "sentiment", Value: 1}},
				},
			})
			return err
		},
	},
	{
		Version: 6,
		Name:    "Create ratings collection",
		Up: func(ctx context.Context, db *mongo.Database) error {
			err := db.CreateCollection(ctx, "ratings")
			if err != nil {
				if cmdErr, ok := err.(mongo.CommandError); ok && cmdErr.Code == 48 {
					return nil
				}
				return err
			}

			_, err = db.Collection("ratings").Indexes().CreateMany(ctx, []mongo.IndexModel{
				{
					Keys:    bson.D{{Key: "restaurantId", Value: 1}},
					Options: options.Index().SetUnique(true),
				},
				{
					Keys: bson.D{{Key: "averageRating", Value: -1}},
				},
			})
			return err
		},
	},
	{
		Version: 7,
		Name:    "Create admins_logs collection",
		Up: func(ctx context.Context, db *mongo.Database) error {
			err := db.CreateCollection(ctx, "admins_logs")
			if err != nil {
				if cmdErr, ok := err.(mongo.CommandError); ok && cmdErr.Code == 48 {
					return nil
				}
				return err
			}

			_, err = db.Collection("admins_logs").Indexes().CreateMany(ctx, []mongo.IndexModel{
				{
					Keys: bson.D{{Key: "adminId", Value: 1}},
				},
				{
					Keys: bson.D{{Key: "actionType", Value: 1}},
				},
				{
					Keys: bson.D{{Key: "createdAt", Value: 1}},
				},
			})
			return err
		},
	},
	{
		Version: 8,
		Name:    "Create favorites collection",
		Up: func(ctx context.Context, db *mongo.Database) error {
			err := db.CreateCollection(ctx, "favorites")
			if err != nil {
				if cmdErr, ok := err.(mongo.CommandError); ok && cmdErr.Code == 48 {
					return nil
				}
				return err
			}

			_, err = db.Collection("favorites").Indexes().CreateMany(ctx, []mongo.IndexModel{
				{
					Keys: bson.D{
						{Key: "userId", Value: 1},
						{Key: "restaurantId", Value: 1},
					},
					Options: options.Index().SetUnique(true),
				},
				{
					Keys: bson.D{{Key: "addedAt", Value: 1}},
				},
			})
			return err
		},
	},
}

func RunMigrations(ctx context.Context) error {
	db := MongoDB.Database("restaurantdb_1")

	for _, migration := range migrations {
		log.Printf("Applying migration %d: %s\n", migration.Version, migration.Name)
		err := migration.Up(ctx, db)
		if err != nil {
			return fmt.Errorf("failed to apply migration %d: %v", migration.Version, err)
		}
		log.Printf("Successfully applied migration %d\n", migration.Version)
	}

	return nil
}
