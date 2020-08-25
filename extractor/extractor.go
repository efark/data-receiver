/*
Package extractor has the interface to extract data from a gin.Context and validate the result.
*/
package extractor

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
)

// Extractor interface can be implemented by struct that have the Extract and Validate methods.
type Extractor interface {
	Extract(c *gin.Context) map[string]string
	Validate(m map[string]string) error
}

// CreateExtractor has a switch to create the right extractor.
func CreateExtractor(class string, params map[string]string) (Extractor, error) {
	var ext Extractor
	var err error
	switch class {
	case "HeaderExtractor":
		ext, err = NewHeaderExtractor(params)
	case "QueryExtractor":
		ext, err = NewQueryExtractor(params)
	default:
		ext, err = NewEmptyExtractor(params)
	}
	return ext, err
}

/*
EmptyExtractor returns an empty map, validation of this map is always right.
*/
type EmptyExtractor struct {
}

// NewEmptyExtractor generates an EmptyExtractor.
func NewEmptyExtractor(_ map[string]string) (EmptyExtractor, error) {
	return EmptyExtractor{}, nil
}

// Extract for an EmptyExtractor returns an empty map.
func (e EmptyExtractor) Extract(c *gin.Context) map[string]string {
	return make(map[string]string)
}

// Validate always return nil.
func (e EmptyExtractor) Validate(_ map[string]string) error {
	return nil
}

/*
HeaderExtractor contains a map of strings which has as key, the key that will be returned and as value the name of the header to look for.
ie: for "signature": "x-signature" -> Header "x-signature"'s value will be retrieved, and will be saved as "signature": "value"
*/
type HeaderExtractor struct {
	config map[string]string
}

// NewHeaderExtractor returns a HeaderExtractor with the parameters from the configuration.
func NewHeaderExtractor(c map[string]string) (HeaderExtractor, error) {
	if len(c) == 0 {
		return HeaderExtractor{}, errors.New("No parameters received.")
	}
	return HeaderExtractor{config: c}, nil
}

// Extract takes the headers from the gin.Context and returns a map with them.
func (h HeaderExtractor) Extract(c *gin.Context) map[string]string {
	m := make(map[string]string)
	for key, value := range h.config {
		m[key] = c.GetHeader(value)
	}
	return m
}

// Validate checks that the values are not empty strings.
func (h HeaderExtractor) Validate(m map[string]string) error {
	for k, v := range m {
		if v == "" {
			return fmt.Errorf("'%s' is empty", k)
		}
	}
	return nil
}

/*
QueryExtractor contains a map of strings which has as key, the key that will be returned and as value the name of the param to look for.
Same logic as Header, but with the query params in the URL.
*/
type QueryExtractor struct {
	config map[string]string
}

// NewQueryExtractor holds the configuration to extract the query parameters from the URL.
func NewQueryExtractor(c map[string]string) (QueryExtractor, error) {
	if len(c) == 0 {
		return QueryExtractor{}, errors.New("No parameters received.")
	}
	return QueryExtractor{config: c}, nil
}

// Extract takes the query parameters from the gin.Context and returns a map with them.
func (q QueryExtractor) Extract(c *gin.Context) map[string]string {
	m := make(map[string]string)
	for key, value := range q.config {
		m[key] = c.Query(value)
	}
	return m
}

// Validate checks that the values are not empty strings.
func (q QueryExtractor) Validate(m map[string]string) error {
	for k, v := range m {
		if v == "" {
			return fmt.Errorf("'%s' is empty", k)
		}
	}
	return nil
}
