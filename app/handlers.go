package app

import (
	"context"
	"errors"
	"fmt"
	"movieapi/app/models"
	"net/http"

	"movieapi/app/config"
	"movieapi/app/helpers"

	"github.com/eefret/gomdb"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// IndexHandler is the handler for the index page
func (a *App) IndexHandler() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintf(writer, "Welcome to the home page!")
	}
}

func (a *App) CreateMovieHandler(movie *models.PostMovie) (*models.Movie, error) {
	bsonPostMovie := helpers.MapPostMovieToBson(movie)

	// Save in DB
	insertOneResult, err := a.DB.InsertOne(config.DbName, config.MovieCollection, bsonPostMovie)
	if err != nil {
		return &models.Movie{}, err
	}

	insertedID := insertOneResult.InsertedID.(primitive.ObjectID)

	response := helpers.MapPostMovieToMovie(insertedID, movie)

	return response, nil
}

func (a *App) UpdateMovieHandler() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		updateMovie := models.UpdateMovie{}
		err := helpers.Parse(writer, request, &updateMovie)
		if err != nil {
			fmt.Printf("cannot parse movie body. err=%v \n", err)
			helpers.SendResponse(writer, request, nil, http.StatusBadRequest)
			return
		}

		vars := mux.Vars(request)
		id, ok := vars["id"]
		if !ok {
			fmt.Printf("title is missing in parameters")
		}

		movieFilters := models.FindMovie{
			Id: id,
		}
		// Get movie from DB
		movieList, err := a.findMovieAux(movieFilters)
		if err != nil {
			fmt.Printf("cannot find movie in DB. err=%v \n", err)
			helpers.SendResponse(writer, request, nil, http.StatusBadRequest)
			return
		}

		if len(movieList.Movies) == 0 {
			err = errors.New("movie not found")
			fmt.Printf("cannot find movie in DB. err=%v \n", err)
			helpers.SendResponse(writer, request, nil, http.StatusBadRequest)
			return
		}

		firtsMovie := movieList.Movies[0]
		firtsMovie.Genres = updateMovie.Genres
		firtsMovie.Rating = updateMovie.Rating

		movieToUpdate := helpers.MapMovieToPostMovie(firtsMovie)

		objId, _ := primitive.ObjectIDFromHex(id)

		filter := bson.M{"_id": bson.M{"$eq": objId}}

		// Update in DB
		updateResult, err := a.DB.UpdateOne(config.DbName, config.MovieCollection, filter, movieToUpdate)
		if err != nil {
			fmt.Printf("Cannot update movie in DB. err=%v \n", err)
			helpers.SendResponse(writer, request, nil, http.StatusInternalServerError)
			return
		}

		if updateResult.ModifiedCount == 0 {
			fmt.Printf("Nothing was modified in DB. err=%v \n", err)
			helpers.SendResponse(writer, request, nil, http.StatusNotFound)
			return
		}

		helpers.SendResponse(writer, request, movieToUpdate, http.StatusOK)
	}
}

func (a *App) FindMovieByTitleHandler() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		vars := mux.Vars(request)
		title, ok := vars["title"]
		if !ok {
			fmt.Printf("title is missing in parameters")
		}

		movieTitle := models.FindMovieByTitle{}
		if title == "" {
			err := errors.New("title is required")
			fmt.Printf("cannot parse movie body. err=%v \n", err)
			helpers.SendResponse(writer, request, nil, http.StatusBadRequest)
			return
		}

		movieTitle.Title = title
		movieFilter := helpers.MapFindMovieByTitleToFindMovie(movieTitle)

		movieList, err := a.findMovieAux(movieFilter)
		if err != nil {
			fmt.Printf("cannot find movie in DB. err=%v \n", err)
			helpers.SendResponse(writer, request, nil, http.StatusInternalServerError)
			return
		}

		if len(movieList.Movies) == 0 {
			movie, err := a.findMovieByTitleFromExternalApi(movieTitle.Title)
			if err != nil {
				fmt.Printf("cannot find movie in external API. err=%v \n", err)
				helpers.SendResponse(writer, request, nil, http.StatusInternalServerError)
				return
			}

			// Create movie in DB
			movieCreated, err := a.CreateMovieHandler(movie)
			if err != nil {
				fmt.Printf("cannot create movie in DB. err=%v \n", err)
				helpers.SendResponse(writer, request, nil, http.StatusInternalServerError)
				return
			}

			helpers.SendResponse(writer, request, movieCreated, http.StatusOK)
			return
		}
		movieResponse := &models.Movie{}
		if len(movieList.Movies) > 0 {
			movieResponse = movieList.Movies[0]
		}

		helpers.SendResponse(writer, request, movieResponse, http.StatusOK)
	}
}

