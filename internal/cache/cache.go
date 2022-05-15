package cache

import "errors"

var ErrorNotFound = errors.New("value not found")

type Cache interface {
	Set(key string, value interface{}) error
	Get(key string) (interface{}, error)
	Clear() error
	Delete(key string) error
}
