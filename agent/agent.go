package agent

import (
	"fmt"
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	types "github.com/ohkrab/krab/agent/graphql"
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

	rootQuery := graphql.ObjectConfig{
		Name: "RootQuery",
		Fields: graphql.Fields{
			"migrationSets": &graphql.Field{
				Type: graphql.NewList(types.MigrationSet),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					sets := []*krab.MigrationSet{}
					for _, v := range a.Config.MigrationSets {
						sets = append(sets, v)
					}

					return sets, nil
				},
			},
			"migrationSet": &graphql.Field{
				Type: types.MigrationSet,
				Args: graphql.FieldConfigArgument{
					"refName": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
					// "migration"
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id, ok := p.Args["refName"].(string)
					if ok {
						return a.Config.MigrationSets[id], nil
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
	http.ListenAndServe(":8888", nil)

	return 0
}
