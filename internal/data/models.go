package data

import (
	"database/sql"
	"errors"
)

// Define a custom ErrRecordNotFound error.
var (
	ErrRecordNotFound = errors.New("models: record not found")
	ErrEditConflict   = errors.New("models: edit conflict")
)

// Create a Models struct which wraps the MovieModel.
type Models struct {
	// Movies MovieModel
	// Set the Movies field to be an interface containing the methods that both the
	// 'real' model and mock model need to support
	Movies interface {
		Insert(movie *Movie) error
		Get(id int64) (*Movie, error)
		Update(movie *Movie) error
		Delete(id int64) error
	}
}

// For ease of use, we also add a New() method which returns a Models struct containging
// a MovieModel.
func NewModels(db *sql.DB) Models {
	return Models{
		Movies: MovieModel{DB: db},
	}

}

// Create a helper function which rturns a Models instance containing the mock models only.
func NewMockModels() Models {
	return Models{
		Movies: MockMovieModel{},
	}
}
