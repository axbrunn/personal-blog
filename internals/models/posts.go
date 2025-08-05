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
	Author     string
	Created_at time.Time
	Updated_at time.Time
	Slug       string
}

type PostModel struct {
	DB *sql.DB
}

// This will insert a new post into the database.
func (m *PostModel) Insert(title string, content string, expires int) (int, error) {
	return 0, nil
}

func (m *PostModel) Get(id int) (Post, error) {
	// Write the SQL statement we want to execute. Again, I've split it over two
	// lines for readability.
	stmt := `SELECT id, title, content, author, created_at, updated_at, slug FROM blog_posts
WHERE id = ?`
	// Use the QueryRow() method on the connection pool to execute our
	// SQL statement, passing in the untrusted id variable as the value for the
	// placeholder parameter. This returns a pointer to a sql.Row object which
	// holds the result from the database.
	row := m.DB.QueryRow(stmt, id)
	// Initialize a new zeroed Snippet struct.
	var p Post
	// Use row.Scan() to copy the values from each field in sql.Row to the
	// corresponding field in the Snippet struct. Notice that the arguments
	// to row.Scan are *pointers* to the place you want to copy the data into,
	// and the number of arguments must be exactly the same as the number of
	// columns returned by your statement.
	err := row.Scan(&p.ID, &p.Title, &p.Content, &p.Author, &p.Created_at, &p.Updated_at, &p.Slug)
	if err != nil {
		// If the query returns no rows, then row.Scan() will return a
		// sql.ErrNoRows error. We use the errors.Is() function check for that
		// error specifically, and return our own ErrNoRecord error
		// instead (we'll create this in a moment).
		if errors.Is(err, sql.ErrNoRows) {
			return Post{}, ErrNoRecord
		} else {
			return Post{}, err
		}
	}
	// If everything went OK, then return the filled Snippet struct.
	return p, nil
}

// This will return the 10 most recently created posts.
func (m *PostModel) Latest() ([]Post, error) {
	return nil, nil
}
