package api

import (
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
)

func GraphQL(s *graphql.Schema) http.Handler {
	return handler.New(&handler.Config{
		Schema:   s,
		Pretty:   false,
		GraphiQL: false,
	})
}
