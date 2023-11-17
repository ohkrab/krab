package web

import (
	"net/http"

	"github.com/a-h/templ"
	"github.com/ohkrab/krab/views"
)

type Renderer struct {
}

func (renderer *Renderer) HTML(w http.ResponseWriter, r *http.Request, info views.LayoutInfo, view templ.Component) {
	views.Layout(info, view).Render(r.Context(), w)
}
