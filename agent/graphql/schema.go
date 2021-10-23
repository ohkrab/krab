package schema

import "github.com/graphql-go/graphql"

var (
	TypeMigrationSet = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "MigrationSet",
			Fields: graphql.Fields{
				"refName": &graphql.Field{
					Type: graphql.String,
				},
				"schema": &graphql.Field{
					Type: graphql.String,
				},
				"migrations": &graphql.Field{
					Type: graphql.NewList(TypeMigration),
				},
			},
		},
	)

	TypeMigration = graphql.NewObject(
		graphql.ObjectConfig{
			Name: "Migration",
			Fields: graphql.Fields{
				"refName": &graphql.Field{
					Type: graphql.String,
				},
				"up": &graphql.Field{
					Type: graphql.String,
				},
				"down": &graphql.Field{
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
