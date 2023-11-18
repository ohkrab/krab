package web

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"strings"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/google/uuid"
	"github.com/ohkrab/krab/krab"
	"github.com/ohkrab/krab/krabdb"
	"github.com/ohkrab/krab/krabenv"
	"github.com/ohkrab/krab/krabhcl"
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
	if krabenv.Auth() == "basic" {
		r.Use(middleware.BasicAuth("Krab", krabenv.HttpBasicAuthData()))
	}
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
	r.Get("/images/logo.svg", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/svg+xml")
		w.Header().Set("Cache-Control", "public, max-age=7776000")
		w.Write(s.EmbeddableResources.Logo)
	})
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Content-type") == "application/json" {
			w.Write([]byte(`{"status":"ok"}`))
		} else {
			http.Redirect(w, r, "/ui/databases", http.StatusFound)
		}
	})
	r.Route("/api", func(r chi.Router) {
		r.Post("/actions/execute", func(w http.ResponseWriter, r *http.Request) {
			r.ParseForm()

			form := PostForm(r.PostForm)

			addr := krabhcl.AddrFromStrings([]string{"action", form.Get("namespace"), form.Get("name")})
			action, ok := s.Config.Actions[addr.OnlyRefNames()]
			if !ok {
				http.Error(w, http.StatusText(404), http.StatusNotFound)
				return
			}
			formArgs := form.GetObject("args")
			inputs := krab.NamedInputs{}
			for _, arg := range action.Arguments.Args {
				inputs[arg.Name] = formArgs.Get(arg.Name)
			}

			conn := &krabdb.DefaultConnection{}
			cmd := krab.CmdAction{
				Action:     action,
				Connection: conn,
			}
			_, err := cmd.Do(context.Background(), krab.CmdOpts{
				NamedInputs: inputs,
			})
			if err != nil {
				log.Println(err)
				http.Error(w, http.StatusText(500), http.StatusInternalServerError)
				return
			}

			http.Redirect(w, r, "/ui/actions", http.StatusFound)
		})
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
						  (pg_database_size(d.datname)::numeric / (SUM(pg_database_size(d.datname)) OVER ()))::double precision AS size_percent,
						  pg_encoding_to_char(d.encoding) AS encoding,
						  d.datcollate AS collation,
						  d.datctype AS character_type
						from pg_database d
						inner join pg_tablespace ts on ts.oid = d.dattablespace
						inner join pg_roles auth on auth.oid = d.datdba
						order by size_percent desc, name`
				return db.SelectContext(r.Context(), &data, sql)
			})
			if err != nil {
				log.Println(err)
				http.Error(w, http.StatusText(500), http.StatusInternalServerError)
				return
			}
			for _, db := range data {
				db.CanConnect = !(db.IsTemplate && db.Name == "template0")
			}
			s.render.HTML(w, r, views.LayoutInfo{Footer: krab.InfoVersion}, views.DatabaseList(data))
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
						inner join pg_roles auth on auth.oid = t.spcowner
						left join default_locations dl on dl.name = t.spcname
						order by name`
				return db.SelectContext(r.Context(), &data, sql)
			})

			if err != nil {
				log.Println(err)
				http.Error(w, http.StatusText(500), http.StatusInternalServerError)
				return
			}
			s.render.HTML(w, r, views.LayoutInfo{Footer: krab.InfoVersion}, views.TablespaceList(data))
		})

		r.Get("/databases/{dbname}/schemas", func(w http.ResponseWriter, r *http.Request) {
			dbName := chi.URLParam(r, "dbname")

			conn := s.Connection.(*krabdb.SwitchableDatabaseConnection)
			conn.DatabaseName = dbName

			data := []*dto.SchemaListItem{}
			err := conn.Get(func(db krabdb.DB) error {
				sql := `select
							n.oid AS id,
							n.nspname AS name,
							n.nspowner AS owner_id,
							r.rolname AS owner_name
						from pg_namespace n
						join pg_roles r on n.nspowner = r.oid`
				return db.SelectContext(r.Context(), &data, sql)
			})
			if err != nil {
				log.Println(err)
				http.Error(w, http.StatusText(500), http.StatusInternalServerError)
				return
			}
			for _, schema := range data {
				schema.DatabaseName = dbName
			}
			s.render.HTML(w, r, views.LayoutInfo{Footer: krab.InfoVersion, Nav: views.NavDatabase, Database: dbName}, views.SchemaList(data))
		})

		r.Get("/databases/{dbname}/schemas/{schema}/tables", func(w http.ResponseWriter, r *http.Request) {
			dbName := chi.URLParam(r, "dbname")
			schema := chi.URLParam(r, "schema")

			conn := s.Connection.(*krabdb.SwitchableDatabaseConnection)
			conn.DatabaseName = dbName

			data := []*dto.TableListItem{}
			err := conn.Get(func(db krabdb.DB) error {
				sql := `select
					schemaname as schema_name,
					tablename as name,
					tableowner as owner_name,
					schemaname IN ('information_schema', 'pg_catalog') AS internal,
					coalesce(tablespace, 'pg_default') as tablespace_name,
					rowsecurity as rls,
					pg_size_pretty(pg_relation_size(format('%I.%I', schemaname, tablename))) AS size,
					(
						pg_relation_size(
							format('%I.%I', schemaname, tablename)
						)::numeric /
						SUM( pg_relation_size(format('%I.%I', schemaname, tablename)) ) OVER ()
					)::double precision AS size_percent,
					(select reltuples::bigint from pg_class where oid = (format('%I.%I', schemaname, tablename)::regclass)) AS estimated_rows
				from pg_tables
				where schemaname = $1
				order by schemaname, pg_relation_size(format('%I.%I', schemaname, tablename)) desc, tablename`
				return db.SelectContext(r.Context(), &data, sql, schema)
			})
			if err != nil {
				log.Println(err)
				http.Error(w, http.StatusText(500), http.StatusInternalServerError)
				return
			}
			for _, table := range data {
				table.DatabaseName = dbName
			}
			s.render.HTML(w, r, views.LayoutInfo{Footer: krab.InfoVersion, Nav: views.NavDatabase, Database: dbName}, views.TableList(data))
		})

		r.Get("/actions", func(w http.ResponseWriter, r *http.Request) {
			data := []*dto.ActionListItem{}
			for _, action := range s.Config.Actions {
				args := []*dto.ActionListItemArgument{}
				for _, arg := range action.Arguments.Args {
					args = append(args, &dto.ActionListItemArgument{
						Name:        arg.Name,
						Type:        arg.Type,
						Description: arg.Description,
					})
				}
				data = append(data, &dto.ActionListItem{
					Namespace:   action.Namespace,
					Name:        action.RefName,
					Description: "",
					Transaction: action.Transaction,
					Arguments:   args,
				})
			}
			s.render.HTML(w, r, views.LayoutInfo{Footer: krab.InfoVersion}, views.ActionList(data))
		})

		r.Get("/actions/new/{namespace}/{name}", func(w http.ResponseWriter, r *http.Request) {
			addr := krabhcl.AddrFromStrings([]string{"action", chi.URLParam(r, "namespace"), chi.URLParam(r, "name")})
			action, ok := s.Config.Actions[addr.OnlyRefNames()]
			if !ok {
				http.Error(w, http.StatusText(404), http.StatusNotFound)
				return
			}
			args := []*dto.ActionFormArgument{}
			for _, arg := range action.Arguments.Args {
				args = append(args, &dto.ActionFormArgument{
					Name:        arg.Name,
					Description: arg.Description,
					Value:       "",
				})
			}
			data := dto.ActionForm{
				ExecutionID: uuid.New().String(),
				Namespace:   action.Namespace,
				Name:        action.RefName,
				Arguments:   args,
			}

			s.render.HTML(w, r, views.LayoutInfo{Footer: krab.InfoVersion}, views.ActionForm(&data))
		})

		r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
			s.render.HTML(w, r, views.LayoutInfo{Blank: true}, views.Error404())
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

type PostForm map[string][]string

func (f PostForm) Get(key string) string {
	if _, ok := f[key]; !ok {
		return ""
	}

	return f[key][0]
}

func (f PostForm) GetObject(key string) PostForm {
	prefix := key + "["
	suffix := "]"
	obj := map[string][]string{}
	for k, v := range f {
		if strings.HasPrefix(k, prefix) && strings.HasSuffix(k, suffix) {
			obj[strings.TrimPrefix(strings.TrimSuffix(k, suffix), prefix)] = v
		}
	}

	return PostForm(obj)
}
