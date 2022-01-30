package models

type Movie struct {
	Id           string   `json:"_id"`
	Title        string   `json:"title"`
	ReleasedYear int32    `json:"releasedYear"`
	Rating       float32  `json:"rating"`
	Genres       []string `json:"genres"`
}

// Needed for post request
type CreateMovie struct {
	Title        string   `json:"title"`
	ReleasedYear int32    `json:"releasedYear"`
	Rating       float32  `json:"rating"`
	Genres       []string `json:"genres"`
}

type UpdateMovie struct {
	Id     string   `json:"_id"`
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
