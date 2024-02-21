package fuze

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
	"time"
)

func TestApp(t *testing.T) {
	a := NewApp()

	a.GET("/path/{id}", func(c *Ctx) {
		fmt.Println(c.Parameters)
	})

	go func() {
		err := a.Run()
		if !errors.Is(err, http.ErrServerClosed) {
			require.NoError(t, err)
		}
	}()

	req, err := http.NewRequest(http.MethodGet, "http://localhost:3000/path/10", nil)
	require.NoError(t, err)

	cl := &http.Client{}

	res, err := cl.Do(req)
	require.NoError(t, err)
	require.Equal(t, http.StatusOK, res.StatusCode)

	time.Sleep(time.Second)

	err = a.GracefulShutdown()
	require.NoError(t, err)
}
