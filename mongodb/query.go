package mongodb

import (
	"reflect"

	"github.com/leetatech/leeta_golang_libraries/query/filter"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// BuildMongoFilterQuery constructs a MongoDB filter query based on the provided request filter
// and a field mapping. It supports both "and" and "or" operators for combining field conditions.
// TODO: handle filter field mapping better
func BuildMongoFilterQuery(requestFilter *filter.Request, fieldMapping map[string]string) bson.M {
	query := bson.M{}

	if requestFilter == nil {
		return query
	}

	// Helper function to build individual field queries
	buildFieldQuery := func(field filter.RequestField) bson.M {
		// Use the mapped field name if it exists, otherwise use the original field name
		fieldName := fieldMapping[field.Name]
		if fieldName == "" {
			fieldName = field.Name
		}
		if field.Operator == filter.CompareOperatorContains || reflect.TypeOf(field.Value).Kind() == reflect.Slice {
			return bson.M{fieldName: bson.M{"$in": field.Value}}
		}
		return bson.M{fieldName: field.Value}
	}

	switch requestFilter.Operator {
	case filter.LogicOperatorAnd:
		for _, field := range requestFilter.Fields {
			fieldQuery := buildFieldQuery(field)
			for key, value := range fieldQuery {
				query[key] = value
			}
		}
	case filter.LogicOperatorOr:
		orConditions := make([]bson.M, len(requestFilter.Fields))
		for i, field := range requestFilter.Fields {
			orConditions[i] = buildFieldQuery(field)
		}
		query["$or"] = orConditions
	}

	return query
}

// GetPaginatedOpts returns MongoDB find options for pagination.
// It calculates the number of documents to skip and the limit based on the given page size and page index.
// If pageIndex or pageSize are less than 1, it sets sensible defaults.
func GetPaginatedOpts(pageSize, pageIndex int64) *options.FindOptions {
	if pageIndex < 1 {
		pageIndex = 1
	}
	if pageSize < 1 {
		pageSize = 10 // or a sensible default
	}

	skip := pageSize * (pageIndex - 1)
	limit := pageSize

	opts := options.Find()
	opts.SetSkip(skip)
	opts.SetLimit(limit)
	return opts
}
