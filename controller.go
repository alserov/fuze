package fuze

import (
	"net/http"
)

type Controller struct {
	http.Handler
	h Handler
	//Group(base string) Handler

	get    map[string]HandlerStruct
	post   map[string]HandlerStruct
	delete map[string]HandlerStruct
	put    map[string]HandlerStruct
}

type Handler interface {
	GET(path string, fn HandlerFunc, mw ...Middleware)
	POST(path string, fn HandlerFunc, mw ...Middleware)
	PUT(path string, fn HandlerFunc, mw ...Middleware)
	DELETE(path string, fn HandlerFunc, mw ...Middleware)
}

func NewController() *Controller {
	return &Controller{
		get:    make(map[string]HandlerStruct),
		put:    make(map[string]HandlerStruct),
		delete: make(map[string]HandlerStruct),
		post:   make(map[string]HandlerStruct),
	}
}

func (c *Controller) GET(path string, fn HandlerFunc, mw ...Middleware) {
	for _, mdlwr := range mw {
		fn = mdlwr(fn)
	}

	params, pathEls := transformPath(path)

	c.get["/"+path] = HandlerStruct{
		fn:             fn,
		pathParameters: params,
		pathElements:   pathEls,
	}
}

func (c *Controller) POST(path string, fn HandlerFunc, mw ...Middleware) {
	for _, mdlwr := range mw {
		fn = mdlwr(fn)
	}

	params, pathEls := transformPath(path)

	c.post["/"+path] = HandlerStruct{
		fn:             fn,
		pathParameters: params,
		pathElements:   pathEls,
	}
}

func (c *Controller) PUT(path string, fn HandlerFunc, mw ...Middleware) {
	for _, mdlwr := range mw {
		fn = mdlwr(fn)
	}

	params, pathEls := transformPath(path)

	c.put["/"+path] = HandlerStruct{
		fn:             fn,
		pathParameters: params,
		pathElements:   pathEls,
	}
}

func (c *Controller) DELETE(path string, fn HandlerFunc, mw ...Middleware) {
	for _, mdlwr := range mw {
		fn = mdlwr(fn)
	}

	params, pathEls := transformPath(path)

	c.delete["/"+path] = HandlerStruct{
		fn:             fn,
		pathParameters: params,
		pathElements:   pathEls,
	}
}

func (c *Controller) getHandlersAmount() int {
	return len(c.get) + len(c.put) + len(c.post) + len(c.delete)
}

//	func (r *router) Group(base string) Handler {
//		return &handler{
//			base: base,
//		}
//	}
//
//	type handler struct {
//		router router
//		base   string
//	}
//
//	func (h *handler) GET(path string, fn http.HandlerFunc, mw ...Middleware) {
//		for _, mdlwr := range mw {
//			fn = mdlwr(fn).ServeHTTP
//		}
//
//		h.router.get[h.base+"/"+path] = fn
//	}
//
//	func (h *handler) POST(path string, fn http.HandlerFunc, mw ...Middleware) {
//		for _, mdlwr := range mw {
//			fn = mdlwr(fn).ServeHTTP
//		}
//
//		h.router.post[h.base+"/"+path] = fn
//	}
//
//	func (h *handler) PUT(path string, fn http.HandlerFunc, mw ...Middleware) {
//		for _, mdlwr := range mw {
//			fn = mdlwr(fn).ServeHTTP
//		}
//
//		h.router.put[h.base+"/"+path] = fn
//	}
//
//	func (h *handler) DELETE(path string, fn http.HandlerFunc, mw ...Middleware) {
//		for _, mdlwr := range mw {
//			fn = mdlwr(fn).ServeHTTP
//		}
//
//		h.router.delete[h.base+"/"+path] = fn
//	}

type Ctx struct {
	Request  *http.Request
	Response http.ResponseWriter

	// Parameters represents a map[string]string where key is name of parameter and value is parameter value
	Parameters Parameters

	// ID is being created automatically, is uuid
	ID string
}

type Parameters map[string]string

type HandlerStruct struct {
	fn             HandlerFunc
	pathParameters map[int]string
	pathElements   []string
}

type HandlerFunc func(c *Ctx)
