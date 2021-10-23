package agent

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/graphql-go/graphql"
	graphqlHandler "github.com/graphql-go/handler"
	schema "github.com/ohkrab/krab/agent/graphql"
	"github.com/ohkrab/krab/cli"
	"github.com/ohkrab/krab/cliargs"
	"github.com/ohkrab/krab/krab"
	"github.com/unrolled/secure"
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
					for _, set := range a.Config.MigrationSets {

						migrations := []schema.Migration{}
						for _, m := range set.Migrations {
							migrations = append(migrations, schema.Migration{
								RefName:     m.RefName,
								Version:     m.Version,
								Transaction: m.ShouldRunInTransaction(),
								Up:          m.Up.SQL,
								Down:        m.Down.SQL,
							})
						}

						sets = append(sets, schema.MigrationSet{
							RefName:    set.RefName,
							Schema:     set.Schema,
							Migrations: migrations,
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

	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"HEAD", "GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{},
		AllowCredentials: false,
		MaxAge:           300,
	}))
	secureMiddleware := secure.New(secure.Options{
		FrameDeny:          true,
		ContentTypeNosniff: true,
		BrowserXssFilter:   true,
		// ContentSecurityPolicy: "script-src $NONCE",
	})

	r.Use(secureMiddleware.Handler)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	schemaHandler := graphqlHandler.New(&graphqlHandler.Config{
		Schema:     &schema,
		Pretty:     true,
		GraphiQL:   false,
		Playground: true,
	})

	r.Post("/graphql", func(w http.ResponseWriter, r *http.Request) {
		schemaHandler.ServeHTTP(w, r)
	})
	r.Get("/graphql", func(w http.ResponseWriter, r *http.Request) {
		schemaHandler.ServeHTTP(w, r)
	})
	http.ListenAndServe(":8888", r)

	return 0
}
