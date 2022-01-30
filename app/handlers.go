package app

import (
	"context"
	"errors"
	"fmt"
	"movieapi/app/models"
	"net/http"

	"log"
	"movieapi/app/config"
	"movieapi/app/helpers"

	"github.com/eefret/gomdb"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// IndexHandler is the handler for the index page
func (a *App) IndexHandler() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintf(writer, "Welcome to the home page!")
	}
}

func (a *App) CreateMovieHandler(movie *models.CreateMovie) (*models.Movie, error) {
	bsonPostMovie := helpers.MapPostMovieToBson(movie)

	// Save in DB
	insertOneResult, err := a.DB.InsertOne(config.DbName, config.MovieCollection, bsonPostMovie)
	if err != nil {
		log.Printf("Cannot save movie in DB. err=%v \n", err)
		return &models.Movie{}, err
	}

	insertedID := insertOneResult.InsertedID.(primitive.ObjectID).Hex()

	response := helpers.MapPostMovieToMovie(insertedID, movie)

	return response, nil
}

func (a *App) UpdateMovieHandler() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		updateMovie := models.UpdateMovie{}
		err := helpers.Parse(writer, request, &updateMovie)
		if err != nil {
			log.Printf("Cannot parse movie body. err=%v \n", err)
			helpers.SendResponse(writer, request, nil, http.StatusBadRequest)
			return
		}

		movieFilters := models.FindMovie{
			Id: updateMovie.Id,
		}
		// Get movie from DB
		movieList, err := a.findMovieAux(movieFilters)
		if err != nil {
			log.Printf("Cannot find movie in DB. err=%v \n", err)
			helpers.SendResponse(writer, request, nil, http.StatusBadRequest)
			return
		}

		firtsMovie := movieList.Movies[0]
		movieToUpdate := helpers.MapMovieToPostMovie(firtsMovie)

		if len(movieList.Movies) == 0 {
			err := errors.New("Movie not found")
			log.Printf("Cannot find movie in DB. err=%v \n", err)
			helpers.SendResponse(writer, request, nil, http.StatusBadRequest)
			return
		}

		objId, _ := primitive.ObjectIDFromHex(updateMovie.Id)

		filter := bson.M{"_id": bson.M{"$eq": objId}}

		// Update in DB
		updateResult, err := a.DB.UpdateOne(config.DbName, config.MovieCollection, filter, movieToUpdate)
		if err != nil {
			log.Printf("Cannot update movie in DB. err=%v \n", err)
			helpers.SendResponse(writer, request, nil, http.StatusInternalServerError)
			return
		}

		if updateResult.ModifiedCount == 0 {
			log.Printf("Nothing was modified in DB. err=%v \n", err)
			helpers.SendResponse(writer, request, nil, http.StatusNotFound)
			return
		}

		helpers.SendResponse(writer, request, movieToUpdate, http.StatusOK)
	}
}

func (a *App) FindMovieByTitleHandler() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		title, ok := request.URL.Query()["title"]
		if !ok || len(title[0]) < 1 {
			log.Println("Url Param 'title' is missing")
			helpers.SendResponse(writer, request, nil, http.StatusBadRequest)
			return
		}

		movieTitle := models.FindMovieByTitle{}
		if title[0] == "" {
			err := errors.New("title is required")
			log.Printf("Cannot parse movie body. err=%v \n", err)
			helpers.SendResponse(writer, request, nil, http.StatusBadRequest)
			return
		}

		movieTitle.Title = title[0]
		movieFilter := helpers.MapFindMovieByTitleToFindMovie(movieTitle)

		movieList, err := a.findMovieAux(movieFilter)
		if err != nil {
			log.Printf("Cannot find movie in DB. err=%v \n", err)
			helpers.SendResponse(writer, request, nil, http.StatusInternalServerError)
			return
		}

		if len(movieList.Movies) == 0 {
			movie, err := a.findMovieByTitleFromExternalApi(movieTitle.Title)
			if err != nil {
				log.Printf("Cannot find movie in external API. err=%v \n", err)
				helpers.SendResponse(writer, request, nil, http.StatusInternalServerError)
				return
			}

			// Create movie in DB
			movieCreated, err := a.CreateMovieHandler(movie)
			if err != nil {
				log.Printf("Cannot create movie in DB. err=%v \n", err)
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

func (a *App) findMovieByTitleFromExternalApi(movieTitle string) (*models.CreateMovie, error) {
	api := gomdb.Init(config.OmbdApiKey)
	query := &gomdb.QueryData{
		Title: movieTitle,
	}
	movieResult, err := api.MovieByTitle(query)
	if err != nil {
		return &models.CreateMovie{}, err
	}

	movie, err := helpers.MapMovieResultToMovie(movieResult)
	if err != nil {
		return &models.CreateMovie{}, err
	}

	return movie, nil
}

func (a *App) FindMovies() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		movieFilters := models.FindMovie{}

		err := helpers.Parse(writer, request, &movieFilters)

		if err != nil {
			log.Printf("Cannot parse movie body. err=%v \n", err)
			helpers.SendResponse(writer, request, nil, http.StatusBadRequest)
			return
		}

		movieList, err := a.findMovieAux(movieFilters)
		if err != nil {
			log.Printf("Cannot find movie in DB. err=%v \n", err)
			helpers.SendResponse(writer, request, nil, http.StatusInternalServerError)
			return
		}

		helpers.SendResponse(writer, request, movieList, http.StatusOK)
	}
}

func (a *App) findMovieAux(movieFilters models.FindMovie) (models.MovieList, error) {
	filter, err := a.setMovieFilter(movieFilters)
	if err != nil {
		log.Printf("Cannot set movie filter. err=%v \n", err)
		return models.MovieList{}, err
	}

	// Find in DB
	cursor, err := a.DB.Query(config.DbName, config.MovieCollection, &filter, nil)
	if err != nil {
		log.Printf("Cannot find movie in DB. err=%v \n", err)
		return models.MovieList{}, err

	}
	var movieList models.MovieList
	var movies []*models.Movie

	// Iterate cursor
	for cursor.Next(context.TODO()) {
		var movieTemp models.Movie
		err := cursor.Decode(&movieTemp)
		if err != nil {
			log.Printf("Cannot decode movie. err=%v \n", err)
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
		log.Printf("Cannot set filter id. err=%v \n", err)
		return nil, err
	}

	helpers.SetFilter(&filter, models.MovieTitle, movieFilters.Title)

	err = helpers.SetRangeFilter(&filter, models.MovieReleasedYear, movieFilters.ReleasedYearInferiorRange, movieFilters.ReleasedYearSuperiorRange, 32)
	if err != nil {
		log.Printf("Cannot set range filter. err=%v \n", err)
		return nil, err
	}

	// Rating
	err = helpers.SetRangeFilter(&filter, models.MovieRating, movieFilters.RatingInferiorRange, movieFilters.RatingSuperiorRange, 32)
	if err != nil {
		log.Printf("Cannot set range filter. err=%v \n", err)
		return nil, err
	}

	// Genre
	err = helpers.SetFilterArrayString(&filter, models.MovieGenres, movieFilters.Genres)
	if err != nil {
		log.Printf("Cannot set filter array string. err=%v \n", err)
		return nil, err
	}

	return &filter, nil
}
