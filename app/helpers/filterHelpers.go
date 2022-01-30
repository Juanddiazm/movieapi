package helpers

import (
	"strconv"
	//"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

)

func SetRangeFilter(filter *bson.M, filterName, inferiorRange, superiorRange string, bigSize int) error {
	if inferiorRange != "" && superiorRange != "" {
		inferiorRangeFloat64, err := strconv.ParseFloat(inferiorRange, bigSize)
		if err != nil {
			return err
		}
		superiorRangeFloat64, err := strconv.ParseFloat(superiorRange, bigSize)
		if err != nil {
			return err
		}
		(*filter)[filterName] = bson.M{"$gte": inferiorRangeFloat64, "$lte": superiorRangeFloat64}
	} else if inferiorRange != "" {
		// Convert administrativeCostsInferiorRange string to float64
		inferiorRangeFloat64, err := strconv.ParseFloat(inferiorRange, bigSize)
		if err != nil {
			return err
		}
		(*filter)[filterName] = bson.M{"$gte": inferiorRangeFloat64}
	} else if superiorRange != "" {
		// Convert administrativeCostsSuperiorRange string to float64
		superiorRangeFloat64, err := strconv.ParseFloat(superiorRange, bigSize)
		if err != nil {
			return err
		}
		(*filter)[filterName] = bson.M{"$lte": superiorRangeFloat64}
	}
	return nil
}

func SetFilter(filter *bson.M, filterName string, filterValue interface{}) {
	if filterValue != "" {
		// log filterName and filterValue in same line
		(*filter)[filterName] = filterValue
	}
}

func SetFilterId(filter *bson.M, filterName string, filterValue string) error {
	if filterValue != "" {
		// Convert filterValue string to ObjectId
		filterValueObjectId, err := primitive.ObjectIDFromHex(filterValue)
		if err != nil {
			return err
		}
		// log filterName and filterValue in same line
		(*filter)[filterName] = filterValueObjectId
	}
	return nil
}

func SetFilterArrayString(filter *bson.M, filterName string, filterValue []string) error {
	if filterValue != nil {
		// filter by filterName using filterValueInt using $in
		(*filter)[filterName] = bson.M{"$in": filterValue}
	}
	return nil
}
