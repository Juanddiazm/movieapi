package database

var (
	dbPrefix = "mongodb://"
	dbHost = "localhost"
	dbPort = "27017"
	mgConnectionString = dbPrefix + dbHost + ":" + dbPort
)
