package models

import (
	"context"
	"errors"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

// Define a new User struct. Notice how the field names and types align
// with the columns in the database "users" table?
type User struct {
	ID             string
	Name           string
	Email          string
	HashedPassword []byte
	Created        time.Time
}

// Define a new UserModel struct which wraps a database connection pool.
type UserModel struct {
	DB *mongo.Client
}

// We'll use the Insert method to add a new record to the "users" table.
func (m *UserModel) Insert(name, email, password string) error {
	// Create a bcrypt hash of the plain-text password.
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return err
	}

	// Create context for operation
	//ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	//defer cancel()

	// Prepare document for insert
	doc := bson.D{
		{Key: "name", Value: name},
		{Key: "email", Value: email},
		{Key: "hashed_password", Value: string(hashedPassword)},
		{Key: "created", Value: time.Now()},
	}

	// Get collection for insert operation
	collection := m.DB.Database("snippetbox").Collection("users")

	// Insert document into collection
	_, err = collection.InsertOne(context.Background(), doc)
	if err != nil {
		// If this returns an error, we use the errors.As() function to check
		// whether the error has the type mongo.WriteException. If it does, the
		// error will be assigned to the mongoWriteException variable. We can then check
		// whether or not the error relates to our users_uc_email key by
		// checking if the error code equals 11000 and the contents of the error
		// message string. If it does, we return an ErrDuplicateEmail error.
		var mongoWriteException mongo.WriteException
		if errors.As(err, &mongoWriteException) {
			for _, we := range mongoWriteException.WriteErrors {
				if we.Code == 11000 && strings.Contains(we.Message, "email") {
					return ErrDuplicateEmail
				}
			}
		}
		return err
	}

	return nil
}

// We'll use the Authenticate method to verify whether a user exists with
// the provided email address and password. This will return the relevant
// user ID if they do.
func (m *UserModel) Authenticate(email, password string) (interface{}, error) {
	// Retrieve the id and hashed password associated with the given email. If
	// no matching email exists we return the ErrInvalidCredentials error.
	var result struct {
		ID             interface{} `bson:"_id"`
		HashedPassword []byte      `bson:"hashed_password"`
	}

	// Create request for searching document
	filter := bson.D{{Key: "email", Value: email}}

	// Execute request for the collection and find one document
	err := m.DB.Database("snippetbox").Collection("users").FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return 0, ErrInvalidCredentials
		}
		return 0, err
	}

	// Check whether the hashed password and plain-text password provided match.
	// If they don't, we return the ErrInvalidCredentials error.
	err = bcrypt.CompareHashAndPassword(result.HashedPassword, []byte(password))
	if err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return 0, ErrInvalidCredentials
		}
		return 0, err
	}

	// Otherwise, the password is correct. Return the user ID.
	return result.ID, nil
}

// We'll use the Exists method to check if a user exists with a specific ID.
func (m *UserModel) Exists(id string) (bool, error) {
	// Convert id to ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return false, err
	}

	// Create filter for searching user by ID
	filter := bson.D{{Key: "_id", Value: objectID}}

	// Execute MongoDB query with limit 1
	count, err := m.DB.Database("snippetbox").Collection("users").CountDocuments(context.Background(), filter, options.Count().SetLimit(1))
	if err != nil {
		return false, err
	}

	// Check user existing
	exists := count > 0

	return exists, err
}
