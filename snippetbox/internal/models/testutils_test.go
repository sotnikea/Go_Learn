package models

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type InsertMany struct {
	Documents  []interface{} `json:"documents"`
	Collection string        `json:"collection"`
}

type InsertOne struct {
	Document   interface{} `json:"document"`
	Collection string      `json:"collection"`
}

type CreateIndexes struct {
	Indexes []struct {
		Key    bson.D `json:"key"`
		Name   string `json:"name"`
		Unique bool   `json:"unique"`
	} `json:"indexes"`
	Collection string `json:"collection"`
}

func newTestDB(t *testing.T) *mongo.Database {
	// Establish connection to MongoDB
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://test_web:pass@localhost:27017"))
	if err != nil {
		t.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Get database for testing
	db := client.Database("test_snippetbox")

	if err := db.CreateCollection(context.TODO(), "init"); err != nil {
		t.Fatalf("Can not create initial collection: %v", err)
	}

	// Execute setup script
	setupScript, err := os.ReadFile("./testdata/setup.json")
	if err != nil {
		dbDrop(t, db, client)
		t.Fatalf("Failed to read setup script: %v", err)
	}

	if err = executeScript(db, setupScript); err != nil {
		dbDrop(t, db, client)
		t.Fatalf("Failed to execute setup script: %v", err)
	}

	// Add cleaning
	t.Cleanup(func() {
		defer dbDrop(t, db, client)

		teardownScript, err := os.ReadFile("./testdata/teardown.json")

		if err != nil {
			t.Fatalf("Failed to read teardown script: %v", err)
		}

		if err = executeScript(db, teardownScript); err != nil {
			t.Fatalf("Failed to execute teardown script: %v", err)
		}
	})

	// Return connection to database
	return db
}

func dbDrop(t *testing.T, db *mongo.Database, client *mongo.Client) {
	if err := db.Drop(context.TODO()); err != nil {
		t.Fatalf("Failed to drop database: %v", err)
	}
	if err := client.Disconnect(context.TODO()); err != nil {
		t.Fatalf("Failed to disconnect MongoDB client: %v", err)
	}
}

func executeScript(db *mongo.Database, script []byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Parse to JSON
	var commands []map[string]interface{}
	if err := json.Unmarshal(script, &commands); err != nil {
		return fmt.Errorf(" Error unmarshalling JSON: %w", err)
	}

	// Execute each command
	for _, cmd := range commands {

		// Handle insertMany
		if insertManyCmd, ok := cmd["insertMany"]; ok {
			insertManyData := InsertMany{}
			mapData, _ := json.Marshal(insertManyCmd)
			if err := json.Unmarshal(mapData, &insertManyData); err != nil {
				return fmt.Errorf("Error unmarshalling insertMany command: %w", err)
			}

			// Check need of _id converting
			for i, doc := range insertManyData.Documents {
				if docMap, ok := doc.(map[string]interface{}); ok {
					if idStr, ok := docMap["_id"].(string); ok {
						objectId, err := primitive.ObjectIDFromHex(idStr)
						if err == nil {
							docMap["_id"] = objectId
						} else {
							return fmt.Errorf("Error converting _id to ObjectID for document %d: %w", i, err)
						}
					}
				}
			}

			collection := db.Collection(insertManyData.Collection)
			_, err := collection.InsertMany(ctx, insertManyData.Documents)
			if err != nil {
				return fmt.Errorf("Error executing insertMany: %w", err)
			}
		}

		// Handle insertOne
		if insertOneCmd, ok := cmd["insertOne"]; ok {
			insertOneData := InsertOne{}
			mapData, _ := json.Marshal(insertOneCmd)
			if err := json.Unmarshal(mapData, &insertOneData); err != nil {
				return fmt.Errorf("Error unmarshalling insertOne command: %w", err)
			}

			// Convert _id from string to ObjectID
			if docMap, ok := insertOneData.Document.(map[string]interface{}); ok {
				if idStr, ok := docMap["_id"].(string); ok {
					objectId, err := primitive.ObjectIDFromHex(idStr)
					if err == nil {
						docMap["_id"] = objectId
					} else {
						return fmt.Errorf("Error converting _id to ObjectID: %w", err)
					}
				}
			}

			collection := db.Collection(insertOneData.Collection)
			_, err := collection.InsertOne(ctx, insertOneData.Document)
			if err != nil {
				return fmt.Errorf("Error executing insertOne: %w", err)
			}
		}

		// Handle createIndexes
		if createIndexesCmd, ok := cmd["createIndexes"]; ok {
			createIndexesData := CreateIndexes{}
			mapData, _ := json.Marshal(createIndexesCmd)
			if err := json.Unmarshal(mapData, &createIndexesData); err != nil {
				return fmt.Errorf("Error unmarshalling createIndexes command: %w", err)
			}
			collection := db.Collection(createIndexesData.Collection)
			indexes := []mongo.IndexModel{}

			for _, index := range createIndexesData.Indexes {
				indexModel := mongo.IndexModel{
					Keys: index.Key,
					Options: options.Index().
						SetName(index.Name).
						SetUnique(index.Unique),
				}
				indexes = append(indexes, indexModel)
			}

			if len(indexes) > 0 {
				_, err := collection.Indexes().CreateMany(ctx, indexes)
				if err != nil {
					return fmt.Errorf("Error creating indexes: %w", err)
				}
			} else {
				return fmt.Errorf("No indexes found to create")
			}
		}
	}

	return nil
}