func (a *App) findMovieByTitleFromExternalApi(movieTitle string) (*models.PostMovie, error) {
	api := gomdb.Init(config.OmbdApiKey)
	query := &gomdb.QueryData{
		Title: movieTitle,
	}
	movieResult, err := api.MovieByTitle(query)
	if err != nil {
		return &models.PostMovie{}, err
	}

	movie, err := helpers.MapMovieResultToMovie(movieResult)
	if err != nil {
		return &models.PostMovie{}, err
	}

	return movie, nil
}

func (a *App) FindMoviesHandler() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

		id := request.URL.Query().Get("id")
		title := request.URL.Query().Get("title")
		releasedYearInferiorRange := request.URL.Query().Get("releasedYearInferiorRange")
		releasedYearSuperiorRange := request.URL.Query().Get("releasedYearSuperiorRange")
		ratingInferiorRange := request.URL.Query().Get("ratingInferiorRange")
		ratingSuperiorRange := request.URL.Query().Get("ratingSuperiorRange")
		// Get all genres from query
		genres := request.URL.Query()["genres"]

		movieFilters := models.FindMovie{
			Id:                        id,
			Title:                     title,
			ReleasedYearInferiorRange: releasedYearInferiorRange,
			ReleasedYearSuperiorRange: releasedYearSuperiorRange,
			RatingInferiorRange:       ratingInferiorRange,
			RatingSuperiorRange:       ratingSuperiorRange,
			Genres:                    genres,
		}

		movieList := models.MovieList{}
		var err error
		movieList, err = a.findMovieAux(movieFilters)
		if err != nil {
			fmt.Printf("cannot find movie in DB. err=%v \n", err)
			helpers.SendResponse(writer, request, nil, http.StatusInternalServerError)
			return
		}

		helpers.SendResponse(writer, request, movieList, http.StatusOK)
	}
}

func (a *App) findMovieAux(movieFilters models.FindMovie) (models.MovieList, error) {
	filter, err := a.setMovieFilter(movieFilters)
	if err != nil {
		fmt.Printf("cannot set movie filter. err=%v \n", err)
		return models.MovieList{}, err
	}
	findOptions := options.Find()

	// Find in DB
	cursor, err := a.DB.Query(config.DbName, config.MovieCollection, &filter, findOptions)
	if err != nil {
		fmt.Printf("cannot find movie in DB. err=%v \n", err)
		return models.MovieList{}, err

	}
	var movieList models.MovieList
	var movies []*models.Movie

	// Iterate cursor
	for cursor.Next(context.TODO()) {
		var movieTemp models.Movie
		err := cursor.Decode(&movieTemp)
		if err != nil {
			fmt.Printf("cannot decode movie. err=%v \n", err)
			return models.MovieList{}, err

		}
		movies = append(movies, &movieTemp)
	}

	movieList.Movies = movies

	return movieList, nil
}

func (a *App) setMovieFilter(movieFilters models.FindMovie) (*bson.M, error) {
	filter := make(bson.M)

	err := helpers.SetFilterId(&filter, models.MovieId, movieFilters.Id)
	if err != nil {
		fmt.Printf("Cannot set filter id. err=%v \n", err)
		return nil, err
	}

	helpers.SetFilter(&filter, models.MovieTitle, movieFilters.Title)

	err = helpers.SetRangeFilter(&filter, models.MovieReleasedYear, movieFilters.ReleasedYearInferiorRange, movieFilters.ReleasedYearSuperiorRange, 32)
	if err != nil {
		fmt.Printf("Cannot set range filter. err=%v \n", err)
		return nil, err
	}

	// Rating
	err = helpers.SetRangeFilter(&filter, models.MovieRating, movieFilters.RatingInferiorRange, movieFilters.RatingSuperiorRange, 32)
	if err != nil {
		fmt.Printf("Cannot set range filter. err=%v \n", err)
		return nil, err
	}

	// Genre
	err = helpers.SetFilterArrayString(&filter, models.MovieGenres, movieFilters.Genres)
	if err != nil {
		fmt.Printf("Cannot set filter array string. err=%v \n", err)
		return nil, err
	}

	return &filter, nil
}
