package fuze

type Middleware func(next HandlerFunc) HandlerFunc
