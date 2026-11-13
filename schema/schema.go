package schema

import (
	"github.com/graphql-go/graphql"
	"github.com/wahed-tech/banking-ops-svc/auth"
	"github.com/wahed-tech/go-lib/v2/errors"
)

var errUnauthenticated = errors.New("unauthenticated")

var currentUserType = graphql.NewObject(graphql.ObjectConfig{
	Name: "CurrentUser",
	Fields: graphql.Fields{
		"userId": &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
		"email":  &graphql.Field{Type: graphql.NewNonNull(graphql.String)},
	},
})

func New() (graphql.Schema, error) {
	query := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"currentUser": &graphql.Field{Type: graphql.NewNonNull(currentUserType), Resolve: resolveCurrentUser},
		},
	})

	s, err := graphql.NewSchema(graphql.SchemaConfig{Query: query})
	if err != nil {
		return graphql.Schema{}, errors.Wrap(err).Here()
	}

	return s, nil
}

func resolveCurrentUser(p graphql.ResolveParams) (any, error) {
	caller, ok := auth.CallerFrom(p.Context)
	if !ok {
		return nil, errUnauthenticated.Here()
	}

	return map[string]any{"userId": caller.UserID, "email": caller.Email}, nil
}
