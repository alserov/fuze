# Fuze 

<div>
    <img src="docs/fuze.jpg" alt="drawing" style="height:100px;"/>
    <h1>💥Fuze - HTTP router and URL matcher for building Go web apps</h1>
</div>

--------

# Examples

## Default server template

```go
package main

import (
	"fmt"
	"github.com/alserov/fuze"
)

func main() {
	a := fuze.NewApp()

	a.GET("/path", func(c *fuze.Ctx) {
		fmt.Println("hello!")
	})

	err := a.Run()
	if err != nil {
		panic(err)
	}
}
 ```  

## Routing

```go

package main

import (
	"fmt"
	"github.com/alserov/fuze"
)

func main() {
	a := fuze.NewApp()

	gr := a.Group("base")

	gr.GET("/{id}", func(c *fuze.Ctx) {
		// your code
	})
	gr.GET("/path/{id}", func(c *fuze.Ctx) {
		// your code
	})

	err := a.Run()
	if err != nil {
		panic(err)
	}
}

```

## Middleware

```go

package main

import (
	"github.com/alserov/fuze"
)

func main() {
	a := fuze.NewApp()

	// or create your own mw that will be of type Middleware
	a.Get("/path", func(c *fuze.Ctx) {
             // your code
	}, fuze.WithRateLimitMW(3, 3))


	err := a.Run()
	if err != nil {
		panic(err)
	}
}

```
