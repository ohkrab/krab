package types

import "github.com/graphql-go/graphql"

var (
	MigrationSet = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "MigrationSet",
			Fields: graphql.Fields{
				"refName": &graphql.Field{
					Type: graphql.String,
				},
				"schema": &graphql.Field{
					Type: graphql.String,
				},
			},
		},
	)

	Migration = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Migration",
			Fields: graphql.Fields{
				"refName": &graphql.Field{
					Type: graphql.String,
				},
				"version": &graphql.Field{
					Type: graphql.String,
				},
				"transaction": &graphql.Field{
					Type: graphql.Boolean,
				},
			},
		},
	)
)
