package models

import (
	"testing"

	"github.com/sotnikea/Go_Learn/tree/main/snippetbox/internal/assert"
)

func TestUserModelExists(t *testing.T) {
	// Skip the test if the "-short" flag is provided when running the test
	if testing.Short() {
		t.Skip("models: skipping integration test")
	}

	// Set up a suite of table-driven tests and expected results
	tests := []struct {
		name   string
		userID string
		want   bool
	}{
		{
			name:   "Valid ID",
			userID: "111111111111111111111111",
			want:   true,
		},

		{
			name:   "Zero ID",
			userID: "",
			want:   false,
		},
		{
			name:   "Non-existent ID",
			userID: "222222222222222222222222",
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Call the newTestDB() helper function to get a connection pool to our test database
			db := newTestDB(t)

			// Create a new instance of the UserModel
			m := UserModel{db}

			// Call the UserModel.Exists() method and check that the return
			// value and error match the expected values for the sub-test
			exists, err := m.Exists(tt.userID)

			assert.Equal(t, exists, tt.want)
			assert.NilError(t, err)
		})
	}
}
