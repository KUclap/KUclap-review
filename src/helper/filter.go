package helper

import (
	"reflect"
	"strconv"
	"strings"
	"time"

	"kuclap-review-api/src/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func CreateQueryFiltering(query *models.ReviewFilterField) bson.M {

	values := reflect.ValueOf(*query)
	types := values.Type()
	length := values.NumField()

	result := make(bson.M, length)

	for i := 0; i < length; i++ {

		if !values.Field(i).IsZero() {

			field := types.Field(i).Tag.Get("bson")
			value := values.Field(i).Interface()

			if types.Field(i).Tag.Get("type") == "match" {

				result[field] = value

			} else if types.Field(i).Tag.Get("type") == "length" {

				if exQuery, ok := result[field]; ok { // key is exist: Usecase's uses $gte and $lte as range

					operation := types.Field(i).Tag.Get("operation")
					exQuery.(bson.M)[operation] = value
					result[field] = exQuery

				} else {

					query := bson.M{}
					operation := types.Field(i).Tag.Get("operation")
					query[operation] = value
					result[field] = query

				}

			} else if types.Field(i).Tag.Get("type") == "text" {

				query := CreateSubstringFiltering(field, *value.(*string))
				result["$and"] = query

			} else { // date type

				// Parsing timestamp (user sent) to UTC
				timestamp, err := strconv.ParseInt(*value.(*string), 10, 64)
				if err != nil {
					panic(err)
				}

				timeParsed := time.Unix(timestamp/1000, 0)

				if exQuery, ok := result[field]; ok { // if key is exist: Usecase's uses $gte and $lte as range

					operation := types.Field(i).Tag.Get("operation")
					exQuery.(bson.M)[operation] = timeParsed
					result[field] = exQuery

				} else {

					query := bson.M{}
					operation := types.Field(i).Tag.Get("operation")
					query[operation] = timeParsed
					result[field] = query

				}

			}

		}
	}

	return result

}

func CreateSubstringFiltering(field string, text string) []bson.M {

	txt := strings.Split(text, " ")
	filter := make([]bson.M, len(txt))

	for i, subtxt := range txt {

		filter[i] = bson.M{field: bson.M{
			"$regex": primitive.Regex{Pattern: ".*" + subtxt + ".*", Options: "i"},
		}}

	}

	return filter

}
