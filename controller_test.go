package fuze

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
	"time"
)

var (
	s *http.Server
	r *Router
)

func init() {
	r = NewRouter(NewController())
	s = &http.Server{
		Addr:    fmt.Sprintf(":%d", 3001),
		Handler: r,
	}
}

func TestRouterGet(t *testing.T) {
	r.c.GET("test", func(c *Ctx) {
		fmt.Println("received get request")
	})

	go func() {
		err := s.ListenAndServe()
		require.NoError(t, err)
	}()

	time.Sleep(time.Millisecond * 200)

	req, err := http.NewRequest(http.MethodGet, "http://127.0.0.1:3001/test", nil)
	require.NoError(t, err)

	client := http.Client{}
	_, err = client.Do(req)
	require.NoError(t, err)
}

func TestRouterGetWithParameters(t *testing.T) {
	r.c.GET("test/{id}/{age}/{country}", func(c *Ctx) {
		require.Equal(t, "5", c.Parameters["id"])
	})

	go func() {
		err := s.ListenAndServe()
		require.NoError(t, err)
	}()

	time.Sleep(time.Millisecond * 200)

	req, err := http.NewRequest(http.MethodGet, "http://127.0.0.1:3001/test/5/20/ua", nil)
	require.NoError(t, err)

	client := http.Client{}
	res, err := client.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)
}

func TestRouterGroupGet(t *testing.T) {
	gr := r.c.Group("test")

	gr.GET("{id}", func(c *Ctx) {
		require.Equal(t, "5", c.Parameters["id"])
	})
	gr.GET("/path/{id}", func(c *Ctx) {
		require.Equal(t, "6", c.Parameters["id"])
	})

	go func() {
		err := s.ListenAndServe()
		require.NoError(t, err)
	}()

	time.Sleep(time.Millisecond * 200)

	req, err := http.NewRequest(http.MethodGet, "http://127.0.0.1:3001/test/5", nil)
	require.NoError(t, err)

	client := http.Client{}
	res, err := client.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)

	req, err = http.NewRequest(http.MethodGet, "http://127.0.0.1:3001/test/path/6", nil)
	require.NoError(t, err)

	res, err = client.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)
}

func TestRouterPost(t *testing.T) {
	r.c.POST("test/{name}", func(c *Ctx) {
		require.Equal(t, "user", c.Parameters["name"])

		var req struct {
			Integer int
			String  string
			Boolean bool
		}
		err := c.Decode(&req)
		require.NoError(t, err)

		require.Equal(t, req.Integer, 1)
		require.Equal(t, req.String, "a")
		require.Equal(t, req.Boolean, true)
	})

	go func() {
		err := s.ListenAndServe()
		require.NoError(t, err)
	}()

	time.Sleep(time.Millisecond * 200)

	body, err := json.Marshal(struct {
		Integer int
		String  string
		Boolean bool
	}{
		Integer: 1,
		String:  "a",
		Boolean: true,
	})
	require.NoError(t, err)

	req, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:3001/test/user", bytes.NewReader(body))
	require.NoError(t, err)

	client := http.Client{}
	res, err := client.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)
}
