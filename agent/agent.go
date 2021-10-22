package agent

import (
	"fmt"
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	"github.com/ohkrab/krab/cli"
	"github.com/ohkrab/krab/cliargs"
	"github.com/ohkrab/krab/krab"
)

// ActionMigrateDown keeps data needed to perform this action.
type Agent struct {
	Config *krab.Config
}

func (a *Agent) Help() string {
	return `Usage: krab agent
  
Starts the agent mode.
`
}

func (a *Agent) Synopsis() string {
	return fmt.Sprintf("Start agent")
}

// Run in CLI.
func (a *Agent) Run(args []string) int {
	ui := cli.DefaultUI()
	flags := cliargs.New(args)
	flags.RequireNonFlagArgs(0)

	err := flags.Parse()
	if err != nil {
		ui.Output(a.Help())
		ui.Error(err.Error())
		return 1
	}
	// Schema
	migrationSetType := graphql.NewObject(
		graphql.ObjectConfig{
			Name: "MigrationSet",
			Fields: graphql.Fields{
				"ref_name": &graphql.Field{
					Type: graphql.String,
				},
				"schema": &graphql.Field{
					Type: graphql.String,
				},
			},
		},
	)

	rootQuery := graphql.ObjectConfig{
		Name: "RootQuery",
		Fields: graphql.Fields{
			"migration_set": &graphql.Field{
				Type: migrationSetType,
				Args: graphql.FieldConfigArgument{
					"id": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id, ok := p.Args["ref_name"]
					if ok {
						return
					}
					return nil, nil
				},
			},
		},
	}
	schemaConfig := graphql.SchemaConfig{
		Query: graphql.NewObject(rootQuery),
	}
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		ui.Error(fmt.Errorf("failed to create new schema, error: %v", err).Error())
		return 1
	}

	h := handler.New(&handler.Config{
		Schema:     &schema,
		Pretty:     true,
		GraphiQL:   false,
		Playground: true,
	})

	http.Handle("/graphql", h)
	http.ListenAndServe(":8080", nil)

	return 0
}
