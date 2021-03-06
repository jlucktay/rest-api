# Random notes and musings

## Articles that are/were very helpful

I constantly find myself winding up on articles written by Dave Cheney. Big fan of his work!

- [Five suggestions for setting up a Go project][dave-cheney-setup]
- A good pattern for providing optional config: [Functional options for friendly APIs][dave-cheney-options]

## Unused code snippets that seem very handy

### Closure to prepare thing once and use many times

``` go
func (a *apiServer) readPayments() httprouter.Handle {
  // thing := prepareThing()
  return func(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
    // use thing
    w.WriteHeader(http.StatusNotImplemented) // 501
  }
}
```

### Middleware

``` go
s.router.HandleFunc("/admin", s.adminOnly(s.handleAdminIndex()))

func (a *apiServer) adminOnly(h http.HandlerFunc) http.HandlerFunc {
  return func(w http.ResponseWriter, r *http.Request) {
    if !currentUser(r).IsAdmin {
      http.NotFound(w, r)
      return
    }
    h(w, r)
  }
}
```

Get in between one or more requests and their respective handler(s).

### Handler-specific type(s)

Maybe the request/response type(s) are only useful in this specific context, so
define and make use of them in situ, rather than thrown in amongst every other
custom type in the universe of packages.

``` go
func (s *server) handleSomething() http.HandlerFunc {
  type request struct {
    Name string
  }
  type response struct {
    Greeting string `json:"greeting"`
  }
  return func(w http.ResponseWriter, r *http.Request) {
    ...
  }
}
```

### Make use of sync.Once to setup dependencies

Preface handler call wrappers with one-time setup.

Subsequent calls will block until the setup completes, and then just pass
straight through.

``` go
func (s *server) handleTemplate(files string...) http.HandlerFunc {
  var (
    init sync.Once
    tpl  *template.Template
    err  error
  )
  return func(w http.ResponseWriter, r *http.Request) {
    init.Do(func(){
      tpl, err = template.ParseFiles(files...)
    })
    if err != nil {
      http.Error(w, err.Error(), http.StatusInternalServerError)
      return
    }
    // use tpl
  }
}
```

[dave-cheney-setup]: https://dave.cheney.net/2014/12/01/five-suggestions-for-setting-up-a-go-project
[dave-cheney-options]: https://dave.cheney.net/2014/10/17/functional-options-for-friendly-apis
