package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Movie struct {
	Id           primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	Title        string             `json:"title"`
	ReleasedYear int32              `json:"releasedYear"`
	Rating       float32            `json:"rating"`
	Genres       []string           `json:"genres"`
}

// Needed for post request
type PostMovie struct {
	Title        string   `json:"title"`
	ReleasedYear int32    `json:"releasedYear"`
	Rating       float32  `json:"rating"`
	Genres       []string `json:"genres"`
}

type UpdateMovie struct {
	Rating float32  `json:"rating"`
	Genres []string `json:"genres"`
}

type FindMovieByTitle struct {
	Title string `json:"title"`
}

type FindMovie struct {
	Id                        string   `json:"_id"`
	Title                     string   `json:"title"`
	ReleasedYearInferiorRange string   `json:"releasedYearInferiorRange"`
	ReleasedYearSuperiorRange string   `json:"releasedYearSuperiorRange"`
	RatingInferiorRange       string   `json:"ratingInferiorRange"`
	RatingSuperiorRange       string   `json:"ratingSuperiorRange"`
	Genres                    []string `json:"genres"`
}

type MovieList struct {
	Movies []*Movie `json:"movies"`
}

var (
	// Movie attributes
	MovieId           string = "_id"
	MovieTitle        string = "title"
	MovieReleasedYear string = "releasedYear"
	MovieRating       string = "rating"
	MovieGenres       string = "genres"
)
