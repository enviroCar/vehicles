package main

import (
	"encoding/json"
	"errors"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"log"
	"net/http"
	"net/url"
	"reflect"
	"runtime"

	"github.com/gorilla/mux"
)

var _ http.Handler = (*Server)(nil)

// HandlerFunc is an http.HandlerFunc with an Context.
type HandlerFunc func(*Context) (interface{}, error)

// Server is the HTTP server.
type Server struct {
	router     *mux.Router
	routeByPtr map[uintptr]*mux.Route
}

// NewServer creates a new Server.
func NewServer() *Server {
	s := &Server{
		routeByPtr: make(map[uintptr]*mux.Route),
		router:     mux.NewRouter().StrictSlash(true),
	}
	s.router.Use(s.loggingMiddleware())
	s.router.Use(mux.CORSMethodMiddleware(s.router))
	s.router.MethodNotAllowedHandler = s.errorHandler(nil, ErrMethodNotAllowed)
	s.router.NotFoundHandler = s.errorHandler(nil, ErrNotFound)
	return s
}

func (s *Server) loggingMiddleware() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Printf("%v %v", r.Method, r.RequestURI)
			next.ServeHTTP(w, r)
		})
	}
}

// Get defines a HTTP GET route.
func (s *Server) Get(path string, handlerFunc HandlerFunc) {
	pc := reflect.ValueOf(handlerFunc).Pointer()
	log.Printf("Registering route: %v\n", path)
	route := s.router.
		Host("{host:.+}").
		Path(path).
		Name(runtime.FuncForPC(pc).Name()).
		Handler(s.handler(handlerFunc)).
		Methods(http.MethodGet)

	s.routeByPtr[pc] = route

}

// Start starts the server.
func (s *Server) Start(addr string) error {
	return http.ListenAndServe(addr, s)
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

// Context is a HTTP context.
type Context struct {
	server  *Server
	Params  map[string]string
	Request *http.Request
	logger  *logrus.Entry
}

// URL returns a URL builder function for the specified handler.
func (c *Context) URL(handler HandlerFunc) func(...string) (*url.URL, error) {
	if route, ok := c.server.routeByPtr[reflect.ValueOf(handler).Pointer()]; ok {
		return func(pairs ...string) (*url.URL, error) {
			params := []string{"host", c.Request.Host}

			params = append(params, pairs...)
			url, err := route.URL(params...)
			if err != nil {
				return nil, err
			}
			prefix := c.Request.Header.Get("X-Forwarded-Prefix")
			if prefix != "" {
				url.Path = prefix + url.Path
			}
			proto := c.Request.Header.Get("X-Forwarded-Proto")
			if proto != "" {
				url.Scheme = proto
			}
			port := c.Request.Header.Get("X-Forwarded-Port")
			if port != "" && url.Port() != port &&
				((url.Scheme == "https" && port != "443") ||
					(url.Scheme == "http" && port != "80")) {
				url.Host = url.Hostname() + ":" + port
			}

			return url, nil
		}
	}
	return func(...string) (*url.URL, error) {
		return nil, errors.New("route not found")
	}
}

func (*Server) errorHandler(ctxlogger *logrus.Entry, err error) http.Handler {

	if ctxlogger == nil {
		ctxlogger = logrus.NewEntry(logrus.StandardLogger())
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		status := http.StatusInternalServerError
		if e, ok := err.(Error); ok {
			status = e.Status()
		}
		content := &struct {
			StatusCode int    `json:"statusCode"`
			StatusText string `json:"statusText"`
			Message    string `json:"message"`
		}{
			StatusCode: status,
			StatusText: http.StatusText(status),
			Message:    err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		if err := json.NewEncoder(w).Encode(content); err != nil {
			ctxlogger.WithError(err).Error("could not encode error response")
		}
	})
}

func (*Server) redirectHandler(redirect *Redirect) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Location", redirect.Location.String())
		w.WriteHeader(redirect.Status)
	})
}

func (*Server) contentHandler(ctxlogger *logrus.Entry, content interface{}) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if content == nil {
			w.WriteHeader(http.StatusNoContent)
		} else {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(content); err != nil {
				ctxlogger.WithError(err).Error("could not encode content response")
			}
		}
	})
}

func (s *Server) handler(f HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		requestId := r.Header.Get("X-Request-ID")
		if requestId == "" {
			requestId = uuid.NewV4().String()
		}

		logger := logrus.New()
		ctxlogger := logger.WithFields(logrus.Fields{
			"request-id": requestId,
		})

		var handler http.Handler
		content, err := f(&Context{s, mux.Vars(r), r, ctxlogger})
		if err != nil {
			handler = s.errorHandler(ctxlogger, err)
		} else if redirect, ok := content.(*Redirect); ok {
			handler = s.redirectHandler(redirect)
		} else {
			handler = s.contentHandler(ctxlogger, content)
		}
		handler.ServeHTTP(w, r)
	}
}
