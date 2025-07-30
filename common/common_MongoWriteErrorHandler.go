package common

import (
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrMongoDuplicateKey = errors.New("duplicate key error")
	ErrMongoUnknown      = errors.New("unknown error")
)

func NewErrDuplicateKey(err error) error {
	return fmt.Errorf("%s: %w", ErrMongoDuplicateKey.Error(), err)
}

func NewErrUnknown(err error) error {
	return fmt.Errorf("%s: %w", ErrMongoUnknown.Error(), err)
}

// Mapper

type MongoErrorMapper struct {
	ErrDuplicateKey error
}

func NewMongoErrorMapper() *MongoErrorMapper {
	return new(MongoErrorMapper)
}

func (m *MongoErrorMapper) WithErrDuplicateKey(err error) *MongoErrorMapper {
	m.ErrDuplicateKey = err
	return m
}

// Handler

func MongoWriteErrorHandler(err error, mapper *MongoErrorMapper) error {
	if err == nil {
		return nil
	}

	var writeException mongo.WriteException
	if errors.As(err, &writeException) {
		for _, we := range writeException.WriteErrors {
			switch we.Code {
			case 11000:
				if mapper != nil && mapper.ErrDuplicateKey != nil {
					return mapper.ErrDuplicateKey
				}
				return NewErrDuplicateKey(err)
			}
		}
	}

	return NewErrUnknown(err)
}
