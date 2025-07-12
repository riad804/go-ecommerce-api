package helpers

import (
	"go.mongodb.org/mongo-driver/bson"
)

func StructToBsonMap(obj interface{}) (bson.M, error) {
	data, err := bson.Marshal(obj)
	if err != nil {
		return nil, err
	}

	var result bson.M
	err = bson.Unmarshal(data, &result)
	return result, err
}

func Reduce[T any, R any](arr []T, initial R, fn func(R, T) R) R {
	result := initial
	for _, v := range arr {
		result = fn(result, v)
	}
	return result
}

func Contains[T comparable](list []T, target T) bool {
	for _, item := range list {
		if item == target {
			return true
		}
	}
	return false
}
