package fuze

import (
	"context"
	"github.com/google/uuid"
	"net/http"
	"net/url"
	"strings"
)

func transformToCtx(w http.ResponseWriter, req *http.Request, p ...Parameters) *Ctx {
	c := &Ctx{
		Request:  req,
		Response: w,
		ID:       uuid.New().String(),
	}

	if len(p) > 0 {
		c.Parameters = p[0]
	}

	return c
}

func removeFirstSlash(path *string) {
	p := *path
	if p[0] == '/' {
		*path = p[1:]
	}
}

// findLikePath researches for the most alike path in handlers map, if the path is found it will return
// Handler, parsed query parameters and true, otherwise HandlerStruct{}, nil, false
func findLikePath(path string, handlers map[string]HandlerStruct) (HandlerStruct, Parameters, bool) {
	if path[0] == '/' {
		path = path[1:]
	}

	pathEls := strings.Split(path, "/")

	chDone := make(chan struct{}, len(handlers))
	chFound := make(chan struct {
		h HandlerStruct
		p Parameters
	}, 1)
	defer func() {
		close(chFound)
		close(chDone)
	}()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	for _, h := range handlers {
		go func(ctx context.Context, h HandlerStruct) {
			if len(h.pathParameters) < 1 {
				select {
				case <-ctx.Done():
				default:
					chDone <- struct{}{}
				}
				return
			}

			if len(pathEls)-len(h.pathElements) != len(h.pathParameters) {
				select {
				case <-ctx.Done():
				default:
					chDone <- struct{}{}
				}
				return
			}

			for _, pathEl := range h.pathElements {
				if !strings.Contains(path, pathEl) {
					select {
					case <-ctx.Done():
					default:
						chDone <- struct{}{}
					}
					return
				}
			}

			p := parseQueryParameter(&url.URL{Path: path}, h.pathParameters)

			select {
			case <-ctx.Done():
			default:
				chFound <- struct {
					h HandlerStruct
					p Parameters
				}{h: h, p: p}
			}
		}(ctx, h)
	}

	doneCounter := 0
	for {
		select {
		case f := <-chFound:
			cancel()
			return f.h, f.p, true
		case <-chDone:
			doneCounter++
			if doneCounter == len(handlers) {
				cancel()
				return HandlerStruct{}, nil, false
			}
		}
	}
}

func transformPath(path string) (map[int]string, []string) {
	res := make(map[int]string)

	if path[0] == '/' {
		path = path[1:]
	}

	s := strings.Split(path, "/")

	removed := 0

	for i := 0; i < len(s); i++ {
		if len(s[i]) > 0 {
			if s[i][0] == '{' && s[i][len(s[i])-1] == '}' {
				res[i+removed] = s[i][1 : len(s[i])-1]
				s = append(s[:i], s[i+1:]...)
				removed++
				i--
			}
		}
	}

	return res, s
}

func parseQueryParameter(u *url.URL, params map[int]string) Parameters {
	if len(params) < 1 {
		return nil
	}

	path := u.Path
	if path[0] == '/' {
		path = path[1:]
	}

	s := strings.Split(path, "/")

	vals := make(map[string]string)

	for k, v := range params {
		vals[v] = s[k]
	}

	return vals
}
