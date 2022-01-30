package helpers

import (
	"movieapi/app/models"
	"strconv"
	"strings"

	"github.com/eefret/gomdb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func MapPostMovieToBson(movie *models.PostMovie) bson.M {
	return bson.M{
		models.MovieTitle:        movie.Title,
		models.MovieReleasedYear: movie.ReleasedYear,
		models.MovieRating:       movie.Rating,
		models.MovieGenres:       movie.Genres,
	}
}

func MapPostMovieToMovie(id primitive.ObjectID, movie *models.PostMovie) *models.Movie {

	return &models.Movie{
		Id:           id,
		Title:        movie.Title,
		ReleasedYear: movie.ReleasedYear,
		Rating:       movie.Rating,
		Genres:       movie.Genres,
	}
}

// MapMovieToPostMovie
func MapMovieToPostMovie(movie *models.Movie) models.PostMovie {
	return models.PostMovie{
		Title:        movie.Title,
		ReleasedYear: movie.ReleasedYear,
		Rating:       movie.Rating,
		Genres:       movie.Genres,
	}
}

func MapFindMovieByTitleToFindMovie(findMovieByTitle models.FindMovieByTitle) models.FindMovie {
	return models.FindMovie{
		Title: findMovieByTitle.Title,
	}
}

func MapMovieResultToMovie(movieResult *gomdb.MovieResult) (*models.PostMovie, error) {
	year := movieResult.Year
	// Convert year from string to int32
	yearInt32, err := strconv.ParseInt(year, 10, 32)
	if err != nil {
		return nil, err
	}

	rating := movieResult.ImdbRating
	// Convert rating from string to float64
	ratingFloat32, err := strconv.ParseFloat(rating, 32)
	if err != nil {
		return nil, err
	}

	genre := movieResult.Genre
	// Convert genre from string to []string separated by comma
	genreArray := strings.Split(genre, ",")
	// Remove the first space from each string element in the array
	for i, v := range genreArray {
		genreArray[i] = strings.TrimSpace(v)
	}

	return &models.PostMovie{
		Title:        movieResult.Title,
		ReleasedYear: int32(yearInt32),
		Rating:       float32(ratingFloat32),
		Genres:       genreArray,
	}, nil
}
