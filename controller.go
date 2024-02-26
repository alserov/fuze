package fuze

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type Controller struct {
	http.Handler
	h struct {
		Handler
		Group
	}

	get    map[string]HandlerStruct
	post   map[string]HandlerStruct
	delete map[string]HandlerStruct
	put    map[string]HandlerStruct
	patch  map[string]HandlerStruct
}

type Handler interface {
	GET(path string, fn HandlerFunc, mw ...Middleware)
	POST(path string, fn HandlerFunc, mw ...Middleware)
	PUT(path string, fn HandlerFunc, mw ...Middleware)
	DELETE(path string, fn HandlerFunc, mw ...Middleware)
}

type Group interface {
	Group(base string) Handler
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

	removeFirstSlash(&path)

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

	removeFirstSlash(&path)

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

	removeFirstSlash(&path)

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

	removeFirstSlash(&path)

	params, pathEls := transformPath(path)

	c.delete["/"+path] = HandlerStruct{
		fn:             fn,
		pathParameters: params,
		pathElements:   pathEls,
	}
}

func (c *Controller) Group(base string) Handler {
	return &group{
		c:    c,
		base: base,
	}
}

type group struct {
	c    *Controller
	base string
}

func (g *group) GET(path string, fn HandlerFunc, mw ...Middleware) {
	for _, mdlwr := range mw {
		fn = mdlwr(fn)
	}

	removeFirstSlash(&path)

	params, pathEls := transformPath(g.base + "/" + path)

	g.c.get[g.base+"/"+path] = HandlerStruct{
		fn:             fn,
		pathParameters: params,
		pathElements:   pathEls,
	}
}

func (g *group) POST(path string, fn HandlerFunc, mw ...Middleware) {
	for _, mdlwr := range mw {
		fn = mdlwr(fn)
	}

	params, pathEls := transformPath(path)

	removeFirstSlash(&path)

	g.c.post["/"+path] = HandlerStruct{
		fn:             fn,
		pathParameters: params,
		pathElements:   pathEls,
	}
}

func (g *group) PUT(path string, fn HandlerFunc, mw ...Middleware) {
	for _, mdlwr := range mw {
		fn = mdlwr(fn)
	}

	params, pathEls := transformPath(path)

	removeFirstSlash(&path)

	g.c.put["/"+path] = HandlerStruct{
		fn:             fn,
		pathParameters: params,
		pathElements:   pathEls,
	}
}

func (g *group) DELETE(path string, fn HandlerFunc, mw ...Middleware) {
	for _, mdlwr := range mw {
		fn = mdlwr(fn)
	}

	params, pathEls := transformPath(path)

	removeFirstSlash(&path)

	g.c.delete["/"+path] = HandlerStruct{
		fn:             fn,
		pathParameters: params,
		pathElements:   pathEls,
	}
}

func (c *Controller) getHandlersAmount() int {
	return len(c.get) + len(c.put) + len(c.post) + len(c.delete)
}

type Ctx struct {
	Request  *http.Request
	Response http.ResponseWriter

	// Parameters represents a map[string]string where key is name of parameter and value is parameter value
	Parameters Parameters

	// ID is being created automatically, is uuid
	ID string
}

// Decode Decodes request body
func (c *Ctx) Decode(target any) error {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return fmt.Errorf("failed to req request body: %w", err)
	}

	if err = json.Unmarshal(body, &target); err != nil {
		return fmt.Errorf("failed to decode request body: %w", err)
	}
	return nil
}

// SendValue Sends value with http status
func (c *Ctx) SendValue(val any, status int) error {
	b, err := json.Marshal(val)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	c.Response.WriteHeader(status)
	if _, err = c.Response.Write(b); err != nil {
		return fmt.Errorf("failed to send response: %w", err)
	}

	return nil
}

// SendStatus Sends http status
func (c *Ctx) SendStatus(status int) {
	c.Response.WriteHeader(status)
}

type Parameters map[string]string

type HandlerStruct struct {
	fn             HandlerFunc
	pathParameters map[int]string
	pathElements   []string
}

type HandlerFunc func(c *Ctx)
