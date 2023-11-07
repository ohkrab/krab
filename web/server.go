package web

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/ohkrab/krab/krab"
	"github.com/ohkrab/krab/krabdb"
	"github.com/ohkrab/krab/views"
	"github.com/ohkrab/krab/web/dto"
)

type Server struct {
	Config     *krab.Config
	Connection krabdb.Connection
	EmbeddableResources
	render Renderer
}

type responseTablespace struct {
	ID        uint64 `json:"ID" db:"id"`
	Name      string `json:"name" db:"name"`
	OwnerID   uint64 `json:"ownerID" db:"owner_id"`
	OwnerName string `json:"ownerName" db:"owner_name"`
	Size      string `json:"size" db:"size"`
	Location  string `json:"location" db:"location"`
}

func (s *Server) Run(args []string) int {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Heartbeat("/health/live"))
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Get("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/x-icon")
		w.Header().Set("Cache-Control", "public, max-age=7776000")
		w.Write(s.EmbeddableResources.Favicon)
	})
	r.Get("/images/logo-white.svg", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/svg+xml")
		w.Header().Set("Cache-Control", "public, max-age=7776000")
		w.Write(s.EmbeddableResources.WhiteLogo)
	})
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})
	r.Route("/ui", func(r chi.Router) {
		r.Use(render.SetContentType(render.ContentTypeJSON))

		r.Get("/databases", func(w http.ResponseWriter, r *http.Request) {
			data := []*dto.DatabaseListItem{}
			err := s.Connection.Get(func(db krabdb.DB) error {
				sql := `select
						  d.oid as id,
						  d.datname AS name,
						  d.datdba AS owner_id,
						  auth.rolname AS owner_name,
						  d.datistemplate AS is_template,
						  d.datconnlimit AS connection_limit,
						  d.dattablespace AS tablespace_id,
						  ts.spcname AS tablespace_name,
						  pg_size_pretty( pg_database_size(d.datname)) as size,
						  pg_encoding_to_char(d.encoding) AS encoding,
						  d.datcollate AS collation,
						  d.datctype AS character_type
						from pg_database d
						inner join pg_tablespace ts on ts.oid = d.dattablespace
						inner join pg_authid auth on auth.oid = d.datdba
						order by name`
				return db.SelectContext(r.Context(), &data, sql)
			})
			if err != nil {
				log.Println(err)
				http.Error(w, http.StatusText(500), http.StatusInternalServerError)
				return
			}
			s.render.HTML(w, r, views.DatabaseList(data))
		})

		r.Get("/tablespaces", func(w http.ResponseWriter, r *http.Request) {
			data := []*dto.TablespaceListItem{}

			err := s.Connection.Get(func(db krabdb.DB) error {
				sql := `
						with default_locations as (
						  select 'pg_default' as name, setting || '/base' as location from pg_settings where name='data_directory'
						  union all
						  select 'pg_global', setting || '/global' from pg_settings where name='data_directory'
						)
						select
						  t.oid AS id,
						  t.spcname AS name,
						  t.spcowner AS owner_id,
						  auth.rolname AS owner_name,
						  pg_size_pretty(pg_tablespace_size(t.oid)) AS size,
						  coalesce(nullif(pg_tablespace_location(t.oid), ''), dl.location) AS location
						from pg_tablespace t
						inner join pg_authid auth on auth.oid = t.spcowner
						left join default_locations dl on dl.name = t.spcname
						order by name`
				return db.SelectContext(r.Context(), &data, sql)
			})

			if err != nil {
				log.Println(err)
				http.Error(w, http.StatusText(500), http.StatusInternalServerError)
				return
			}
			s.render.HTML(w, r, views.TablespaceList(data))
		})

		r.Get("/actions", func(w http.ResponseWriter, r *http.Request) {
			data := []*dto.ActionListItem{
				{
					Namespace:   "db",
					Name:        "create",
					Description: "Create database",
					Transaction: false,
					Arguments: []*dto.ActionListItemArgument{
						{
							Name:        "name",
							Type:        "string",
							Description: "Database name",
						},
						{
							Name:        "user",
							Type:        "string",
							Description: "Database user",
						},
					},
				},
				{
					Namespace:   "user",
					Name:        "create",
					Description: "Create user",
					Transaction: true,
					Arguments: []*dto.ActionListItemArgument{
						{
							Name:        "user",
							Type:        "string",
							Description: "Database user",
						},
						{
							Name:        "password",
							Type:        "string",
							Description: "Database password",
						},
					},
				},
			}
			s.render.HTML(w, r, views.ActionList(data))
		})

		r.Get("/actions/new/{namespace}/{name}", func(w http.ResponseWriter, r *http.Request) {
			data := dto.ActionForm{
				ExecutionID: uuid.New().String(),
				Namespace:   chi.URLParam(r, "namespace"),
				Name:        chi.URLParam(r, "name"),
				Arguments: []*dto.ActionFormArgument{
					{
						Name:        "name",
						Description: "Database name",
						Value:       "test",
					},
					{
						Name:        "user",
						Description: "Database user",
						Value:       "aa",
					},
				},
			}
			s.render.HTML(w, r, views.ActionForm(&data))
		})
	})

	server := &http.Server{Addr: ":8888", Handler: r}
	notifyCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer stop()

	go func() {
		err := server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	<-notifyCtx.Done()
	server.Shutdown(context.Background())

	return 0
}

func (a *Server) Help() string {
	return `Usage: krab agent

Starts server.
`
}

func (a *Server) Synopsis() string {
	return "Krab server"
}

type ErrResponse struct {
	Err            error `json:"-"` // low-level runtime error
	HTTPStatusCode int   `json:"-"` // http response status code

	StatusText string `json:"status"`          // user-level status message
	AppCode    int64  `json:"code,omitempty"`  // application-specific error code
	ErrorText  string `json:"error,omitempty"` // application-level error message, for debugging
}

func ErrRender(err error) render.Renderer {
	return &ErrResponse{
		Err:            err,
		HTTPStatusCode: 422,
		StatusText:     "Error rendering response.",
		ErrorText:      err.Error(),
	}
}

var ErrNotFound = &ErrResponse{HTTPStatusCode: 404, StatusText: "Resource not found."}

func (e *ErrResponse) Render(w http.ResponseWriter, r *http.Request) error {
	render.Status(r, e.HTTPStatusCode)
	return nil
}
