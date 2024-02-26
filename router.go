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
		break
	case http.MethodPost:
		h, ok := r.c.post[path]
		if !ok {
			var p Parameters
			if h, p, ok = findLikePath(path, r.c.post); ok {
				h.fn(transformToCtx(w, req, p))
				return
			}
			w.WriteHeader(http.StatusNotFound)
			return
		}
		h.fn(transformToCtx(w, req))
		break
	case http.MethodPut:
		h, ok := r.c.put[path]
		if !ok {
			var p Parameters
			if h, p, ok = findLikePath(path, r.c.put); ok {
				h.fn(transformToCtx(w, req, p))
				return
			}
			w.WriteHeader(http.StatusNotFound)
			return
		}
		h.fn(transformToCtx(w, req))
		break
	case http.MethodDelete:
		h, ok := r.c.delete[path]
		if !ok {
			var p Parameters
			if h, p, ok = findLikePath(path, r.c.delete); ok {
				h.fn(transformToCtx(w, req, p))
				return
			}
			w.WriteHeader(http.StatusNotFound)
			return
		}
		h.fn(transformToCtx(w, req))
		break
	case http.MethodPatch:
		h, ok := r.c.patch[path]
		if !ok {
			var p Parameters
			if h, p, ok = findLikePath(path, r.c.patch); ok {
				h.fn(transformToCtx(w, req, p))
				return
			}
			w.WriteHeader(http.StatusNotFound)
			return
		}
		h.fn(transformToCtx(w, req))
		break
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}
