package mocks

import (
	"time"

	"snippetbox.bogdandev.de/internal/models"
)

var mockSnippet = &models.Snippet{
	Id:      1,
	Title:   "An old silent pond",
	Content: "An old silent pond...",
	Created: time.Now(),
	Expires: time.Now(),
}

type SnippetModel struct{}

func (m *SnippetModel) Insert(title, content string, expiresInDays int) (int64, error) {
	return 2, nil
}

func (m *SnippetModel) Get(id int64) (*models.Snippet, error) {
	switch id {
	case 1:
		return mockSnippet, nil
	default:
		return nil, models.ErrNoRecord
	}
}

func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	return []*models.Snippet{mockSnippet}, nil
}
