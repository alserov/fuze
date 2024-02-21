package fuze

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
	"time"
)

var (
	s *http.Server
	r *Controller
)

func init() {
	r = NewController()
	s = &http.Server{
		Addr:    fmt.Sprintf(":%d", 3000),
		Handler: r,
	}
}

func TestRouterGet(t *testing.T) {
	r.GET("test", func(c *Ctx) {
		fmt.Println("received get request")
	})

	go func() {
		err := s.ListenAndServe()
		require.NoError(t, err)
	}()

	time.Sleep(time.Millisecond * 200)

	req, err := http.NewRequest(http.MethodGet, "http://127.0.0.1:3000/test", nil)
	require.NoError(t, err)

	client := http.Client{}
	_, err = client.Do(req)
	require.NoError(t, err)
}

func TestRouterGetWithMiddleware(t *testing.T) {
	r.GET("test/{id}/test", func(c *Ctx) {
		fmt.Println("received get request")
		require.Equal(t, "5", c.Parameters["id"])
	}, func(next HandlerFunc) HandlerFunc {
		return func(c *Ctx) {
			fmt.Println("middleware works")
			next(c)
		}
	})

	go func() {
		err := s.ListenAndServe()
		require.NoError(t, err)
	}()

	time.Sleep(time.Millisecond * 200)

	req, err := http.NewRequest(http.MethodGet, "http://127.0.0.1:3000/test/5/test", nil)
	require.NoError(t, err)

	client := http.Client{}
	_, err = client.Do(req)
	require.NoError(t, err)
}

func TestRateLimiterMW(t *testing.T) {
	r.GET("test", func(c *Ctx) {
		fmt.Println("received get request")
	}, WithRateLimitMW(3, 3))

	go func() {
		err := s.ListenAndServe()
		require.NoError(t, err)
	}()

	time.Sleep(time.Millisecond * 200)

	req, err := http.NewRequest(http.MethodGet, "http://127.0.0.1:3000/test", nil)
	require.NoError(t, err)

	client := http.Client{}

	var res *http.Response
	for i := 0; i < 3; i++ {
		res, err = client.Do(req)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, res.StatusCode)
	}

	res, err = client.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusTooManyRequests, res.StatusCode)

	time.Sleep(time.Second)

	res, err = client.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)
}
