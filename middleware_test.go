package fuze

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
	"time"
)

func TestWithRateLimitMW(t *testing.T) {
	r.c.GET("test", func(c *Ctx) {
		fmt.Println("received get request")
	}, WithRateLimitMW(3, 3))

	go func() {
		err := s.ListenAndServe()
		require.NoError(t, err)
	}()

	time.Sleep(time.Millisecond * 200)

	req, err := http.NewRequest(http.MethodGet, "http://127.0.0.1:3001/test", nil)
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

func TestRouterGetWithMiddleware(t *testing.T) {
	r.c.GET("test/{id}/test", func(c *Ctx) {
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

	req, err := http.NewRequest(http.MethodGet, "http://127.0.0.1:3001/test/5/test", nil)
	require.NoError(t, err)

	client := http.Client{}
	_, err = client.Do(req)
	require.NoError(t, err)
}
