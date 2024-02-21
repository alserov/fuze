package fuze

import (
	"golang.org/x/time/rate"
	"net/http"
)

// WithRateLimitMW
// lim represents the rate (tokens per second)
// burst is the burst size (maximum number of tokens allowed in the bucket)
func WithRateLimitMW(lim int, burst int) Middleware {
	l := rate.NewLimiter(rate.Limit(lim), burst)
	return func(next HandlerFunc) HandlerFunc {
		return func(c *Ctx) {
			if !l.Allow() {
				c.Response.WriteHeader(http.StatusTooManyRequests)
				return
			}
			next(c)
		}
	}
}

// WithCookies
// checks if all the required cookies present
// names are names of cookies that must be in request
func WithCookies(names ...string) Middleware {
	return func(next HandlerFunc) HandlerFunc {
		return func(c *Ctx) {
			for _, n := range names {
				if _, err := c.Request.Cookie(n); err != nil {
					c.Response.WriteHeader(http.StatusMethodNotAllowed)
					return
				}
			}
			next(c)
		}
	}
}
