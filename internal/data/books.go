package data

import (
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"
)

// By default, the keys in the JSON object are equal to the field names in the struct ( ID,
// CreatedAt, Title and so on).
type Book struct {
	ID        int64     `json:"id"`                       // Unique integer ID for the book
	CreatedAt time.Time `json:"-"`                        // Timestamp for when the book is added to our database, "-" directive, hidden in response
	Title     string    `json:"title"`                    // Book title
	Year      int32     `json:"year,omitempty"`           // Book release year, "omitempty" - hide from response if empty
	Cost      int32     `json:"runtime,omitempty,string"` // Book runtime (in minutes), "string" - convert int to string
	Genres    []string  `json:"genres,omitempty"`         // Slice of genres for the book (romance, comedy, etc.)
	Version   int32     `json:"version"`                  // The version number starts at 1 and will be incremented each
	Amount    int       `json:"amount"`

	// time the book information is updated
}

// Define a BookModel struct type which wraps a sql.DB connection pool.
type BookModel struct {
	DB *sql.DB
}

// method for inserting a new record in the books table.
func (m BookModel) Insert(book *Book) error {
	query := `
		INSERT INTO books(title, year, cost, genres)
		VALUES ($1, $2, $3, $4)
		
		RETURNING id, created_at, version`

	return m.DB.QueryRow(query, &book.Title, &book.Year, &book.Cost, pq.Array(&book.Genres)).Scan(&book.ID, &book.CreatedAt, &book.Version)
}

// method for fetching a specific record from the books table.
func (m BookModel) Get(id int64) (*Book, error) {
	if id < 1 {
		return nil, ErrRecordNotFound
	}

	query := `
		SELECT *
		FROM books
		WHERE id = $1`

	var book Book

	err := m.DB.QueryRow(query, id).Scan(
		&book.ID,
		&book.CreatedAt,
		&book.Title,
		&book.Year,
		&book.Cost,
		pq.Array(&book.Genres),
		&book.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &book, nil

}

// method for updating a specific record in the books table.
func (m BookModel) Update(book *Book) error {
	query := `
		UPDATE books
		SET title = $1, year = $2, runtime = $3, genres = $4, version = version + 1
		WHERE id = $5
		RETURNING version`

	args := []interface{}{
		book.Title,
		book.Year,
		book.Cost,
		pq.Array(book.Genres),
		book.ID,
	}

	return m.DB.QueryRow(query, args...).Scan(&book.Version)
}

// method for deleting a specific record from the books table.
func (m BookModel) Delete(id int64) error {
	if id < 1 {
		return ErrRecordNotFound
	}
	// Construct the SQL query to delete the record.
	query := `
		DELETE FROM books
		WHERE id = $1`

	result, err := m.DB.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}
