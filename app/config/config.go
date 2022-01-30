package config

var (
	DbPrefix           string = "mongodb://"
	DbHost             string = "localhost"
	DbPort             string = "27017"
	DbName             string = "movie"
	MgConnectionString string = DbPrefix + DbHost + ":" + DbPort
	MovieCollection    string = "movies"
	OmbdApiKey         string = "40c55e7c"
)
