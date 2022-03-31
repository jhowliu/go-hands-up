// package localcache could help you cache everything you want
package localcache

import "errors"

var ErrKeyNonExist = errors.New("key is not exist")

// Cache is the interface that provides Get() and Set() two methods
type Cache interface {
	// Get the value from cache with key
	Get(string) (interface{}, error)
	// Set the value from cache with key
	Set(string, interface{}) error

	evict(string)
}
