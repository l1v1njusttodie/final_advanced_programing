package data

import (
	"database/sql"
	"errors"
)

// Define a custom ErrRecordNotFound error. We'll return this from our Get() method when
// looking up a movie that doesn't exist in our database.
var (
	ErrRecordNotFound = errors.New("record (row, entry) not found")
	ErrEditConflict   = errors.New("edit conflict")
)

// Create a Models struct which wraps the BookModel
// kind of enveloping
type Models struct {
	Books  BookModel
	Users  UserModel
	Tokens TokenModel // used to generate activation tokens
}

// method which returns a Models struct containing the initialized BookModel.
func NewModels(db *sql.DB) Models {
	return Models{
		Books:  BookModel{DB: db},
		Users:  UserModel{DB: db},
		Tokens: TokenModel{DB: db}, // new TokenModel initilization
	}
}
