package agent

import (
	"fmt"
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
	schema "github.com/ohkrab/krab/agent/graphql"
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
				Type: graphql.NewList(schema.TypeMigrationSet),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					sets := []schema.MigrationSet{}
					for _, v := range a.Config.MigrationSets {
						sets = append(sets, schema.MigrationSet{
							RefName:    v.RefName,
							Schema:     v.Schema,
							Migrations: nil,
						})
					}

					return sets, nil
				},
			},
			"migrationSet": &graphql.Field{
				Type: schema.TypeMigrationSet,
				Args: graphql.FieldConfigArgument{
					"refName": &graphql.ArgumentConfig{
						Type: graphql.String,
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					id, ok := p.Args["refName"].(string)
					if ok {
						set := a.Config.MigrationSets[id]

						migrations := []schema.Migration{}
						for _, v := range set.Migrations {
							migrations = append(migrations, schema.Migration{
								RefName:     v.RefName,
								Version:     v.Version,
								Transaction: v.ShouldRunInTransaction(),
								Up:          v.Up.SQL,
								Down:        v.Down.SQL,
							})
						}

						return schema.MigrationSet{
							RefName:    set.RefName,
							Schema:     set.Schema,
							Migrations: migrations,
						}, nil
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
