package fuze

import (
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

var (
	a *App
)

func init() {
	a = NewApp()

	a.GET("/path/{id}", func(c *Ctx) {})
	a.POST("/path/{id}", func(c *Ctx) {})
	a.DELETE("/path/{id}", func(c *Ctx) {})

	go func() {
		err := a.Run()
		if err != nil {
			panic(err)
		}
	}()
}

func BenchmarkController_GET(b *testing.B) {
	req, err := http.NewRequest(http.MethodGet, "http://localhost:3000/path/10", nil)
	require.NoError(b, err)

	cl := &http.Client{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		res, err := cl.Do(req)
		require.NoError(b, err)
		require.Equal(b, http.StatusOK, res.StatusCode)
	}
}

func BenchmarkController_POST(b *testing.B) {
	req, err := http.NewRequest(http.MethodPost, "http://localhost:3000/path/10", nil)
	require.NoError(b, err)

	cl := &http.Client{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		res, err := cl.Do(req)
		require.NoError(b, err)
		require.Equal(b, http.StatusOK, res.StatusCode)
	}
}

func BenchmarkController_DELETE(b *testing.B) {
	req, err := http.NewRequest(http.MethodDelete, "http://localhost:3000/path/10", nil)
	require.NoError(b, err)

	cl := &http.Client{}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		res, err := cl.Do(req)
		require.NoError(b, err)
		require.Equal(b, http.StatusOK, res.StatusCode)
	}
}
