package models

import (
	"database/sql"
	"errors"
	"time"
)

type Post struct {
	ID         int
	Title      string
	Content    string
	Excerpt    string
	Author     string
	Created_at time.Time
	Updated_at time.Time
	Slug       string
}

type PostModel struct {
	DB *sql.DB
}

func (m *PostModel) Insert(title string, content string, excerpt string, author string, slug string) (string, error) {
	stmt := `INSERT INTO blog_posts (title, content, excerpt, author, slug)
	         VALUES (?, ?, ?, ?, ?)`

	_, err := m.DB.Exec(stmt, title, content, excerpt, author, slug)
	if err != nil {
		return "", err
	}

	return slug, nil
}

func (m *PostModel) Get(slug string) (Post, error) {
	stmt := `SELECT id, title, content, author, created_at, updated_at, slug FROM blog_posts
WHERE slug = ?`

	row := m.DB.QueryRow(stmt, slug)

	var p Post

	err := row.Scan(&p.ID, &p.Title, &p.Content, &p.Author, &p.Created_at, &p.Updated_at, &p.Slug)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Post{}, ErrNoRecord
		} else {
			return Post{}, err
		}
	}
	// If everything went OK, then return the filled Snippet struct.
	return p, nil
}

func (m *PostModel) Latest() ([]Post, error) {
	stmt := `SELECT id, title, author, excerpt, created_at, updated_at, slug FROM blog_posts ORDER BY id DESC LIMIT 10`

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var posts []Post

	for rows.Next() {
		var p Post

		err = rows.Scan(&p.ID, &p.Title, &p.Author, &p.Excerpt, &p.Created_at, &p.Updated_at, &p.Slug)
		if err != nil {
			return nil, err
		}

		posts = append(posts, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}
