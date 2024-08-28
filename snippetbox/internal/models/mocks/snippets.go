package mocks

import (
	"time"

	"github.com/sotnikea/Go_Learn/tree/main/snippetbox/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var mockSnippet = models.Snippet{
	ID:      "111111111111111111111111",
	Title:   "An old silent pond",
	Content: "An old silent pond...",
	Created: time.Now(),
	Expires: time.Now(),
}

type SnippetModel struct{}

func (m *SnippetModel) Insert(title string, content string, expires int) (interface{}, error) {
	objectID, _ := primitive.ObjectIDFromHex("222222222222222222222222")
	return objectID, nil
}

func (m *SnippetModel) Get(id string) (models.Snippet, error) {
	switch id {
	case "111111111111111111111111":
		return mockSnippet, nil
	default:
		return models.Snippet{}, models.ErrNoRecord
	}
}

func (m *SnippetModel) Latest() ([]models.Snippet, error) {
	return []models.Snippet{mockSnippet}, nil
}
