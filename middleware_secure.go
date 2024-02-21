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
