package models

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SnippetModelInterface interface {
	Insert(title string, content string, expires int) (interface{}, error)
	Get(id string) (Snippet, error)
	Latest() ([]Snippet, error)
}

// Define a Snippet type to hold the data for an individual snippet
type Snippet struct {
	ID      string `bson:"_id,omitempty"`
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

// Define a SnippetModel type which wraps a mongo.Client connection pool
type SnippetModel struct {
	DB *mongo.Database
}

// This will insert a new snippet into the database.
func (m *SnippetModel) Insert(title string, content string, expires int) (interface{}, error) {
	// Create context for operation
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Prepare document for insert
	doc := bson.D{
		{Key: "title", Value: title},
		{Key: "content", Value: content},
		{Key: "created", Value: time.Now()},
		{Key: "expires", Value: time.Now().Add(time.Duration(expires) * time.Hour * 24)},
	}

	// Get collection for insert operation
	collection := m.DB.Collection("snippets")

	// Insert document into collection
	result, err := collection.InsertOne(ctx, doc)
	if err != nil {
		return nil, err
	}

	// Return id of inserted document
	return result.InsertedID, nil

}

// This will return a specific snippet based on its id.
func (m *SnippetModel) Get(id string) (Snippet, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Create empty Snippet for saving results
	var s Snippet

	// Transform id to ObjectID
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return Snippet{}, err
	}

	// Create request for searching document
	filter := bson.D{
		{Key: "_id", Value: objID},
		{Key: "expires", Value: bson.D{{Key: "$gt", Value: time.Now()}}},
	}

	// Execute request for the collection and find one document
	err = m.DB.Collection("snippets").FindOne(ctx, filter).Decode(&s)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return Snippet{}, ErrNoRecord
		}
		return Snippet{}, err
	}

	// Return result
	return s, nil
}

// This will return the 10 most recently created snippets
func (m *SnippetModel) Latest() ([]Snippet, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Create empty array for snippets
	var snippets []Snippet

	// Search only not expired document
	filter := bson.D{
		{Key: "expires", Value: bson.D{{Key: "$gt", Value: time.Now()}}},
	}

	// Get last 10 documents ordered by the creation time
	opts := options.Find().SetSort(bson.D{{Key: "created", Value: -1}}).SetLimit(10)

	// Execute request
	cursor, err := m.DB.Collection("snippets").Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}

	// Close cursor
	defer cursor.Close(ctx)

	// Decode all documents from cursor to Snippet structure
	for cursor.Next(ctx) {
		var s Snippet
		if err := cursor.Decode(&s); err != nil {
			return nil, err
		}
		snippets = append(snippets, s)
	}

	//Check for errors while get documents from collection
	if err := cursor.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}
