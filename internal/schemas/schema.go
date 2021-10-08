package schemas

import (
	"io/ioutil"

	"github.com/graph-gophers/graphql-go"
)

var (
	opts = []graphql.SchemaOpt{graphql.UseFieldResolvers()}
)

// ParseSchema parses the schema
// Reads and parses the schema from file.
// Associates root resolver. Panics if can't read.
func ParseSchema(path string, resolver interface{}) (*graphql.Schema, error) {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	schemaString := string(b)
	schema, err := graphql.ParseSchema(
		schemaString,
		resolver,
		opts...,
	)
	if err != nil {
		return nil, err
	}

	return schema, nil
}
