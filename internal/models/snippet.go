package models

import (
	"database/sql"
	"errors"
	"time"
)

type Snippet struct {
	Id      int64
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}

type SnippetModel struct {
	DB *sql.DB
}

func (sm *SnippetModel) Insert(title string, content string, expiresInDays int) (int64, error) {
	now := time.Now().UTC()
	expires := now.AddDate(0, 0, expiresInDays)

	res, err := sm.DB.Exec(
		`INSERT INTO snippets (title, content, created, expires)
			VALUES (?, ?, ?, ?)`,
		title, content, now, expires,
	)

	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()

	if err != nil {
		return 0, err
	}

	return id, nil
}

func (sm *SnippetModel) Get(id int64) (*Snippet, error) {
	row := sm.DB.QueryRow(`SELECT id, title, content, created, expires FROM snippets WHERE id = ?`, id)

	var snippet Snippet

	err := row.Scan(
		&snippet.Id,
		&snippet.Title,
		&snippet.Content,
		&snippet.Created,
		&snippet.Expires,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Snippet doesn't exist
			return nil, ErrNoRecord
		}
		return nil, err
	}

	if time.Now().After(snippet.Expires) {
		// Snippet is expired
		return nil, ErrNoRecord
	}

	return &snippet, nil
}

func (sm *SnippetModel) Latest() ([]*Snippet, error) {
	return nil, nil
}
