package fuze

import (
	"net/http"
)

type Router struct {
	http.Handler

	c *Controller
}

func NewRouter(ctrl *Controller) *Router {
	return &Router{
		c: ctrl,
	}
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	path := req.URL.Path
	method := req.Method

	switch method {
	case http.MethodGet:
		h, ok := r.c.get[path]
		if !ok {
			var p Parameters
			if h, p, ok = findLikePath(path, r.c.get); ok {
				h.fn(transformToCtx(w, req, p))
				return
			}
			w.WriteHeader(http.StatusNotFound)
			return
		}
		h.fn(transformToCtx(w, req))
	case http.MethodPost:
		h, ok := r.c.post[path]
		if !ok {
			var p Parameters
			if h, p, ok = findLikePath(path, r.c.get); ok {
				h.fn(transformToCtx(w, req, p))
				return
			}
			w.WriteHeader(http.StatusNotFound)
			return
		}
		h.fn(transformToCtx(w, req))
	}
}
